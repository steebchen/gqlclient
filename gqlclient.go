// Package gqlclient provides a simple GraphQL client with queries, mutations and subscriptions.
package gqlclient

import (
	"net/http"
)

// Client for graphql requests
type GQLClient struct {
	url  string
	http *http.Client
}

// New creates a graphql http
func New(url string) *GQLClient {
	c := &GQLClient{
		url:  url,
		http: http.DefaultClient,
	}

	return c
}

func (c *GQLClient) WithHTTPClient(client *http.Client) *GQLClient {
	c.http = client
	return c
}
