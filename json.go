package gqlclient

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type rawJsonError struct {
	json.RawMessage
}

func (r *rawJsonError) Error() string {
	return string(r.RawMessage)
}

func unpack(data interface{}, dest interface{}) error {
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      dest,
		TagName:     "json",
		ErrorUnused: true,
		ZeroFields:  true,
	})
	if err != nil {
		return fmt.Errorf("mapstructure: %s", err.Error())
	}

	err = d.Decode(data)

	return err
}
