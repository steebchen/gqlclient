package gqlclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Message string
}

type ResponseData struct {
	Data       interface{}
	Errors     []Error
	Extensions map[string]interface{}
}

func (c *GQLClient) Send(dest interface{}, query string, variables map[string]interface{}) (*ResponseData, error) {
	resp, err := c.RawSend(query, variables)
	if err != nil {
		return resp, err
	}

	// we want to unpack even if there is an error, so we can see partial responses
	unpackErr := unpack(resp.Data, dest)

	if resp.Errors != nil {
		return resp, &rawJsonError{}
	}

	return resp, unpackErr
}

func (c *GQLClient) RawSend(query string, variables map[string]interface{}) (*ResponseData, error) {
	req := &Request{
		Query:     query,
		Variables: variables,
	}

	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encode: %s", err.Error())
	}

	rawResponse, err := c.http.Post(c.url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("post: %s", err.Error())
	}
	defer func() {
		_ = rawResponse.Body.Close()
	}()

	responseBody, err := ioutil.ReadAll(rawResponse.Body)

	if err != nil {
		return nil, fmt.Errorf("read: %s", err.Error())
	}

	if rawResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d: %s", rawResponse.StatusCode, responseBody)
	}

	// decode it into map string first, let mapstructure do the final decode
	// mapstructure is way stricter about unknown fields, can handle embedded structs and more
	respDataRaw := &ResponseData{}
	err = json.Unmarshal(responseBody, &respDataRaw)
	if err != nil {
		return nil, fmt.Errorf("decode: %s", err.Error())
	}

	return respDataRaw, nil
}
