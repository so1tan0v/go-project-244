package diff

import (
	"maps"
	"reflect"
	"slices"
)

func BuildDiff(left, right map[string]any) []DiffNode {
	allKeysMap := map[string]struct{}{}

	for k := range left {
		allKeysMap[k] = struct{}{}
	}

	for k := range right {
		allKeysMap[k] = struct{}{}
	}

	allKeys := make([]string, 0, len(allKeysMap))
	for k := range allKeysMap {
		allKeys = append(allKeys, k)
	}

	slices.Sort(allKeys)

	result := make([]DiffNode, 0, len(allKeys))
	for _, key := range allKeys {
		leftVal, leftOk := left[key]
		rightVal, rightOk := right[key]

		switch {
		case leftOk && !rightOk:
			result = append(result, DiffNode{Key: key, Type: NodeRemoved, OldValue: leftVal})
		case !leftOk && rightOk:
			result = append(result, DiffNode{Key: key, Type: NodeAdded, NewValue: rightVal})
		default:
			if isObject(leftVal) && isObject(rightVal) {
				ln := asObject(leftVal)
				rn := asObject(rightVal)
				children := BuildDiff(ln, rn)

				result = append(result, DiffNode{Key: key, Type: NodeNested, Children: children})

				continue
			}

			if valuesEqual(leftVal, rightVal) {
				result = append(result, DiffNode{Key: key, Type: NodeUnchanged, OldValue: leftVal, NewValue: rightVal})

				continue
			}

			result = append(result, DiffNode{Key: key, Type: NodeUpdated, OldValue: leftVal, NewValue: rightVal})
		}
	}

	return result
}

func isObject(v any) bool {
	if v == nil {
		return false
	}

	_, ok := v.(map[string]any)
	if ok {
		return true
	}

	_, ok = v.(map[string]interface{})

	return ok
}

func asObject(v any) map[string]any {
	if v == nil {
		return map[string]any{}
	}

	if m, ok := v.(map[string]any); ok {
		return maps.Clone(m)
	}

	if m, ok := v.(map[string]interface{}); ok {
		return maps.Clone(m)
	}

	return map[string]any{}
}

func valuesEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}
