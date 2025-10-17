package diff

import (
	jsonformatter "code/internal/drivers/formmaters/json"
	"code/internal/drivers/formmaters/plain"
	"code/internal/drivers/formmaters/stylish"
	"code/internal/drivers/parsers/jsonparser"
	"code/internal/drivers/parsers/yamlparser"
	"code/internal/interfaces"

	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	usecase "code/internal/usecase/diff"

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

	var parser interfaces.Parser
	switch ext1 {
	case ".json":
		parser = jsonparser.Parser{}
	case ".yml":
		fallthrough
	case ".yaml":
		parser = yamlparser.Parser{}
	default:
		return fmt.Errorf("unsupported extension: %s", ext1)
	}

	format := command.String("format")
	if format == "" {
		format = "stylish"
	}

	var formatter interfaces.Formatter
	switch format {
	case "stylish":
		formatter = stylish.Formatter{}
	case "plain":
		formatter = plain.Formatter{}
	case "json":
		formatter = jsonformatter.Formatter{}
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}

	svc := usecase.NewService(parser, formatter)
	out, err := svc.GenerateDiff(leftRaw, rightRaw)
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}
