package interfaces

type Logger interface {
	Write(...any) error
}
