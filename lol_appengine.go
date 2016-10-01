// +build appengine

package lol

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

var (
	// DefaultClientFactory is a ClientProviderFunc which use urlfetch.
	DefaultClientFactory ClientFactory = func(ctx context.Context) *http.Client {
		return urlfetch.Client(ctx)
	}
)
