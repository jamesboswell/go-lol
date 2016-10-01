// Package lol provides a client for league legends rest api.
package lol

//go:generate go run go-lol-generator/main.go

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

var (
	// ErrNotSupportedRegion is returned if operation is not supported in a region.
	ErrNotSupportedRegion = errors.New("go-lol: operation does not work for such region")
)

// ClientFactory is used to get a http client.
// This must NOT return nil.
type ClientFactory func(context.Context) *http.Client

// Client is a league of legend api fetcher.
type Client struct {
	StaticClient
	apiKey string
}

type StaticClient struct {
	getClient ClientFactory
}

// New creates a new league of legends client.
func New(clientFactory ClientFactory, key string) *Client {
	return &Client{
		StaticClient: NewStatic(clientFactory),
		apiKey:       key,
	}
}

func NewStatic(clientFactory ClientFactory) StaticClient {
	if clientFactory == nil {
		clientFactory = DefaultClientFactory
	}

	return StaticClient{
		getClient: clientFactory,
	}
}

func (c StaticClient) doRequest(ctx context.Context, method, urlStr string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	httpClient := c.getClient(ctx)
	return ctxhttp.Do(ctx, httpClient, req)
}

// HTTPError represents an error returned from riot api server.
type HTTPError struct {
	URL *url.URL
	// Code is the HTTP response status code and will always be populated.
	Code int `json:"code"`
	// Body is the raw response returned by the server.
	// It is often but not always JSON, depending on how the request fails.
	Body []byte
	// Header contains the response header fields from the server.
	Header http.Header
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("go-lol: riot api returned HTTP error %d url: %s", e.Code, e.URL)
}

// verifyResponse returns nil if no error found.
func verifyResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	err := HTTPError{
		Code:   resp.StatusCode,
		URL:    resp.Request.URL,
		Header: resp.Header,
	}

	if resp.Body != nil {
		err.Body, _ = ioutil.ReadAll(resp.Body)
	}
	return err
}
