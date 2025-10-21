package code

import (
	jsonformatter "code/internal/drivers/formmaters/json"
	"code/internal/drivers/formmaters/plain"
	"code/internal/drivers/formmaters/stylish"
	"code/internal/drivers/parsers/jsonparser"
	"code/internal/drivers/parsers/yamlparser"
	"code/internal/interfaces"

	usecase "code/internal/usecase/diff"

	"fmt"
	"os"
	"path/filepath"
)

func GenDiff(file1, file2, format string) (string, error) {
	if file1 == "" || file2 == "" {
		return "", fmt.Errorf("you should pass file paths")
	}

	leftRaw, err := os.ReadFile(file1)
	if err != nil {
		return "", err
	}
	rightRaw, err := os.ReadFile(file2)
	if err != nil {
		return "", err
	}

	ext1 := filepath.Ext(file1)
	ext2 := filepath.Ext(file2)
	if ext1 == "" || ext1 != ext2 {
		return "", fmt.Errorf("files must have the same supported extension: got %s and %s", ext1, ext2)
	}

	parser, err := pickParser(ext1)
	if err != nil {
		return "", err
	}

	if format == "" {
		format = "stylish"
	}

	formatter, err := pickFormatter(format)
	if err != nil {
		return "", err
	}

	svc := usecase.NewService(parser, formatter)
	return svc.GenerateDiff(leftRaw, rightRaw)
}

func pickParser(ext string) (interfaces.Parser, error) {
	switch ext {
	case ".json":
		return jsonparser.Parser{}, nil
	case ".yml", ".yaml":
		return yamlparser.Parser{}, nil
	default:
		return nil, fmt.Errorf("unsupported extension: %s", ext)
	}
}

func pickFormatter(format string) (interfaces.Formatter, error) {
	switch format {
	case "stylish":
		return stylish.Formatter{}, nil
	case "plain":
		return plain.Formatter{}, nil
	case "json":
		return jsonformatter.Formatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
