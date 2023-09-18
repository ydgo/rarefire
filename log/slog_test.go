package log

import "testing"

func TestDebug(t *testing.T) {
	l := newSlog(DebugLevel)
	l.SetLevel(InfoLevel)
	l.Debug("hello world", "name", "go")
	l.Info("hello world", "name", "go")
	l.Error("hello world", "name", "go")
	l.With("uid", 1523).Error("hello world", "name", "c++")
}

func TestDefaultLog(t *testing.T) {
	Debug("hello world", "name", "go")
	Info("hello world", "name", "go")
	SetLevel(ErrorLevel)
	Error("hello world", "name", "go")
}
