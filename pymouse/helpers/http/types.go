package http

type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Params  map[string]string
	Body    []byte
}
