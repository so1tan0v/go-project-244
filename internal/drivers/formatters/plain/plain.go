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

	if err := f.writeNodes(&sb, nodes, ""); err != nil {
		return "", err
	}

	return strings.TrimRight(sb.String(), "\n"), nil
}

func (f Formatter) writeNodes(sb *strings.Builder, nodes []diff.DiffNode, key string) error {
	sorted := make([]diff.DiffNode, len(nodes))
	copy(sorted, nodes)

	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	for i, n := range sorted {
		beforeLen := sb.Len()
		keyValue := getKey(key, n.Key)

		if err := f.parseKey(&n, sb, keyValue); err != nil {
			return err
		}

		if sb.Len() > beforeLen && i != len(sorted)-1 {
			if _, err := fmt.Fprintf(sb, "\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f Formatter) parseKey(n *diff.DiffNode, sb *strings.Builder, keyValue string) error {
	switch n.Type {
	case diff.NodeNested:
		return f.writeNodes(sb, n.Children, keyValue+".")
	case diff.NodeUnchanged:
		// continue
	case diff.NodeRemoved:
		_, err := fmt.Fprintf(sb, "Property '%s' was removed", keyValue)
		if err != nil {
			return err
		}
	case diff.NodeAdded:
		_, err := fmt.Fprintf(sb, "Property '%s' was added with value: %s", keyValue, f.stringify(n.NewValue))
		if err != nil {
			return err
		}
	case diff.NodeUpdated:
		_, err := fmt.Fprintf(sb, "Property '%s' was updated. From %s to %s", keyValue, f.stringify(n.OldValue), f.stringify(n.NewValue))
		if err != nil {
			return err
		}
	}

	return nil
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

func getKey(key, i string) string {
	return key + i
}
