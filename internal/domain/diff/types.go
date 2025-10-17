package diff

type NodeType int

const (
	NodeUnchanged NodeType = iota
	NodeAdded
	NodeRemoved
	NodeUpdated
	NodeNested
)

type DiffNode struct {
	Key      string     `json:"key"`
	Type     NodeType   `json:"-"`
	OldValue any        `json:"oldValue"`
	NewValue any        `json:"newValue"`
	Children []DiffNode `json:"children"`
}
