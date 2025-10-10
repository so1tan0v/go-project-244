package json_processor

import (
	"code/internal/drivers/helper"
	"code/internal/interfaces"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type JsonProcessor struct{}

var _ interfaces.Processor[map[string]any] = (*JsonProcessor)(nil)

func (j JsonProcessor) GetContent(filePath string) (map[string]any, error) {
	r := map[string]any{}

	f, err := os.Open(filePath)
	if err != nil {
		return r, err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	fc, err := io.ReadAll(f)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(fc, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

func (j JsonProcessor) Compare(file1, file2 map[string]any, args ...any) []interfaces.CompareResult {
	result := make([]interfaces.CompareResult, 0)

	for k, v1 := range file1 {
		v2, found := file2[k]

		if !found {
			result = append(result, interfaces.CompareResult{Key: k, Value: v1, Flag: args[0].(string)})

			continue
		}

		if reflect.TypeOf(v1).String() == "map" {
			m1, ok := v1.(map[string]interface{})
			if !ok {
				result = append(result, interfaces.CompareResult{Key: k, Value: v1})

				continue
			}

			m2, ok := v2.(map[string]interface{})
			if !ok {
				result = append(result, interfaces.CompareResult{Key: k, Value: v1})

				continue
			}

			result = append(result, interfaces.CompareResult{Key: k, Value: j.Compare(m1, m2, args)})

			continue
		}

		if v1 != v2 {
			result = append(result, interfaces.CompareResult{Key: k, Value: v1, Flag: args[0].(string)})

			continue
		}

		result = append(result, interfaces.CompareResult{Key: k, Value: v1})
	}

	return result
}

func (j JsonProcessor) GetComparisonResult(diffResult []interfaces.CompareResult) (string, error) {
	sort.SliceStable(diffResult, func(i, j int) bool {
		return diffResult[i].Key < diffResult[j].Key
	})

	diffResult = helper.FilterSlices(diffResult, func(k int, i interfaces.CompareResult) bool {
		if i.Flag == "" && k == slices.IndexFunc(diffResult, func(result interfaces.CompareResult) bool {
			return result.Key == i.Key
		}) {
			return false
		}

		return true
	})

	var sb strings.Builder
	for k, result := range diffResult {
		sb.WriteString(j.getStringResult(result, 2))

		if k < len(diffResult)-1 {
			sb.WriteString("\n")
		}
	}

	return fmt.Sprintf("{\n%s\n}", sb.String()), nil
}

func (j JsonProcessor) getStringResult(diffResult interfaces.CompareResult, tabs int) string {
	var flag string = ""
	if diffResult.Flag != "" {
		flag = diffResult.Flag + " "
	} else {
		flag = "  "
	}

	return fmt.Sprintf(
		"%s%s%s: %s",
		strings.Repeat(" ", tabs),
		flag,
		diffResult.Key,
		j.parseValue(diffResult.Value),
	)
}

func (j JsonProcessor) parseValue(value any) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%v", value)
	}
}
