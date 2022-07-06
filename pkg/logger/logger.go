package logger

type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}
