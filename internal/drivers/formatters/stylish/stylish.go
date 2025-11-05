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
		f.parseValue(&n, sb, indent, depth)

		if i < len(sorted)-1 {
			_, err := fmt.Fprint(sb, "\n")
			if err != nil {
				return
			}
		}
	}
}

func (f Formatter) parseValue(n *diff.DiffNode, sb *strings.Builder, indent string, depth int) {
	switch n.Type {
	case diff.NodeNested:
		f.parseNodeNested(sb, n, indent, depth)
	case diff.NodeUnchanged:
		f.parseNodeUnchanged(sb, n, indent, depth)
	case diff.NodeRemoved:
		f.parseNodeRemoved(sb, n, indent, depth)
	case diff.NodeAdded:
		f.parseNodeAdded(sb, n, indent, depth)
	case diff.NodeUpdated:
		f.parseNodeUpdated(sb, n, indent, depth)
	}
}

func (f Formatter) stringify(v any, depth int) string {
	switch m := v.(type) {
	case map[string]any:
		return f.stringifyObject(m, depth)
	case string:
		return m
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", m)
	}
}

func (f Formatter) stringifyObject(m map[string]any, depth int) string {
	var sb strings.Builder

	if _, err := fmt.Fprint(&sb, "{\n"); err != nil {
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

		if _, err := fmt.Fprintf(&sb, "%s%s: %s", indent, k, f.stringify(val, depth+2)); err != nil {
			return ""
		}

		if i < len(keys)-1 {
			if _, err := fmt.Fprintf(&sb, "\n"); err != nil {
				return ""
			}
		}
	}

	if _, err := fmt.Fprint(&sb, "\n"+strings.Repeat(" ", (depth+1)*2)+"}"); err != nil {
		return ""
	}

	return sb.String()
}
