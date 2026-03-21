package log

type Level string

var logger Logger

func SetLogger(l Logger) {
	logger = l
}

func SetLevel(level Level) {
	if logger == nil {
		panic("logger is not set; call 'SetDefault' first")
	}
	logger.SetLevel(level)
}

func Debug(msg string, args ...any) {
	if logger == nil {
		panic("logger is not set; call 'SetDefault' first")
	}
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	if logger == nil {
		panic("logger is not set; call 'SetDefault' first")
	}
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	if logger == nil {
		panic("logger is not set; call 'SetDefault' first")
	}
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	if logger == nil {
		panic("logger is not set; call 'SetDefault' first")
	}
	logger.Error(msg, args...)
}

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelDebug Level = "debug"
)
