// Package gqlclient provides a simple GraphQL client with queries, mutations and subscriptions.
package gqlclient

import (
	"net/http"
)

// Client is the GraphQL client which is returned by New()
type Client struct {
	url  string
	http *http.Client
}

// New creates a graphql http
func New(url string) *Client {
	c := &Client{
		url:  url,
		http: http.DefaultClient,
	}

	return c
}

// WithHTTPClient uses a given http client for all requests
func (c *Client) WithHTTPClient(client *http.Client) *Client {
	c.http = client
	return c
}
