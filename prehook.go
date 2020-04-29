package httpclient

import "net/http"

// Prehook defines a hook called before a HTTP request
type Prehook func(request *http.Request)
