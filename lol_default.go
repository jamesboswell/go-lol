// +build !appengine

package lol

import (
	"net/http"

	"golang.org/x/net/context"
)

var (
	// DefaultClientFactory is a simple ClientProviderFunc which returns http.DefaultClient
	DefaultClientFactory ClientFactory = func(context.Context) *http.Client {
		return http.DefaultClient
	}
)
