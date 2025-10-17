package json

import (
	"code/internal/domain/diff"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Formatter struct{}

func (f Formatter) Format(nodes []diff.DiffNode) (string, error) {
	var sb strings.Builder

	f.writeNodes(&sb, nodes)

	return sb.String(), nil
}

func (f Formatter) writeNodes(sb *strings.Builder, nodes []diff.DiffNode) {
	sorted := make([]diff.DiffNode, len(nodes))

	copy(sorted, nodes)

	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	prettyJSON, err := json.MarshalIndent(nodes, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	_, err = fmt.Fprint(sb, string(prettyJSON))
	if err != nil {
		return
	}
}
