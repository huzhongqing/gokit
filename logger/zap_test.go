package logger

import "testing"

func TestZap(t *testing.T) {
	Zap.Debug("debug")
	Zap.Info("info")
	Zap.Warn("warn")
	Zap.Error("error")
	Zap.Panic("panic")
	Zap.Fatal("fatal")
}
