package interfaces

// Differ provides a format-agnostic diff operation over two file paths.
// Implementations can use any internal types and parsing strategies.
type Differ interface {
	Diff(pathToFile1, pathToFile2 string) (string, error)
}
