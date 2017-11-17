package logs

import (
	"testing"
	"time"
)

func testAllLevel(l *AtcLogger) {
	l.Trace("%v", "Trace info.")
	l.Info("%v", "Info info.")
	l.Notice("%v", "Notice info.")
	l.Debug("%v", "Debug info.")
}

func TestStdout(t *testing.T) {
	l1 := NewLogger(10000)
	l1.SetHandler("stdout")
	l1.SetLevel(LevelDebug)
	testAllLevel(l1)
	time.Sleep(2 * time.Second)
}

func TestFile(t *testing.T) {
	l := NewLogger(10000)
	l.SetHandler("file", `{"filename":"test.log","perm":"0660"}`)
	l.SetLevel(LevelDebug)
	testAllLevel(l)
	time.Sleep(1 * time.Second)
}
