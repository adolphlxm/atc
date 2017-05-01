package logs

import (
	//"sync/atomic"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

//Log message level
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelFatal
)

const TimeFormat = "2006/01/02 15:04:05.000000"

var LevelName [7]string = [7]string{"Trace", "Debug", "Info", "Notice", "Warn", "Error", "Fatal"}

type LoggerHandlerType func() IAtcLogger

type IAtcLogger interface {
	Output(msg string) error
}

var LoggerHandler = make(map[string]LoggerHandlerType)

func Register(name string, handler LoggerHandlerType) {
	if LoggerHandler == nil {
		panic("ATC logs: Register handler is nil")
	}
	if _, found := LoggerHandler[name]; found {
		panic("ATC logs: Register failed for handler " + name)
	}
	LoggerHandler[name] = handler
}

type AtcLogger struct {
	mu      sync.Mutex
	handler map[string]IAtcLogger

	skip  int
	level int

	msg   chan string
	close int32
}

func NewLogger(channellen int64) *AtcLogger {
	loger := &AtcLogger{
		handler: make(map[string]IAtcLogger),
		level:   LevelFatal,
		skip:    2,
		msg:     make(chan string, channellen),
	}

	go loger.Run()

	return loger
}

func (l *AtcLogger) SetSkip(skip int) {
	l.skip = skip
}

func (l *AtcLogger) SetLevel(level int) {
	l.level = level
}

func (l *AtcLogger) SetHandler(name string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if handler, ok := LoggerHandler[name]; ok {
		l.handler[name] = handler()
	} else {
		return fmt.Errorf("ATC logs: %q handler setting fail.", name)
	}

	return nil
}

func (l *AtcLogger) Run() {
	for {
		select {
		case msg := <-l.msg:
			for _, ll := range l.handler {
				err := ll.Output(msg)
				if err != nil {
					fmt.Println("ATC logs: Run the handler to fail.")
				}
			}
		}
	}
}

func (l *AtcLogger) Output(level int, msg string) error {
	now := time.Now().Format(TimeFormat)
	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.level {
		return nil
	}

	_, file, line, ok := runtime.Caller(l.skip)
	if !ok {
		file = "???"
		line = 0
	}
	_, filename := path.Split(file)
	msg = fmt.Sprintf("[ATC] [%s] %s %s#%d: %s", LevelName[level], now, filename, line, msg)

	l.msg <- msg
	return nil
}

func (l *AtcLogger) Trace(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelTrace, msg)
}

func (l *AtcLogger) Debug(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelDebug, msg)
}

func (l *AtcLogger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelInfo, msg)
}

func (l *AtcLogger) Notice(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelNotice, msg)
}

func (l *AtcLogger) Warn(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelWarn, msg)
}

func (l *AtcLogger) Error(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	l.Output(LevelError, msg)
}
func (l *AtcLogger) Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelError, msg)
}

func (l *AtcLogger) Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelFatal, msg)
	os.Exit(1)
}
