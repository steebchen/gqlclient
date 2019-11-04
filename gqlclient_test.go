package gqlclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	query := `query GetUser { user(id: $id) { id name } }`
	expectQuery := `{"query":"` + query + `","variables":{"id":"1"},"operationName":null}`

	srv := mockServer(t, expectQuery, map[string]interface{}{
		"data": map[string]interface{}{
			"id":   "1",
			"name": "bob",
		},
	})

	client := New(srv.URL)

	var data struct {
		ID   string
		Name string
	}

	resp, err := client.Send(context.Background(), &data, query, map[string]interface{}{
		"id": "1",
	})

	require.Equal(t, nil, err)

	require.Equal(t, []Error(nil), resp.Errors)
	require.Equal(t, "bob", data.Name)
}

func TestClient_WithHTTPClient(t *testing.T) {
	query := `query GetUser { user(id: $id) { id name } }`
	expectQuery := `{"query":"` + query + `","variables":{"id":"1"},"operationName":null}`

	srv := mockServer(t, expectQuery, map[string]interface{}{
		"data": map[string]interface{}{
			"id":   "1",
			"name": "bob",
		},
	})

	client := New(srv.URL)

	client.WithHTTPClient(http.DefaultClient)

	var data struct {
		ID   string
		Name string
	}

	_, err := client.Send(context.Background(), &data, query, map[string]interface{}{
		"id": "1",
	})

	require.Equal(t, nil, err)
}
