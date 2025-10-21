package diff

import (
	"code"
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func Handler(_ context.Context, command *cli.Command) error {
	file1 := command.Args().Get(0)
	file2 := command.Args().Get(1)
	format := command.String("format")

	if format == "" {
		format = "stylish"
	}

	out, err := code.GenDiff(file1, file2, format)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}
