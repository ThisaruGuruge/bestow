package log

type Logger interface {
	SetLevel(level Level)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}
