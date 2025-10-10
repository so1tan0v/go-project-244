package diff_checker

import (
	jsonprocessor "code/internal/drivers/diff-checker/json-processor"
	"code/internal/drivers/helper"
	"code/internal/interfaces"
	"fmt"
	"path/filepath"
	"strings"
)

type jsonDiffer struct {
	p jsonprocessor.JsonProcessor
}

func (d jsonDiffer) Diff(pathToFile1, pathToFile2 string) (string, error) {
	json1, err := d.p.GetContent(pathToFile1)
	if err != nil {
		return "", err
	}

	json2, err := d.p.GetContent(pathToFile2)
	if err != nil {
		return "", err
	}

	diff := helper.MergeSlices(d.p.Compare(json1, json2, "-"), d.p.Compare(json2, json1, "+"))

	return d.p.GetComparisonResult(diff)
}

func GetDiff(pathToFile1, pathToFile2 string) (string, error) {
	ext1 := strings.ToLower(filepath.Ext(pathToFile1))
	ext2 := strings.ToLower(filepath.Ext(pathToFile2))

	if ext1 != ext2 {
		return "", fmt.Errorf("differenct format file")
	}

	differs := map[string]interfaces.Differ{
		".json": jsonDiffer{p: jsonprocessor.JsonProcessor{}},
	}

	differ, ok := differs[ext1]
	if !ok {
		return "", fmt.Errorf("Unsupported diff format file")
	}

	return differ.Diff(pathToFile1, pathToFile2)
}
