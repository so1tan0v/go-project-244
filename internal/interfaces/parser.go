package interfaces

type Parser interface {
	Parse(data []byte) (map[string]any, error)
}
