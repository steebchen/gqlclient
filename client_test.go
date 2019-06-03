package gqlclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	h := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic(err)
		}

		expect := `{"query":"query GetUser { user(id: $id) { id name } }","variables":{"id":"1"}}`
		require.Equal(t, expect, string(b))

		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"id":   "1",
				"name": "bob",
			},
		})

		if err != nil {
			panic(err)
		}
	}))

	client := New(h.URL).WithHTTPClient(http.DefaultClient)

	var data struct {
		ID   string
		Name string
	}

	resp, err := client.Send(&data, `query GetUser { user(id: $id) { id name } }`, map[string]interface{}{
		"id": "1",
	})

	require.NoError(t, err)

	require.Equal(t, []Error(nil), resp.Errors)
	require.Equal(t, "bob", data.Name)
}
