package main

import (
	"context"
	"crypto/tls"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"time"
)

type ApiClient struct {
	client      *http.Client
	rateLimiter *rate.Limiter
}

func NewClient(apiRateLimit *rate.Limiter) *ApiClient {
	return &ApiClient{
		client: &http.Client{
			Timeout: Timeout * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		rateLimiter: apiRateLimit,
	}
}

func (c *ApiClient) requestWithClient(req *http.Request, httpClient *http.Client) (*http.Response, error) {
	ctx := context.Background()
	err := c.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *ApiClient) fetchPage(proxy string, targetUrl string, country *string) (*http.Response, error) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}

	transport := http.Transport{
		Proxy:           http.ProxyURL(proxyUrl),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Transport: &transport,
		Timeout:   Timeout * time.Second,
	}
	apiRequest, _ := http.NewRequest("GET", targetUrl, nil)

	browserHeaders := getRandomBrowserHeaders(country)
	for headerKey, headerValue := range browserHeaders {
		apiRequest.Header.Set(headerKey, headerValue.(string))
	}

	return c.requestWithClient(apiRequest, httpClient)
}
