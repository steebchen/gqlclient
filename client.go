package gqlclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/steebchen/gqlclient/structs"
)

// Error is a GraphQL Error
type Error struct {
	Message    string
	Path       []string
	Extensions map[string]interface{}
}

// Response is the payload for a GraphQL response
type Response struct {
	Data       interface{}
	Errors     []Error
	Extensions map[string]interface{}
}

// request is the payload for GraphQL queries
type gqlRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

// Send a GraphQL request to struct or map
func (c *Client) Send(ctx context.Context, dest interface{}, query string, variables interface{}) (*Response, error) {
	unboxedVars, err := structs.StructToMap(variables)
	if err != nil {
		return nil, errors.Wrap(err, "StructToMap failed")
	}

	resp, err := c.Raw(ctx, query, unboxedVars)
	if err != nil {
		return nil, err
	}

	if resp.Errors != nil {
		return resp, nil
	}

	// unpack even if there is an error so we can see partial responses
	unpackErr := structs.Unpack(resp.Data, dest)

	return resp, unpackErr
}

// Raw sends a basic GraphQL request without any struct types.
// Parameter `variables` can be either a map or a struct
func (c *Client) Raw(ctx context.Context, query string, variables map[string]interface{}) (*Response, error) {
	var err error

	payload := &gqlRequest{
		Query:     query,
		Variables: variables,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "raw encode")
	}

	req, err := http.NewRequest("post", c.url, bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")
	req = req.WithContext(ctx)

	rawResponse, err := c.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "raw post")
	}
	defer func() {
		_ = rawResponse.Body.Close()
	}()

	responseBody, err := ioutil.ReadAll(rawResponse.Body)

	if err != nil {
		return nil, errors.Wrap(err, "raw read")
	}

	if rawResponse.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status code %d with response %s", rawResponse.StatusCode, responseBody)
	}

	// decode it into map string first, let mapstructure do the final decode
	// mapstructure is way stricter about unknown fields, can handle embedded structs and more
	respDataRaw := &Response{}
	err = json.Unmarshal(responseBody, &respDataRaw)
	if err != nil {
		return nil, errors.Wrap(err, "raw decode")
	}

	return respDataRaw, nil
}
