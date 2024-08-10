package log

type Logger interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
}