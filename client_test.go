package gqlclient

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func mockServer(t *testing.T, expectQuery string, response map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic(err)
		}

		require.Equal(t, expectQuery, string(b))

		err = json.NewEncoder(w).Encode(response)

		if err != nil {
			panic(err)
		}
	}))
}

func TestClient_Send(t *testing.T) {
	query := `query GetUser { user(id: $id) { id name } }`
	variables := map[string]interface{}{
		"id": "1",
	}
	expectQuery := `{"query":"` + query + `","variables":{"id":"1"},"operationName":null}`

	t.Run("struct dest", func(t *testing.T) {
		type structType struct {
			ID   string
			Name string
		}
		var structDest structType

		srv := mockServer(t, expectQuery, map[string]interface{}{
			"data": map[string]interface{}{
				"id":   "1",
				"name": "bob",
			},
		})
		instance := New(srv.URL)

		_, err := instance.Send(context.Background(), &structDest, query, variables)

		require.Equal(t, nil, err)
		require.Equal(t, structType{
			ID:   "1",
			Name: "bob",
		}, structDest)
	})

	t.Run("map dest", func(t *testing.T) {
		var mapDest map[string]interface{}

		srv := mockServer(t, expectQuery, map[string]interface{}{
			"data": map[string]interface{}{
				"id":   "1",
				"name": "bob",
			},
		})
		instance := New(srv.URL)

		_, err := instance.Send(context.Background(), &mapDest, query, variables)

		require.Equal(t, nil, err)
		require.Equal(t, map[string]interface{}{
			"id":   "1",
			"name": "bob",
		}, mapDest)
	})
}

func TestClient_Send_context(t *testing.T) {
	query := `query GetUser { user(id: $id) { id name } }`
	variables := map[string]interface{}{
		"id": "1",
	}
	expectQuery := `{"query":"` + query + `","variables":{"id":"1"},"operationName":null}`

	srv := mockServer(t, expectQuery, map[string]interface{}{
		"data": map[string]interface{}{
			"id":   "1",
			"name": "bob",
		},
	})
	instance := New(srv.URL)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	defer cancel()

	_, err := instance.Raw(ctx, query, variables)

	require.Equal(t, true, os.IsTimeout(errors.Cause(err)))
}

func TestClient_Send_Variations(t *testing.T) {
	query := `query GetUser { user(id: $id) { id name } }`

	type args struct {
		query     string
		variables interface{}
	}

	tests := []struct {
		name    string
		args    *args
		want    *Response
		wantErr bool
	}{
		{
			name: "map variables",
			args: &args{
				query: query,
				variables: map[string]interface{}{
					"id": "1",
				},
			},
			want: &Response{
				Data: map[string]interface{}{
					"id":   "1",
					"name": "bob",
				},
			},
		},
		{
			name: "struct variables",
			args: &args{
				query: query,
				variables: struct {
					ID string `json:"id"`
				}{
					ID: "1",
				},
			},
			want: &Response{
				Data: map[string]interface{}{
					"id":   "1",
					"name": "bob",
				},
			},
		},
		{
			name: "disallow non-object type",
			args: &args{
				query:     query,
				variables: "nope",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			expectQuery := `{"query":"` + query + `","variables":{"id":"1"},"operationName":null}`

			srv := mockServer(t, expectQuery, map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "1",
					"name": "bob",
				},
			})
			instance := New(srv.URL)

			var dest map[string]interface{}
			got, err := instance.Send(context.Background(), &dest, tt.args.query, tt.args.variables)

			if !tt.wantErr {
				require.Equal(t, nil, err)
			}

			if tt.wantErr && err == nil {
				t.Fatalf("want err but got nil. result: %+v", got)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestClient_Errors(t *testing.T) {
	query := `query GetUser { user(id: $id) { id name } }`

	tests := []struct {
		name         string
		mockResponse map[string]interface{}
		want         *Response
	}{
		{
			name: "basic error",
			mockResponse: map[string]interface{}{
				"errors": []map[string]interface{}{{
					"message": "error message",
				}},
			},
			want: &Response{
				Errors: []Error{{
					Message: "error message",
				}},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			expectQuery := `{"query":"` + query + `","variables":{"id":"1"},"operationName":null}`

			srv := mockServer(t, expectQuery, tt.mockResponse)
			instance := New(srv.URL)

			var dest map[string]interface{}
			got, err := instance.Send(context.Background(), &dest, query, map[string]interface{}{
				"id": "1",
			})

			require.Equal(t, nil, err)

			require.Equal(t, tt.want, got)
		})
	}
}
