package rapidhttp

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

// Request makes an HTTP request using fasthttp
func Request(options RequestOptions) (*fasthttp.Response, error) {
	AcRequest := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(AcRequest)

	AcResponse := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(AcResponse)

	AcRequest.SetRequestURI(options.URL)
	AcRequest.Header.SetMethod(options.Method)

	// Set Request Headers
	for key, value := range options.Headers {
		AcRequest.Header.Set(key, value)
	}

	// Set Request Params
	QueryArgs := AcRequest.URI().QueryArgs()
	for key, value := range options.Params {
		QueryArgs.Set(key, value)
	}

	// Treat the "Body" if the POST method is used
	if options.Method == "POST" {
		AcRequest.SetBody(options.Body)
	}

	// Makes the Request according to the passed method
	switch options.Method {
	case "GET", "POST", "OPTIONS":
		err := fasthttp.Do(AcRequest, AcResponse)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("[http/Request]: HTTP Method not supported: %s", options.Method)
	}

	return AcResponse, nil
}
