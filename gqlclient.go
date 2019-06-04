// Package gqlclient provides a simple GraphQL client with queries, mutations and subscriptions.
package gqlclient

import (
	"net/http"
)

// Instance is gqlclient struct returned with New()
type Instance struct {
	url  string
	http *http.Client
}

// New creates a graphql http
func New(url string) *Instance {
	c := &Instance{
		url:  url,
		http: http.DefaultClient,
	}

	return c
}

// WithHTTPClient uses a given http client for all requests
func (c *Instance) WithHTTPClient(client *http.Client) *Instance {
	c.http = client
	return c
}
