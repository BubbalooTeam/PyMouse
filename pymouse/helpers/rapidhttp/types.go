package rapidhttp

type HTTPGetStruct struct {
	Params map[string]string // Params for your request with method "GET".
}

type HTTPPostStruct struct {
	Json map[string]string // Json param for your request with method "POST".
}

type HTTPStruct struct {
	Method     string            // "GET" or "POST".
	URL        string            // URL param for your request.
	Headers    map[string]string // Headers params for your request.
	GETParams  *HTTPGetStruct    // Get params for your request.
	POSTParams *HTTPPostStruct   // Post params for your request.
}
