package rapidhttp

type RequestOptions struct {
	URL     string
	Method  string
	Headers map[string]string
	Params  map[string]string
	Body    []byte
}
