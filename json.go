package gqlclient

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
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
		return errors.Wrap(err, "mapstructure")
	}

	err = d.Decode(data)

	return err
}
