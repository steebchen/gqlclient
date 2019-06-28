package gqlclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	srv := mockServer(t)

	client := New(srv.URL)

	var data struct {
		ID   string
		Name string
	}

	resp, err := client.Send(context.Background(), &data, `query GetUser { user(id: $id) { id name } }`, map[string]interface{}{
		"id": "1",
	})

	require.NoError(t, err)

	require.Equal(t, []Error(nil), resp.Errors)
	require.Equal(t, "bob", data.Name)
}

func TestClient_WithHTTPClient(t *testing.T) {
	srv := mockServer(t)

	client := New(srv.URL)

	client.WithHTTPClient(http.DefaultClient)

	var data struct {
		ID   string
		Name string
	}

	_, err := client.Send(context.Background(), &data, `query GetUser { user(id: $id) { id name } }`, map[string]interface{}{
		"id": "1",
	})

	require.NoError(t, err)
}
