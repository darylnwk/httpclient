package httpclient

import "net/http"

// Posthook defines a hook called after a HTTP request
type Posthook func(response *http.Response, err error)
