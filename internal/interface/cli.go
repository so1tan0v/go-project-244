package _interface

// Cli - Интерфейс приложения
type Cli interface {
	Init() error
	Run() error
	GenerateFlags() error
}
