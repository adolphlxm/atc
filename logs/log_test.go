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
	l1.SetLogger(AdapterStdout)
	testAllLevel(l1)
	time.Sleep(1 * time.Millisecond)
}

func TestFile(t *testing.T) {
	l := NewLogger(10000)
	l.SetLogger(AdapterFile, &File{LogDir:"./", MaxSize:1000,Buffersize:1000,FlushInterval:5})
	l.SetLevel(LevelDebug)
	testAllLevel(l)
	time.Sleep(1 * time.Millisecond)
	l.Flush()
}

func TestLogs(t *testing.T){
	Trace("%v", "Trace info.")
	Info("%v", "Info info.")
	Notice("%v", "Notice info.")
	Debug("%v", "Debug info.")
	time.Sleep(1 * time.Millisecond)
}