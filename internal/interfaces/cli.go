package interfaces

import (
	"context"

	"github.com/urfave/cli/v3"
)

// Cli - Интерфейс приложения
type Cli interface {
	Init() error
	Run(ctx context.Context, args []string) error
	GenerateFlags() error
	AddAction(func(ctx context.Context, command *cli.Command) error)
}
