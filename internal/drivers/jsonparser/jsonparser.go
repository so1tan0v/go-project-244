package jsonparser

import (
	"encoding/json"
)

type Parser struct{}

func (p Parser) Parse(data []byte) (map[string]any, error) {
	var m map[string]any

	if err := json.Unmarshal(data, &m); err != nil {
		return map[string]any{}, err
	}

	return m, nil
}

func (p Parser) FormatName() string {
	return "json"
}
