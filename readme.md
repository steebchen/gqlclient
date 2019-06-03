# gqlclient [![GoDoc](https://godoc.org/github.com/steebchen/gqlclient?status.png)](http://godoc.org/github.com/steebchen/gqlclient)

Package gqlclient provides a GraphQL client implementation.

- Simple, familiar API
- Use strong Go types for response and variables
- Simple error handling

Coming soon:

- Uploads
- Subscriptions

## Installation

Make sure you have a working Go environment, preferably with Go modules.

To install graphql, simply run:

```
$ go get github.com/steebchen/gqlclient
```

## Usage

```go
package main

import (
	"github.com/steebchen/gqlclient"
)

func main() {
	client := gqlclient.New("https://metaphysics-production.artsy.net")
	
	var resp struct {
		Article struct {
			ID    string
			Title string
		}
	}
	
	_, err := client.Send(&resp, `
		query Article($id: String!) {
			article(id: $id) {
				id
				title
			}
		}
	`, map[string]interface{}{
		"id": "55bfed9275de7b060098b9bc",
	})
}
```
