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

func TestBuildDiffObjectInObject(t *testing.T) {
	left := map[string]any{}
	right := map[string]any{}

	left["key_object"] = map[string]any{
		"key_string": "value",
		"key_nil":    nil,
		"key_int":    1,
		"key_bool":   true,
		"key_equal":  true,
	}

	right["key_object"] = map[string]any{
		"key_string": "value1",
		"key_nil":    1,
		"key_int":    2,
		"key_bool":   true,
		"key_equal":  true,
	}

	result := BuildDiff(left, right)

	for _, d := range result {
		assert.Len(t, d.Children, 5)
		assert.Equal(t, d.Key, "key_object")
		assert.Equal(t, d.Type, NodeNested)

		for _, d1 := range d.Children {
			leftKeyVal, leftKeyOk := left["key_object"]
			rightKeyVal, rightKeyOk := right["key_object"]

			if !leftKeyOk || !rightKeyOk {
				t.Error("expected key_object")
			}

			// Приводим к map[string]any
			leftMap, ok := leftKeyVal.(map[string]any)
			if !ok {
				t.Fatalf("left['key_object'] is not a map[string]any")
			}
			rightMap, ok := rightKeyVal.(map[string]any)
			if !ok {
				t.Fatalf("right['key_object'] is not a map[string]any")
			}

			leftVal, leftOk := leftMap[d1.Key]
			rightVal, rightOk := rightMap[d1.Key]

			switch {
			case leftOk && !rightOk:
				assert.Equal(t, d1.Type, NodeRemoved)
			case !leftOk && rightOk:
				assert.Equal(t, d1.Type, NodeAdded)
			default:
				if isObject(leftVal) && isObject(rightVal) {
					assert.Equal(t, d1.Type, NodeNested)
					continue
				}

				if valuesEqual(leftVal, rightVal) {
					assert.Equal(t, d1.Type, NodeUnchanged)
					continue
				}

				assert.Equal(t, d1.Type, NodeUpdated)
			}
		}

	}
}
