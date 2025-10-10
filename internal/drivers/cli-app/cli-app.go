package cli_app

import (
	"code/internal/interfaces"
	"context"

	"github.com/urfave/cli/v3"
)

type CliApp struct {
	Cli    cli.Command
	Flags  []cli.Flag
	Action func(context.Context, *cli.Command) error
}

var _ interfaces.Cli = (*CliApp)(nil)

func NewCliApp() *CliApp {
	return &CliApp{}
}

func (c *CliApp) Init() error {
	if err := c.GenerateFlags(); err != nil {
		return err
	}

	c.Cli = cli.Command{
		Name: "gendiff",
		//Version: "0.0.1",
		Usage:  "Compares two configuration files and shows a difference.",
		Flags:  c.Flags,
		Action: c.Action,
	}

	return nil
}

func (c *CliApp) GenerateFlags() error {
	c.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "output format",
			Value:   "stylish",
		},
	}

	return nil
}

func (c *CliApp) AddAction(f func(ctx context.Context, command *cli.Command) error) {
	c.Action = f
}

func (c *CliApp) Run(ctx context.Context, args []string) error {
	if err := c.Cli.Run(ctx, args); err != nil {
		return err
	}

	return nil
}
