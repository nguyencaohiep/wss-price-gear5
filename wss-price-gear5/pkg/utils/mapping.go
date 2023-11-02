package utils

import (
	"bytes"
	"encoding/json"
)

// Mapping from Map[string]any to any/struct{}
func Mapping(in, out any) error {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(in)
	if err != nil {
		return err
	}
	err = json.NewDecoder(buf).Decode(out)
	if err != nil {
		return err
	}
	return nil
}
