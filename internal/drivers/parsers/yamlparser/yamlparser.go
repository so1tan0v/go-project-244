package yamlparser

import (
	"gopkg.in/yaml.v3"
)

type Parser struct{}

func (p Parser) Parse(data []byte) (map[string]any, error) {
	var m map[string]any

	if err := yaml.Unmarshal(data, &m); err != nil {
		return map[string]any{}, err
	}

	return m, nil
}
