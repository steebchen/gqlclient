package gqlclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Request struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

func (c *GQLClient) MustSend(dest interface{}, query string, variables map[string]interface{}) *ResponseData {
	data, err := c.Send(dest, query, variables)

	if err != nil {
		panic(err)
	}

	return data
}

type Error struct {
	Message    string
	Path       []string
	Extensions map[string]interface{}
}

type ResponseData struct {
	Data       interface{}
	Errors     []Error
	Extensions map[string]interface{}
}

func (c *GQLClient) Send(dest interface{}, query string, variables map[string]interface{}) (*ResponseData, error) {
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

func (c *GQLClient) Raw(query string, variables map[string]interface{}) (*ResponseData, error) {
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
	respDataRaw := &ResponseData{}
	err = json.Unmarshal(responseBody, &respDataRaw)
	if err != nil {
		return nil, errors.Wrap(err, "raw decode")
	}

	return respDataRaw, nil
}
