package diff_checker

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

func Handler(_ context.Context, command *cli.Command) error {
	file1 := command.Args().Get(0)
	file2 := command.Args().Get(1)

	if file1 == "" || file2 == "" {
		return errors.New("you should pass file paths")
	}
	diff, err := GetDiff(file1, file2)

	fmt.Println(diff)

	return err
}
