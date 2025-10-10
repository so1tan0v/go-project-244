package interfaces

import "code/internal/domain/diff"

type Formatter interface {
	Name() string
	Format(nodes []diff.DiffNode) (string, error)
}
