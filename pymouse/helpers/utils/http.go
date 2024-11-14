package utils

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
	"pymouse/pymouse/config"

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
			log.Printf("[utils/GetHTTPTransport][Error]: An error occurred while parsing your proxy while trying to configure the common HTTPClient: %v", err)
			err = http2.ConfigureTransport(HTTPTransport)
			if err != nil {
				log.Printf("[utils/GetHTTPTransport][Error]: An error occurred while trying to configure the common HTTPClient: %v", err)
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
		log.Printf("[utils/GetHTTPTransport][Error]: An error occurred while trying to configure the common HTTPClient: %v", err)
		return nil
	}

	return HTTPTransport
}

func GetHTTPClient() *http.Client {
	HTTPTransport = GetHTTPTransport()
	if HTTPTransport != nil {
		HTTPClient = &http.Client{
			Transport: HTTPTransport,
		}
		return HTTPClient
	}
	return &http.Client{}
}
