package code

import (
	jsonformatter "code/internal/drivers/formatters/json"
	"code/internal/drivers/formatters/plain"
	"code/internal/drivers/formatters/stylish"
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

	leftRaw, err := getFileContent(file1)
	if err != nil {
		return "", err
	}
	rightRaw, err := getFileContent(file2)
	if err != nil {
		return "", err
	}

	if err := validateFiles(file1, file2); err != nil {
		return "", err
	}

	parser, err := pickParser(filepath.Ext(file1))
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

func getFileContent(filePath string) ([]byte, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, err
	}

	fullpath := filepath.Clean(filePath)

	return os.ReadFile(fullpath)
}

func validateFiles(filePath1 string, filePath2 string) error {
	ext1 := filepath.Ext(filePath1)
	ext2 := filepath.Ext(filePath2)

	if err := validateFileExtension(filePath1); err != nil {
		return err
	}

	if err := validateFileExtension(filePath2); err != nil {
		return err
	}

	if ext1 == "" || ext1 != ext2 {
		return fmt.Errorf("files must have the same supported extension: got %s and %s", ext1, ext2)
	}

	return nil
}

func validateFileExtension(filePath string) error {
	ext := filepath.Ext(filePath)
	if ext == "" {
		return fmt.Errorf("file extension is required")
	}

	if ext != ".json" && ext != ".yml" && ext != ".yaml" {
		return fmt.Errorf("unsupported extension: %s", ext)
	}

	return nil
}
