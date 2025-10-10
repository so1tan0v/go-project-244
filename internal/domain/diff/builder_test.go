package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDiff(t *testing.T) {
	left := map[string]interface{}{}
	right := map[string]interface{}{}

	left["key_string"] = "value"
	right["key_string"] = "value1"

	left["key_nil"] = nil
	right["key_nil"] = 1

	left["key_int"] = 2
	right["key_int"] = 1

	left["key_bool"] = true
	right["key_bool"] = false

	left["key_equal"] = true
	right["key_equal"] = true

	left["key_left_new"] = 1.2

	right["key_right_new"] = 1.2

	result := BuildDiff(left, right)

	for _, d := range result {
		assert.Equal(t, d.OldValue, left[d.Key])
		assert.Equal(t, d.NewValue, right[d.Key])

		leftVal, leftOk := left[d.Key]
		rightVal, rightOk := right[d.Key]

		switch {
		case leftOk && !rightOk:
			assert.Equal(t, d.Type, NodeRemoved)
		case !leftOk && rightOk:
			assert.Equal(t, d.Type, NodeAdded)
		default:
			if isObject(leftVal) && isObject(rightVal) {
				assert.Equal(t, d.Type, NodeNested)
				continue
			}

			if valuesEqual(leftVal, rightVal) {
				assert.Equal(t, d.Type, NodeUnchanged)
				continue
			}

			assert.Equal(t, d.Type, NodeUpdated)
		}
	}
}
