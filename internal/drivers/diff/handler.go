package diff

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"code/internal/drivers/jsonparser"
	"code/internal/drivers/stylishformatter"
	usecase "code/internal/usecase/gendiff"

	"github.com/urfave/cli/v3"
)

func Handler(_ context.Context, command *cli.Command) error {
	file1 := command.Args().Get(0)
	file2 := command.Args().Get(1)

	if file1 == "" || file2 == "" {
		return errors.New("you should pass file paths")
	}

	leftRaw, err := os.ReadFile(file1)
	if err != nil {
		return err
	}
	rightRaw, err := os.ReadFile(file2)
	if err != nil {
		return err
	}

	ext1 := filepath.Ext(file1)
	ext2 := filepath.Ext(file2)
	if ext1 == "" || ext1 != ext2 {
		return fmt.Errorf("files must have the same supported extension: got %s and %s", ext1, ext2)
	}

	var parser jsonparser.Parser
	switch ext1 {
	case ".json":
	default:
		return fmt.Errorf("unsupported extension: %s", ext1)
	}

	format := command.String("format")
	if format == "" {
		format = "stylish"
	}

	switch format {
	case "stylish":
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	var formatter stylishformatter.Formatter

	svc := usecase.NewService(parser, formatter)
	out, err := svc.GenerateDiff(leftRaw, rightRaw)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}
