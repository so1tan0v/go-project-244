package stylish

import (
	"code/internal/domain/diff"
	"fmt"
	"sort"
	"strings"
)

type Formatter struct{}

func (f Formatter) Format(nodes []diff.DiffNode) (string, error) {
	var sb strings.Builder

	_, err := fmt.Fprintf(&sb, "{\n")
	if err != nil {
		return "", err
	}

	f.writeNodes(&sb, nodes, 1)

	_, err = fmt.Fprintf(&sb, "\n}")
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}

func (f Formatter) writeNodes(sb *strings.Builder, nodes []diff.DiffNode, depth int) {
	indent := strings.Repeat(" ", depth*2)

	sorted := make([]diff.DiffNode, len(nodes))

	copy(sorted, nodes)

	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	for i, n := range sorted {
		switch n.Type {
		case diff.NodeNested:
			_, err := fmt.Fprintf(sb, "%s  %s: {\n", indent, n.Key)
			if err != nil {
				return
			}

			f.writeNodes(sb, n.Children, depth+2)

			_, err = fmt.Fprint(sb, "\n"+indent+"  }")
			if err != nil {
				return
			}
		case diff.NodeUnchanged:
			_, err := fmt.Fprintf(sb, "%s  %s: %s", indent, n.Key, f.stringify(n.OldValue, depth))
			if err != nil {
				return
			}
		case diff.NodeRemoved:
			_, err := fmt.Fprintf(sb, "%s- %s: %s", indent, n.Key, f.stringify(n.OldValue, depth))
			if err != nil {
				return
			}
		case diff.NodeAdded:
			_, err := fmt.Fprintf(sb, "%s+ %s: %s", indent, n.Key, f.stringify(n.NewValue, depth))
			if err != nil {
				return
			}
		case diff.NodeUpdated:
			_, err := fmt.Fprintf(sb, "%s- %s: %s\n", indent, n.Key, f.stringify(n.OldValue, depth))
			if err != nil {
				return
			}

			_, err = fmt.Fprintf(sb, "%s+ %s: %s", indent, n.Key, f.stringify(n.NewValue, depth))
			if err != nil {
				return
			}
		}

		if i < len(sorted)-1 {
			_, err := fmt.Fprint(sb, "\n")
			if err != nil {
				return
			}
		}
	}
}

func (f Formatter) stringify(v any, depth int) string {
	switch m := v.(type) {
	case map[string]any:
		var sb strings.Builder
		_, err := fmt.Fprint(&sb, "{\n")
		if err != nil {
			return ""
		}

		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		for i, k := range keys {
			val := m[k]
			indent := strings.Repeat(" ", (depth+3)*2)

			_, err := fmt.Fprintf(&sb, "%s%s: %s", indent, k, f.stringify(val, depth+2))
			if err != nil {
				return ""
			}

			if i < len(keys)-1 {
				_, err2 := fmt.Fprintf(&sb, "\n")
				if err2 != nil {
					return ""
				}
			}
		}

		_, err = fmt.Fprint(&sb, "\n"+strings.Repeat(" ", (depth+1)*2)+"}")
		if err != nil {
			return ""
		}

		return sb.String()
	case string:
		return m
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", m)
	}
}
