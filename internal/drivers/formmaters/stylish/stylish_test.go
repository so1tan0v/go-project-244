package stylish

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
			name:  "empty diff",
			nodes: []diff.DiffNode{},
			expected: `{

}`,
		},
		{
			name: "unchanged",
			nodes: []diff.DiffNode{
				{Key: "a", Type: diff.NodeUnchanged, OldValue: 42},
			},
			expected: `{
    a: 42
}`,
		},
		{
			name: "added",
			nodes: []diff.DiffNode{
				{Key: "b", Type: diff.NodeAdded, NewValue: "new"},
			},
			expected: `{
  + b: "new"
}`,
		},
		{
			name: "removed",
			nodes: []diff.DiffNode{
				{Key: "c", Type: diff.NodeRemoved, OldValue: true},
			},
			expected: `{
  - c: true
}`,
		},
		{
			name: "updated",
			nodes: []diff.DiffNode{
				{Key: "d", Type: diff.NodeUpdated, OldValue: 1, NewValue: 2},
			},
			expected: `{
  - d: 1
  + d: 2
}`,
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
			expected: `{
    parent: {
      + child: "value"
    }
}`,
		},
		{
			name: "mixed types with sorting",
			nodes: []diff.DiffNode{
				{Key: "z_removed", Type: diff.NodeRemoved, OldValue: "z"},
				{Key: "a_added", Type: diff.NodeAdded, NewValue: "a"},
				{Key: "m_unchanged", Type: diff.NodeUnchanged, OldValue: nil},
			},
			expected: `{
  + a_added: "a"
    m_unchanged: null
  - z_removed: "z"
}`,
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
			expected: `{
  + config: {
      port: 8080
      ssl: true
    }
}`,
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
			expected: `{
  - db: {
      host: "localhost"
      port: 5432
    }
  + db: {
      host: "127.0.0.1"
      port: 5433
    }
}`,
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
