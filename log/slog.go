package log

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	logger *slog.Logger
	level  *slog.LevelVar
}

func newSlog(level Level) Logger {
	levelVar := &slog.LevelVar{}
	levelVar.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{Level: levelVar}
	return &SlogLogger{logger: slog.New(slog.NewTextHandler(os.Stdout, opts)), level: levelVar}
}

func (l *SlogLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *SlogLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *SlogLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}
func (l *SlogLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

func (l *SlogLogger) SetLevel(level Level) {
	l.level.Set(slog.Level(level))
}

func (l *SlogLogger) With(args ...interface{}) Logger {
	newLog := l.logger.With(args...)
	return &SlogLogger{logger: newLog, level: l.level}
}
