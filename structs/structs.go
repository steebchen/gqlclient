package structs

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Unpack a map `origin` to a struct. `dest` must be a pointer to a struct.
func Unpack(origin interface{}, dest interface{}) error {
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      dest,
		TagName:     "json",
		ErrorUnused: true,
		ZeroFields:  true,
	})
	if err != nil {
		return errors.Wrap(err, "mapstructure")
	}

	err = d.Decode(origin)

	return err
}

// StructToMap converts a struct to a map[string]interface{}
func StructToMap(in interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, errors.Wrap(err, "variables json marshaling failed")
	}

	var out map[string]interface{}
	err = json.Unmarshal(data, &out)
	if err != nil {
		return nil, errors.Wrap(err, "variables unmarshaling failed")
	}

	return out, nil
}
