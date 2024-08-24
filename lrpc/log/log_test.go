package log

import "testing"

func TestLog(t *testing.T) {
	Debug("debug", "hello", "kitty")
	Info("info", "hello", "kitty")
	Warnf("warnf %s", "ttoto")
	Error("error", "hello", "kitty")
}
