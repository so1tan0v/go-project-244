package plain

import (
	"code/internal/domain/diff"
	"fmt"
	"sort"
	"strings"
)

type Formatter struct{}

func (f Formatter) Format(nodes []diff.DiffNode) (string, error) {
	var sb strings.Builder

	f.writeNodes(&sb, nodes, "")

	return strings.TrimRight(sb.String(), "\n"), nil
}

func (f Formatter) writeNodes(sb *strings.Builder, nodes []diff.DiffNode, key string) {
	sorted := make([]diff.DiffNode, len(nodes))

	copy(sorted, nodes)

	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	getKey := func(i string) string {
		return key + i
	}

	for i, n := range sorted {
		beforeLen := sb.Len()
		keyValue := getKey(n.Key)

		switch n.Type {
		case diff.NodeNested:
			f.writeNodes(sb, n.Children, keyValue+".")
		case diff.NodeUnchanged:
			continue
		case diff.NodeRemoved:
			_, err := fmt.Fprintf(sb, "Property '%s' was removed", keyValue)
			if err != nil {
				return
			}
		case diff.NodeAdded:
			_, err := fmt.Fprintf(sb, "Property '%s' was added with value: %s", keyValue, f.stringify(n.NewValue))
			if err != nil {
				return
			}
		case diff.NodeUpdated:
			_, err := fmt.Fprintf(sb, "Property '%s' was updated. From %s to %s", keyValue, f.stringify(n.OldValue), f.stringify(n.NewValue))
			if err != nil {
				return
			}
		}

		wroteSomething := sb.Len() > beforeLen
		if wroteSomething && i != len(sorted)-1 {
			_, err := fmt.Fprintf(sb, "\n")
			if err != nil {
				return
			}
		}
	}
}

func (f Formatter) stringify(v any) string {
	switch m := v.(type) {
	case map[string]any:
		return "[complex value]"
	case string:
		return fmt.Sprintf("'%s'", m)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", m)
	}
}
