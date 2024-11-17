package rapidhttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"pymouse/pymouse/config"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
)

var (
	HTTPTransport *http.Transport
	HTTPClient    *http.Client
)

func GetHTTPTransport() *http.Transport {
	HTTPTransport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	if config.Socks5Proxy != "" {
		ProxyURL, err := url.Parse(config.Socks5Proxy)
		if err != nil {
			logrus.Errorf("An error occurred while parsing your proxy while trying to configure the common HTTPClient: %v", err)
			err = http2.ConfigureTransport(HTTPTransport)
			if err != nil {
				logrus.Errorf("An error occurred while trying to configure the common HTTPClient: %v", err)
				return nil
			}
			return HTTPTransport
		}
		HTTPTransport = &http.Transport{
			Proxy: http.ProxyURL(ProxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	err := http2.ConfigureTransport(HTTPTransport)
	if err != nil {
		logrus.Errorf("An error occurred while trying to configure the common HTTPClient: %v", err)
		return nil
	}

	return HTTPTransport
}

func GetHTTPClient() *http.Client {
	HTTPTransport = GetHTTPTransport()
	if HTTPTransport != nil {
		HTTPClient = &http.Client{
			Transport: HTTPTransport,
			Timeout:   time.Duration(15) * time.Second,
		}
		return HTTPClient
	}
	return &http.Client{
		Timeout: time.Duration(15) * time.Second,
	}
}

func Request(Client *http.Client, HTTPParams HTTPStruct) (*http.Response, error) {
	if Client == nil {
		logrus.Error("HTTP Client is not found.")
		return nil, fmt.Errorf("HTTP Client was not Passed to the Request Function to Process")
	}

	parsedURL := HTTPParams.URL
	if strings.Contains(HTTPParams.Method, "GET") && HTTPParams.GETParams != nil {
		params := url.Values{}
		for k, v := range HTTPParams.GETParams.Params {
			params.Add(k, v)
		}
		if len(params) > 0 {
			parsedURL += "?" + params.Encode()
		}
	}

	if strings.Contains(HTTPParams.Method, "POST") && HTTPParams.POSTParams != nil {
		jsonData, err := json.Marshal(HTTPParams.POSTParams.Json)
		if err != nil {
			logrus.Errorf("Error in Marshaling POSTParams: %v", err)
			return nil, err
		}
		request, err := http.NewRequest(HTTPParams.Method, parsedURL, bytes.NewBuffer(jsonData))
		if err != nil {
			logrus.Errorf("Error in creating a new request: %v", err)
			return nil, err
		}
		if HTTPParams.Headers != nil {
			for k, v := range HTTPParams.Headers {
				request.Header.Set(k, v)
			}
		}

		return Client.Do(request)
	}

	request, err := http.NewRequest(HTTPParams.Method, parsedURL, nil)
	if err != nil {
		logrus.Errorf("Error in creating a new request: %v", err)
		return nil, err
	}

	if HTTPParams.Headers != nil {
		for k, v := range HTTPParams.Headers {
			request.Header.Set(k, v)
		}
	}

	return Client.Do(request)
}
