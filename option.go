package httpclient

import (
	"net/http"
	"time"
)

// Option defines `Client` option
type Option func(*Client)

// OptionDelay sets HTTP client delay between attempts
func OptionDelay(delay time.Duration) Option {
	return func(client *Client) {
		client.Retryer.Delay = delay
	}
}

// OptionTimeout sets HTTP client timeout
func OptionTimeout(timeout time.Duration) Option {
	return func(client *Client) {
		client.Client.Timeout = timeout
	}
}

// OptionAttempts sets retry attempts
func OptionAttempts(attempts uint) Option {
	return func(client *Client) {
		client.Retryer.Attempts = attempts
	}
}

// OptionHTTPClient sets HTTP client
func OptionHTTPClient(httpClient *http.Client) Option {
	return func(client *Client) {
		client.Client = httpClient
	}
}

// OptionBackoff sets `Client` retry backoff
func OptionBackoff(fn func(n uint, delay time.Duration) time.Duration) Option {
	return func(client *Client) {
		client.Retryer.Backoff = fn
	}
}

// OptionJitter sets `Client` retry jitter
func OptionJitter(fn func(backoff time.Duration) time.Duration) Option {
	return func(client *Client) {
		client.Retryer.Jitter = fn
	}
}
