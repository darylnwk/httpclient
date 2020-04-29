package httpclient_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/darylnwk/retry"
	"github.com/stretchr/testify/assert"

	"github.com/darylnwk/httpclient"
)

func TestClient_NewClient(t *testing.T) {
	client := httpclient.NewClient(
		httpclient.OptionTimeout(10*time.Second),
		httpclient.OptionAttempts(1),
		httpclient.OptionHTTPClient(&http.Client{Timeout: 20 * time.Second}),
	)

	assert.Equal(t, 20*time.Second, client.Client.Timeout)
	assert.Equal(t, uint(1), client.Retryer.Attempts)

	client = httpclient.NewClient(
		httpclient.OptionHTTPClient(&http.Client{Timeout: 20 * time.Second}),
		httpclient.OptionTimeout(10*time.Second),
		httpclient.OptionAttempts(2),
	)

	assert.Equal(t, 10*time.Second, client.Client.Timeout)
	assert.Equal(t, uint(2), client.Retryer.Attempts)

	client = httpclient.NewClient(
		httpclient.OptionTimeout(10*time.Second),
		httpclient.OptionAttempts(3),
		httpclient.OptionDelay(time.Second),
		httpclient.OptionBackoff(retry.NoBackoff),
		httpclient.OptionJitter(retry.NoJitter),
	)

	assert.Equal(t, 10*time.Second, client.Client.Timeout)
	assert.Equal(t, uint(3), client.Retryer.Attempts)
}

func TestClient_Do(t *testing.T) {
	var (
		url = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"hello": "world"}`))
		})).URL
		client        = httpclient.NewClient()
		request, _    = http.NewRequest(http.MethodPost, url, strings.NewReader(`{"foo": "bar"}`))
		response, err = client.Do(request)
		b, _          = ioutil.ReadAll(response.Body)
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, `{"hello": "world"}`, string(b))
}

func TestClient_DoHTTPDoRequestError(t *testing.T) {
	var (
		client        = httpclient.NewClient()
		request, _    = http.NewRequest(http.MethodPost, "", strings.NewReader(`{"foo": "bar"}`))
		response, err = client.Do(request)
	)

	assert.EqualError(t, err, "httpclient: request occurred with errors: Post : unsupported protocol scheme \"\"")
	assert.Nil(t, response)
}

func TestClient_DoWithInternalServerError(t *testing.T) {
	var (
		url = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})).URL
		client = httpclient.NewClient(
			httpclient.OptionAttempts(2),
		)
		request, _    = http.NewRequest(http.MethodPost, url, strings.NewReader(`{"foo": "bar"}`))
		response, err = client.Do(request)
	)

	assert.EqualError(t, err, "httpclient: request occurred with errors: retrying on 500 Internal Server Error; retrying on 500 Internal Server Error")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
}

func TestClient_DoWithHook(t *testing.T) {
	var (
		url = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})).URL
		client = httpclient.NewClient(
			httpclient.OptionAttempts(2),
			httpclient.OptionAddPrehook(func(req *http.Request) {
				assert.Equal(t, url, req.URL.String())
			}),
			httpclient.OptionAddPosthook(func(resp *http.Response, err error) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}),
		)
		request, _    = http.NewRequest(http.MethodPost, url, strings.NewReader(`{"foo": "bar"}`))
		response, err = client.Do(request)
	)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
