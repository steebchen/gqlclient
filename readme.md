# gqlclient [![GoDoc](https://godoc.org/github.com/steebchen/gqlclient?status.png)](http://godoc.org/github.com/steebchen/gqlclient) [![Go Report Card](https://goreportcard.com/badge/github.com/steebchen/gqlclient)](https://goreportcard.com/report/github.com/steebchen/gqlclient) [![Actions Status](https://wdp9fww0r9.execute-api.us-west-2.amazonaws.com/production/badge/steebchen/gqlclient)](https://wdp9fww0r9.execute-api.us-west-2.amazonaws.com/production/results/steebchen/gqlclient)

The package gqlclient provides a GraphQL client implementation.

<p align="center">
	<a target="_blank" href="https://github.com/MariaLetta/free-gophers-pack">
		<img src="./gopher.svg"  alt="GraphQL Gopher" height="350" />
	</a>
</p>

Reasons to use gqlclient:

- Simple, familiar API
- Use strong Go types for variables and response data
- Simple error handling
- Supports GraphQL Errors with Extensions

*Note*: This package already works quite well, but it is under heavy development to work towards a v1.0 release. Before that, the API may have breaking changes even with minor versions. 

Coming soon:

- Uploads
- Subscriptions
- More options (http headers, request context)

## Installation

Make sure you have a working Go environment, preferably with Go modules.

To install graphql, simply run:

```
$ go get github.com/steebchen/gqlclient
```

## Quickstart

The recommended way is to use structs depending on your schema for best type-safety.

```go
package main

import (
	"log"
	"github.com/steebchen/gqlclient"
)

func main() {
	client := gqlclient.New("https://metaphysics-production.artsy.net")
	
	var data struct {
		Article struct {
			ID  string
			Title string
		}
	}
	
	type variables struct{
		ID string `json:"id"`
	}

	query := `
		query Article($id: String!) {
			article(id: $id) {
				id
				title
			}
		}
	`

	_, err := client.Send(&data, query, variables{
		ID: "55bfed9275de7b060098b9bc",
	})

	if err != nil {
		panic(err)
	}

	log.Printf("data: %+v", data)
	// Output:
	// Article: {
	//   ID: 55bfed9275de7b060098b9bc
	//   Title: How the 1960s’ Most Iconic Artists Made Art Contemporary
	// }
}
```

If you don't want to use structs, you use `Raw()` to use maps for both input (variables) and output (response data).

```go
resp, err := client.Raw(query, map[string]interface{}{
  "id": "55bfed9275de7b060098b9bc",
})

if err != nil {
	panic(err)
}

log.Printf("data: %+v", resp.Data)
// Output:
// data: map[
//   article: map[
//	   id: 55bfed9275de7b060098b9bc
//	   title: How the 1960s’ Most Iconic Artists Made Art Contemporary
//   ]
// ]
```

Both `Send()` and `Raw()` always return a GraphQL [`Response`](https://godoc.org/github.com/steebchen/gqlclient#Response), so you can access GraphQL Errors and Extensions.
