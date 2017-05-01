package logs

import (
	"testing"
	"time"
)

func testAllLevel(l *AtcLogger) {
	l.Trace("%v", "Trace info.")
	l.Debug("%v", "Debug info.")
	l.Info("%v", "Info info.")
	l.Notice("%v", "Notice info.")
}

func TestStdout(t *testing.T) {
	l1 := NewLogger(10000)
	l1.SetHandler("stdout")
	testAllLevel(l1)
	time.Sleep(1 * time.Second)
}
