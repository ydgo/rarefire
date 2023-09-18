package log

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	SetLevel(level Level)
	With(args ...interface{}) Logger
}

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var defaultLog = newSlog(InfoLevel)

func Debug(msg string, args ...interface{}) {
	defaultLog.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	defaultLog.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	defaultLog.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	defaultLog.Error(msg, args...)
}

func SetLevel(level Level) {
	defaultLog.SetLevel(level)
}

func With(args ...interface{}) Logger {
	return defaultLog.With(args...)
}
