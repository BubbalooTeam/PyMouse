package http

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func Request(options RequestOptions) (*fasthttp.Request, *fasthttp.Response, error) {
	AcRequest := fasthttp.AcquireRequest()
	AcResponse := fasthttp.AcquireResponse()

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

	// Treat the "Body" if the POST method is used.
	if options.Method == "POST" {
		AcRequest.SetBody(options.Body)
	}

	// Makes the Request according to the passed method.
	if options.Method == "GET" ||
		options.Method == "POST" ||
		options.Method == "OPTIONS" {
		err := fasthttp.Do(AcRequest, AcResponse)
		if err != nil {
			return nil, nil, err
		}
	} else {
		return nil, nil, fmt.Errorf("[http/Request]: HTTP Method not Supported: %s", options.Method)
	}
	return AcRequest, AcResponse, nil
}
