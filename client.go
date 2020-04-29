package httpclient

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/darylnwk/retry"
)

const (
	defaultRetryAttempts = 1
	defaultTimeout       = 30 * time.Second
)

// Client defines a HTTP client
type Client struct {
	Client    *http.Client
	Retryer   retry.Retryer
	Prehooks  []Prehook
	Posthooks []Posthook
}

// NewClient initialises a new `Client`
func NewClient(opts ...Option) *Client {
	client := &Client{
		Client: &http.Client{
			Timeout: defaultTimeout,
		},
		Retryer: retry.Retryer{
			Attempts: defaultRetryAttempts,
		},
		Prehooks:  []Prehook{},
		Posthooks: []Posthook{},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Do performs HTTP request and returns HTTP response if exists
func (client *Client) Do(request *http.Request) (*http.Response, error) {
	var (
		response *http.Response
		timeout  bool

		success, errs = client.Retryer.Do(func() error {
			var err error

			for _, prehook := range client.Prehooks {
				prehook(request)
			}

			response, err = client.Client.Do(request)

			for _, posthook := range client.Posthooks {
				posthook(response, err)
			}

			// Retry only on 5xx status codes
			if response != nil && response.StatusCode >= http.StatusInternalServerError {
				return fmt.Errorf("retrying on %s", response.Status)
			}

			return err
		})
	)

	if !success {
		// Check if last error is a timeout error
		if err, ok := errs[len(errs)-1].(net.Error); ok && err.Timeout() {
			timeout = true
		}

		return response, &httpError{
			err:     fmt.Sprintf("httpclient: request occurred with errors: %s", errs),
			timeout: timeout,
		}
	}

	return response, nil
}
