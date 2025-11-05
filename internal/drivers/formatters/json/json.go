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

	if err := f.WriteNodes(&sb, nodes); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func (f Formatter) WriteNodes(sb *strings.Builder, nodes []diff.DiffNode) error {
	sorted := make([]diff.DiffNode, len(nodes))

	copy(sorted, nodes)

	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	wrap := map[string]any{
		"diff": nodes,
	}

	prettyJSON, err := json.MarshalIndent(wrap, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	_, err = fmt.Fprint(sb, string(prettyJSON))
	if err != nil {
		return err
	}

	return nil
}
