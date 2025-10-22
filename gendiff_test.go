package code

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func projectPath(parts ...string) string {
	// absolute path per workspace root
	base := "/Users/user/projects/Backend/Go/go-project-244"
	return filepath.Join(append([]string{base}, parts...)...)
}

func TestGenDiffJSONStylish(t *testing.T) {
	file1 := projectPath("examples", "simple", "file1.json")
	file2 := projectPath("examples", "simple", "file2.json")

	out, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err)
	expected := `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`
	assert.Equal(t, expected, out)
}

func TestGenDiffYAMLDefaultFormatStylish(t *testing.T) {
	file1 := projectPath("examples", "simple", "file1.yml")
	file2 := projectPath("examples", "simple", "file2.yml")

	out, err := GenDiff(file1, file2, "")
	require.NoError(t, err)
	expected := `{
  - follow: false
    host: hexlet.io
  - proxy: 123.234.53.22
  - timeout: 50
  + timeout: 20
  + verbose: true
}`
	assert.Equal(t, expected, out)
}

func TestGenDiffPlainFormat(t *testing.T) {
	file1 := projectPath("examples", "simple", "file1.json")
	file2 := projectPath("examples", "simple", "file2.json")

	out, err := GenDiff(file1, file2, "plain")
	require.NoError(t, err)
	expected := "Property 'follow' was removed\n" +
		"Property 'proxy' was removed\n" +
		"Property 'timeout' was updated. From 50 to 20\n" +
		"Property 'verbose' was added with value: true"
	assert.Equal(t, expected, out)
}

func TestGenDiffJSONFormatJSON(t *testing.T) {
	file1 := projectPath("examples", "simple", "file1.json")
	file2 := projectPath("examples", "simple", "file2.json")

	out, err := GenDiff(file1, file2, "json")
	require.NoError(t, err)

	var wrap struct {
		Diff []map[string]any `json:"diff"`
	}
	require.NoError(t, json.Unmarshal([]byte(out), &wrap))
	assert.Greater(t, len(wrap.Diff), 0)
}

func TestGenDiffErrorMissingPaths(t *testing.T) {
	out, err := GenDiff("", "", "")
	assert.Equal(t, "", out)
	require.Error(t, err)
	assert.Equal(t, "you should pass file paths", err.Error())
}

func TestGenDiffErrorReadFile(t *testing.T) {
	file1 := projectPath("examples", "simple", "no_such.json")
	file2 := projectPath("examples", "simple", "file2.json")

	out, err := GenDiff(file1, file2, "stylish")
	assert.Equal(t, "", out)
	require.Error(t, err)
}

func TestGenDiffErrorReadRightFile(t *testing.T) {
	file1 := projectPath("examples", "simple", "file1.json")
	file2 := projectPath("examples", "simple", "no_such.json")

	out, err := GenDiff(file1, file2, "stylish")
	assert.Equal(t, "", out)
	require.Error(t, err)
}

func TestGenDiffErrorExtMismatch(t *testing.T) {
	file1 := projectPath("examples", "simple", "file1.json")
	file2 := projectPath("examples", "simple", "file2.yml")

	out, err := GenDiff(file1, file2, "stylish")
	assert.Equal(t, "", out)
	require.Error(t, err)
	assert.Equal(t, "files must have the same supported extension: got .json and .yml", err.Error())
}

func TestGenDiffErrorUnsupportedExt(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.txt")
	f2 := filepath.Join(dir, "b.txt")
	require.NoError(t, os.WriteFile(f1, []byte("k: v"), 0o644))
	require.NoError(t, os.WriteFile(f2, []byte("k: v2"), 0o644))

	out, err := GenDiff(f1, f2, "stylish")
	assert.Equal(t, "", out)
	require.Error(t, err)
	assert.Equal(t, "unsupported extension: .txt", err.Error())
}

func TestGenDiffErrorUnsupportedFormat(t *testing.T) {
	file1 := projectPath("examples", "simple", "file1.json")
	file2 := projectPath("examples", "simple", "file2.json")

	out, err := GenDiff(file1, file2, "xml")
	assert.Equal(t, "", out)
	require.Error(t, err)
	assert.Equal(t, "unsupported format: xml", err.Error())
}

func TestPickParserDirect(t *testing.T) {
	// valid
	_, err := pickParser(".json")
	require.NoError(t, err)
	_, err = pickParser(".yml")
	require.NoError(t, err)
	_, err = pickParser(".yaml")
	require.NoError(t, err)

	// invalid
	_, err = pickParser(".xml")
	require.Error(t, err)
	assert.Equal(t, "unsupported extension: .xml", err.Error())
}

func TestPickFormatterDirect(t *testing.T) {
	// valid
	_, err := pickFormatter("stylish")
	require.NoError(t, err)
	_, err = pickFormatter("plain")
	require.NoError(t, err)
	_, err = pickFormatter("json")
	require.NoError(t, err)

	// invalid
	_, err = pickFormatter("xml")
	require.Error(t, err)
	assert.Equal(t, "unsupported format: xml", err.Error())
}
