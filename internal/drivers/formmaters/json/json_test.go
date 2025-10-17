package json

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
			expected: `[]`,
		},
		{
			name: "unchanged",
			nodes: []diff.DiffNode{
				{Key: "a", Type: diff.NodeUnchanged, OldValue: 42},
			},
			expected: `[
  {
    "key": "a",
    "oldValue": 42,
    "newValue": null,
    "children": null
  }
]`,
		},
		{
			name: "added",
			nodes: []diff.DiffNode{
				{Key: "b", Type: diff.NodeAdded, NewValue: "new"},
			},
			expected: `[
  {
    "key": "b",
    "oldValue": null,
    "newValue": "new",
    "children": null
  }
]`,
		},
		{
			name: "removed",
			nodes: []diff.DiffNode{
				{Key: "c", Type: diff.NodeRemoved, OldValue: true},
			},
			expected: `[
  {
    "key": "c",
    "oldValue": true,
    "newValue": null,
    "children": null
  }
]`,
		},
		{
			name: "updated",
			nodes: []diff.DiffNode{
				{Key: "d", Type: diff.NodeUpdated, OldValue: 1, NewValue: 2},
			},
			expected: `[
  {
    "key": "d",
    "oldValue": 1,
    "newValue": 2,
    "children": null
  }
]`,
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
			expected: `[
  {
    "key": "parent",
    "oldValue": null,
    "newValue": null,
    "children": [
      {
        "key": "child",
        "oldValue": null,
        "newValue": "value",
        "children": null
      }
    ]
  }
]`,
		},
		{
			name: "mixed types with sorting",
			nodes: []diff.DiffNode{
				{Key: "z_removed", Type: diff.NodeRemoved, OldValue: "z"},
				{Key: "a_added", Type: diff.NodeAdded, NewValue: "a"},
				{Key: "m_unchanged", Type: diff.NodeUnchanged, OldValue: nil},
			},
			expected: `[
  {
    "key": "z_removed",
    "oldValue": "z",
    "newValue": null,
    "children": null
  },
  {
    "key": "a_added",
    "oldValue": null,
    "newValue": "a",
    "children": null
  },
  {
    "key": "m_unchanged",
    "oldValue": null,
    "newValue": null,
    "children": null
  }
]`,
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
			expected: `[
  {
    "key": "config",
    "oldValue": null,
    "newValue": {
      "port": 8080,
      "ssl": true
    },
    "children": null
  }
]`,
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
			expected: `[
  {
    "key": "db",
    "oldValue": {
      "host": "localhost",
      "port": 5432
    },
    "newValue": {
      "host": "127.0.0.1",
      "port": 5433
    },
    "children": null
  }
]`,
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
