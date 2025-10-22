package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsObjectNil(t *testing.T) {
	m := asObject(nil)
	assert.NotNil(t, m)
	assert.Len(t, m, 0)
}

func TestAsObjectMapStringAny(t *testing.T) {
	src := map[string]any{"a": 1}
	m := asObject(src)
	assert.Equal(t, src, m)
	// ensure clone (modifying result doesn't change original)
	m["a"] = 2
	assert.Equal(t, 1, src["a"])
}

func TestAsObjectMapStringInterface(t *testing.T) {
	src := map[string]interface{}{"b": true}
	m := asObject(src)
	assert.Equal(t, map[string]any{"b": true}, m)
	m["b"] = false
	assert.Equal(t, true, src["b"])
}

func TestAsObjectOtherType(t *testing.T) {
	m := asObject(123)
	assert.NotNil(t, m)
	assert.Len(t, m, 0)
}
