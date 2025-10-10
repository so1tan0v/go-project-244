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
	Key      string
	Type     NodeType
	OldValue any
	NewValue any
	Children []DiffNode
}
