package jsonparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse_Success(t *testing.T) {
	p := Parser{}

	jsonString := `{"name": "test"}`
	jsonBytes := []byte(jsonString)

	result, err := p.Parse(jsonBytes)
	assert.Nil(t, err)
	assert.Equal(t, "test", result["name"])
}

func TestParser_Parse_Error(t *testing.T) {
	p := Parser{}

	jsonString := `name": "test`
	jsonBytes := []byte(jsonString)

	result, err := p.Parse(jsonBytes)
	assert.NotNil(t, err)

	assert.Equal(t, map[string]any{}, result)
}
