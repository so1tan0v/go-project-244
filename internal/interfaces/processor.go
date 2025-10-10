package interfaces

type Processor[T any] interface {
	Compare(file1, file2 T, args ...any) []CompareResult
	GetComparisonResult(c []CompareResult) (string, error)
	GetContent(filePath string) (T, error)
}
