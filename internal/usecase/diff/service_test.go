package diff

import (
	"code/internal/drivers/jsonparser"
	"code/internal/drivers/stylishformatter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	formatter := stylishformatter.Formatter{}
	parser := jsonparser.Parser{}

	sn := NewService(parser, formatter)
	ss := Service{parser: parser, formatter: formatter}

	assert.Equal(t, sn.parser, ss.parser)
}

func TestService_GenerateDiff_Success(t *testing.T) {
	formatter := stylishformatter.Formatter{}
	parser := jsonparser.Parser{}

	sn := NewService(parser, formatter)

	json1 := []byte(`{"test": "value"}`)
	json2 := []byte(`{"test": "value2"}`)

	result, err := sn.GenerateDiff(json1, json2)

	assert.Nil(t, err)
	assert.Equal(t, `{
  - test: "value"
  + test: "value2"
}`, result)
}

func TestService_GenerateDiff_ParseError1(t *testing.T) {
	formatter := stylishformatter.Formatter{}
	parser := jsonparser.Parser{}

	sn := NewService(parser, formatter)

	json1 := []byte(`test": "value"}`)
	json2 := []byte(`{"test": "value2"}`)

	result, err := sn.GenerateDiff(json1, json2)

	assert.Equal(t, "", result)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid character 'e' in literal true (expecting 'r')", err.Error())
}

func TestService_GenerateDiff_ParseError2(t *testing.T) {
	formatter := stylishformatter.Formatter{}
	parser := jsonparser.Parser{}

	sn := NewService(parser, formatter)

	json1 := []byte(`{"test": "value"}`)
	json2 := []byte(`test": "value2"}`)

	result, err := sn.GenerateDiff(json1, json2)

	assert.Equal(t, "", result)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid character 'e' in literal true (expecting 'r')", err.Error())
}
