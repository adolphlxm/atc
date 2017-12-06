package logs

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

// Log message level
const (
	LevelFatal = iota
	LevelError
	LevelWarn
	LevelNotice
	LevelInfo
	LevelTrace
	LevelDebug
)

const TimeFormat = "2006/01/02 15:04:05.000000"

// Name for adapter with ATC official support
const (
	AdapterStdout = "stdout"
	AdapterFile   = "file"
)

var LevelName [7]string = [7]string{"Fatal", "Error", "Warn", "Notice", "Info", "Trace", "Debug"}

type LoggerFunc func() IAtcLogger

type IAtcLogger interface {
	Init(config string) error
	Output(msg string) error
}

var adapters = make(map[string]LoggerFunc)

func Register(adapterName string, handler LoggerFunc) {
	if adapters == nil {
		panic("ATC logs: Register LoggerFunc is nil")
	}
	if _, found := adapters[adapterName]; found {
		panic("ATC logs: Register failed for LoggerFunc " + adapterName)
	}

	adapters[adapterName] = handler
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
		level:   LevelDebug,
		skip:    3,
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

func (l *AtcLogger) SetLogger(adapterName string, configs ...string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	cf := append(configs, "{}")[0]

	if handler, ok := adapters[adapterName]; ok {
		l.handler[adapterName] = handler()
		err := l.handler[adapterName].Init(cf)
		if err != nil {
			return fmt.Errorf("ATC logs: %q handler fail, err:%v.", adapterName, err.Error())
		}
	} else {
		return fmt.Errorf("ATC logs: %q handler setting fail.", adapterName)
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
					fmt.Printf("ATC logs: Output handler fail, err:%v\n", err.Error())
				}
			}
		}
	}
}

func (l *AtcLogger) Output(level int, msg string) error {
	now := time.Now().Format(TimeFormat)
	l.mu.Lock()
	defer l.mu.Unlock()
	if level > l.level {
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

func (l *AtcLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelError, msg)
}

func (l *AtcLogger) Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelFatal, msg)
	os.Exit(1)
}

// Defaultlogs is the default ServeMux used by Serve.
var defaultlogs = NewLogger(10000)

func SetLogger(adapterName string, configs ...string) {
	defaultlogs.SetLogger(adapterName, configs...)
}

func SetLevel(level int) {
	defaultlogs.SetLevel(level)
}

func Trace(format string, v ...interface{}) {
	defaultlogs.Trace(format, v...)
}

func Debug(format string, v ...interface{}) {
	defaultlogs.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	defaultlogs.Info(format, v...)
}

func Notice(format string, v ...interface{}) {
	defaultlogs.Notice(format, v...)
}

func Warn(format string, v ...interface{}) {
	defaultlogs.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	defaultlogs.Error(format, v...)
}

func Fatal(format string, v ...interface{}) {
	defaultlogs.Fatal(format, v...)
}
