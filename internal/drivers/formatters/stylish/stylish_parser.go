package stylish

import (
	"code/internal/domain/diff"
	"fmt"
	"strings"
)

func (f Formatter) parseNodeNested(sb *strings.Builder, n *diff.DiffNode, indent string, depth int) {
	if _, err := fmt.Fprintf(sb, "%s  %s: {\n", indent, n.Key); err != nil {
		return
	}

	f.writeNodes(sb, n.Children, depth+2)

	if _, err := fmt.Fprint(sb, "\n"+indent+"  }"); err != nil {
		return
	}
}

func (f Formatter) parseNodeUnchanged(sb *strings.Builder, n *diff.DiffNode, indent string, depth int) {
	if _, err := fmt.Fprintf(sb, "%s  %s: %s", indent, n.Key, f.stringify(n.OldValue, depth)); err != nil {
		return
	}
}

func (f Formatter) parseNodeRemoved(sb *strings.Builder, n *diff.DiffNode, indent string, depth int) {
	if _, err := fmt.Fprintf(sb, "%s- %s: %s", indent, n.Key, f.stringify(n.OldValue, depth)); err != nil {
		return
	}
}

func (f Formatter) parseNodeAdded(sb *strings.Builder, n *diff.DiffNode, indent string, depth int) {
	if _, err := fmt.Fprintf(sb, "%s+ %s: %s", indent, n.Key, f.stringify(n.NewValue, depth)); err != nil {
		return
	}
}

func (f Formatter) parseNodeUpdated(sb *strings.Builder, n *diff.DiffNode, indent string, depth int) {
	if _, err := fmt.Fprintf(sb, "%s- %s: %s\n", indent, n.Key, f.stringify(n.OldValue, depth)); err != nil {
		return
	}

	if _, err := fmt.Fprintf(sb, "%s+ %s: %s", indent, n.Key, f.stringify(n.NewValue, depth)); err != nil {
		return
	}
}
