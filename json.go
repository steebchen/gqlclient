package gqlclient

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

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
