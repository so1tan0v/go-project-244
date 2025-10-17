package plain

import (
	"testing"

	"code/internal/domain/diff"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStylishFormatter_Format(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []diff.DiffNode
		expected string
	}{
		{
			name:     "empty diff",
			nodes:    []diff.DiffNode{},
			expected: ``,
		},
		{
			name: "unchanged",
			nodes: []diff.DiffNode{
				{Key: "a", Type: diff.NodeUnchanged, OldValue: 42},
			},
			expected: ``,
		},
		{
			name: "added",
			nodes: []diff.DiffNode{
				{Key: "b", Type: diff.NodeAdded, NewValue: "new"},
			},
			expected: `Property "b" was added with value: "new"`,
		},
		{
			name: "removed",
			nodes: []diff.DiffNode{
				{Key: "c", Type: diff.NodeRemoved, OldValue: true},
			},
			expected: `Property "c" was removed`,
		},
		{
			name: "updated",
			nodes: []diff.DiffNode{
				{Key: "d", Type: diff.NodeUpdated, OldValue: 1, NewValue: 2},
			},
			expected: `Property "d" was updated. From 1 to 2`,
		},
		{
			name: "nested object",
			nodes: []diff.DiffNode{
				{
					Key:  "parent",
					Type: diff.NodeNested,
					Children: []diff.DiffNode{
						{Key: "child", Type: diff.NodeAdded, NewValue: "value"},
					},
				},
			},
			expected: `Property "parent.child" was added with value: "value"`,
		},
		{
			name: "mixed types with sorting",
			nodes: []diff.DiffNode{
				{Key: "z_removed", Type: diff.NodeRemoved, OldValue: "z"},
				{Key: "a_added", Type: diff.NodeAdded, NewValue: "a"},
				{Key: "m_unchanged", Type: diff.NodeUnchanged, OldValue: nil},
			},
			expected: `Property "a_added" was added with value: "a"
Property "z_removed" was removed`,
		},
		{
			name: "nested object with map value",
			nodes: []diff.DiffNode{
				{
					Key:      "config",
					Type:     diff.NodeAdded,
					NewValue: map[string]any{"port": 8080, "ssl": true},
				},
			},
			expected: `Property "config" was added with value: [complex value]`,
		},
		{
			name: "updated nested object",
			nodes: []diff.DiffNode{
				{
					Key:      "db",
					Type:     diff.NodeUpdated,
					OldValue: map[string]any{"host": "localhost", "port": 5432},
					NewValue: map[string]any{"host": "127.0.0.1", "port": 5433},
				},
			},
			expected: `Property "db" was updated. From [complex value] to [complex value]`,
		},
	}

	formatter := Formatter{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatter.Format(tt.nodes)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
