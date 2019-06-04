package gqlclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Request is the payload for GraphQL queries
type Request struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

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

// MustSend is the same as Send, but panics if an error occurs
func (c *Instance) MustSend(dest interface{}, query string, variables map[string]interface{}) *Response {
	data, err := c.Send(dest, query, variables)

	if err != nil {
		panic(err)
	}

	return data
}

// Send a GraphQL request and unmarshal it to dest
func (c *Instance) Send(dest interface{}, query string, variables map[string]interface{}) (*Response, error) {
	resp, err := c.Raw(query, variables)
	if err != nil {
		return nil, err
	}

	if resp.Errors != nil {
		return resp, nil
	}

	// unpack even if there is an error so we can see partial responses
	unpackErr := unpack(resp.Data, dest)

	return resp, unpackErr
}

// Raw sends a basic GraphQL request with generic types
func (c *Instance) Raw(query string, variables map[string]interface{}) (*Response, error) {
	req := &Request{
		Query:     query,
		Variables: variables,
	}

	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "raw encode")
	}

	rawResponse, err := c.http.Post(c.url, "application/json", bytes.NewBuffer(requestBody))
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
