package httpclient

type httpError struct {
	err     string
	timeout bool
}

func (e *httpError) Error() string {
	return e.err
}

func (e *httpError) Timeout() bool {
	return e.timeout
}

func (e *httpError) Temporary() bool {
	return true
}
