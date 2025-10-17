package interfaces

import "code/internal/domain/diff"

type Formatter interface {
	Format(nodes []diff.DiffNode) (string, error)
}
