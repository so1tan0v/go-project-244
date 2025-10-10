package stylishformatter

import (
	"code/internal/domain/diff"
	"fmt"
	"sort"
	"strings"
)

type Formatter struct{}

func (f Formatter) Name() string { return "stylish" }

func (f Formatter) Format(nodes []diff.DiffNode) (string, error) {
	var sb strings.Builder

	sb.WriteString("{\n")
	f.writeNodes(&sb, nodes, 1)
	sb.WriteString("\n}")

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
			sb.WriteString(fmt.Sprintf("%s  %s: {\n", indent, n.Key))
			f.writeNodes(sb, n.Children, depth+2)
			sb.WriteString("\n" + indent + "  }")
		case diff.NodeUnchanged:
			sb.WriteString(fmt.Sprintf("%s  %s: %s", indent, n.Key, f.stringify(n.OldValue, depth)))
		case diff.NodeRemoved:
			sb.WriteString(fmt.Sprintf("%s- %s: %s", indent, n.Key, f.stringify(n.OldValue, depth)))
		case diff.NodeAdded:
			sb.WriteString(fmt.Sprintf("%s+ %s: %s", indent, n.Key, f.stringify(n.NewValue, depth)))
		case diff.NodeUpdated:
			sb.WriteString(fmt.Sprintf("%s- %s: %s\n", indent, n.Key, f.stringify(n.OldValue, depth)))
			sb.WriteString(fmt.Sprintf("%s+ %s: %s", indent, n.Key, f.stringify(n.NewValue, depth)))
		}

		if i < len(sorted)-1 {
			sb.WriteString("\n")
		}
	}
}

func (f Formatter) stringify(v any, depth int) string {
	switch m := v.(type) {
	case map[string]any:
		var sb strings.Builder
		sb.WriteString("{\n")

		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		for i, k := range keys {
			val := m[k]
			indent := strings.Repeat(" ", (depth+2)*2)

			sb.WriteString(fmt.Sprintf("%s%s: %s", indent, k, f.stringify(val, depth+2)))

			if i < len(keys)-1 {
				sb.WriteString("\n")
			}
		}
		sb.WriteString("\n" + strings.Repeat(" ", (depth+1)*2) + "}")

		return sb.String()
	case string:
		return fmt.Sprintf("\"%s\"", m)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", m)
	}
}
