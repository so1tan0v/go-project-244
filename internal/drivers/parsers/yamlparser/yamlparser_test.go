package yamlparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse_Success(t *testing.T) {
	p := Parser{}

	yamlString := `name: test`
	jsonBytes := []byte(yamlString)

	result, err := p.Parse(jsonBytes)
	assert.Nil(t, err)
	assert.Equal(t, "test", result["name"])
}

func TestParser_Parse_Error(t *testing.T) {
	p := Parser{}

	yamlString := `
name""""": "test
`
	jsonBytes := []byte(yamlString)

	result, err := p.Parse(jsonBytes)
	assert.NotNil(t, err)

	assert.Equal(t, map[string]any{}, result)
}
