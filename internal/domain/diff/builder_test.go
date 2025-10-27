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
			testDefault(t, leftVal, rightVal, d.Type)
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
			leftVal, leftOk := getValueFromObjectInObject(t, left, "key_object", d1.Key)
			rightVal, rightOk := getValueFromObjectInObject(t, right, "key_object", d1.Key)

			switch {
			case leftOk && !rightOk:
				assert.Equal(t, d1.Type, NodeRemoved)
			case !leftOk && rightOk:
				assert.Equal(t, d1.Type, NodeAdded)
			default:
				testDefault(t, leftVal, rightVal, d1.Type)
			}
		}
	}
}

func testDefault(t *testing.T, leftVal, rightVal any, typeNode NodeType) {
	t.Helper()

	if isObject(leftVal) && isObject(rightVal) {
		assert.Equal(t, typeNode, NodeNested)

		return
	}

	if valuesEqual(leftVal, rightVal) {
		assert.Equal(t, typeNode, NodeUnchanged)

		return
	}

	assert.Equal(t, typeNode, NodeUpdated)
}

func getValueFromObjectInObject(t *testing.T, obj map[string]any, key string, subKey string) (any, bool) {
	t.Helper()

	objKeyVal, objKeyOk := obj[key]
	if !objKeyOk {
		t.Error("expected key")
	}

	objMap, ok := objKeyVal.(map[string]any)
	if !ok {
		t.Fatalf("obj['%s'] is not a map[string]any", key)
	}

	objVal, objOk := objMap[subKey]

	return objVal, objOk
}
