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

var LevelName [7]string = [7]string{"F", "E", "W", "N", "I", "T", "D"}

type LoggerFunc func() IAtcLogger

type IAtcLogger interface {
	Init(config interface{}) error
	Output(msg []byte) error
	Flush()
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

	msg   chan []byte
	close int32
}

func NewLogger(channellen int64) *AtcLogger {
	loger := &AtcLogger{
		handler: make(map[string]IAtcLogger),
		level:   LevelDebug,
		skip:    3,
		msg:     make(chan []byte, channellen),
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

func (l *AtcLogger) SetLogger(adapterName string, configs ...interface{}) error {
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

	//
	if level < LevelNotice || level == LevelDebug {
		_, file, line, ok := runtime.Caller(l.skip)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		msg = fmt.Sprintf("[%s] [%s %s:%d] %s\n", LevelName[level], now, filename, line, msg)
	} else {
		msg = fmt.Sprintf("[%s] [%s] %s\n", LevelName[level], now, msg)
	}

	l.msg <- []byte(msg)
	return nil
}

func (l *AtcLogger) Flush(){
	for _, ll := range l.handler {
		ll.Flush()
	}
}

func (l *AtcLogger) Trace(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.Output(LevelTrace, msg)
}
func (l *AtcLogger) Tracef(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelTrace, msg)
}

func (l *AtcLogger) Debug(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.Output(LevelDebug, msg)
}
func (l *AtcLogger) Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelDebug, msg)
}

func (l *AtcLogger) Info(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.Output(LevelInfo, msg)
}
func (l *AtcLogger) Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelInfo, msg)
}

func (l *AtcLogger) Notice(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.Output(LevelNotice, msg)
}
func (l *AtcLogger) Noticef(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelNotice, msg)
}

func (l *AtcLogger) Warn(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.Output(LevelWarn, msg)
}
func (l *AtcLogger) Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelWarn, msg)
}

func (l *AtcLogger) Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.Output(LevelError, msg)
}
func (l *AtcLogger) Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelError, msg)
}

func (l *AtcLogger) Fatal(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.Output(LevelFatal, msg)
	os.Exit(1)
}
func (l *AtcLogger) Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Output(LevelFatal, msg)
}

// Defaultlogs is the default ServeMux used by Serve.
var defaultlogs = NewLogger(10000)
func SetLogger(adapterName string, configs ...interface{}) error {
	return defaultlogs.SetLogger(adapterName, configs...)
}

func SetLevel(level int) {
	defaultlogs.SetLevel(level)
}

func Trace(args ...interface{}) {
	defaultlogs.Trace(args...)
}
func Tracef(format string, v ...interface{}) {
	defaultlogs.Tracef(format, v...)
}

func Debug(args ...interface{}) {
	defaultlogs.Debug(args...)
}
func Debugf(format string, v ...interface{}) {
	defaultlogs.Debugf(format, v...)
}

func Info(args ...interface{}) {
	defaultlogs.Info(args...)
}
func Infof(format string, v ...interface{}) {
	defaultlogs.Infof(format, v...)
}

func Notice(args ...interface{}) {
	defaultlogs.Notice(args...)
}
func Noticef(format string, v ...interface{}) {
	defaultlogs.Noticef(format, v...)
}

func Warn(args ...interface{}) {
	defaultlogs.Warn(args...)
}
func Warnf(format string, v ...interface{}) {
	defaultlogs.Warnf(format, v...)
}

func Error(args ...interface{}) {
	defaultlogs.Error(args...)
}
func Errorf(format string, v ...interface{}) {
	defaultlogs.Errorf(format, v...)
}

func Fatal(args ...interface{}) {
	defaultlogs.Fatal(args...)
}
func Fatalf(format string, v ...interface{}) {
	defaultlogs.Fatalf(format, v...)
}

func Flush(){
	defaultlogs.Flush()
}