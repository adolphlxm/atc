package logs

import (
	"os"
	"path/filepath"
	"time"
	"fmt"
	"bufio"
	"sync"
	"bytes"
	"runtime"
	"errors"
)

var (
	pid      = os.Getpid()
	program  = filepath.Base(os.Args[0])
)

type File struct {
	mu      sync.Mutex
	*bufio.Writer
	logFile *os.File

	LogDir string `json:"logdir"`

	MaxSize uint64 `json:"maxsize"`
	Buffersize int `json:"buffersize"`
	FlushInterval uint64 `json:"flushinterval"`

	nbytes uint64 // The number of bytes written to this file

}

const flushInterval = 10
// bufferSize sizes the buffer associated with each log file. It's large
// so that log records can accumulate without the logging thread blocking
// on disk I/O. The flushDaemon will block instead.
const bufferSize = 256 * 1024
const maxsize = 1024 * 1024 * 1800

func NewFileLog() IAtcLogger {
	file := &File{
		MaxSize: maxsize,
		Buffersize:bufferSize,
		FlushInterval:flushInterval,
	}

	go file.flushDaemon()

	return file
}

func (this *File) Init(config interface{}) error {
	// Parsing the struct
	logConf, ok := config.(*File)
	if !ok {
		return errors.New("Parsing struct fails.")
	}

	this.LogDir = logConf.LogDir
	if logConf.MaxSize > 0 {
		this.MaxSize = logConf.MaxSize
	}
	if logConf.Buffersize > 0 {
		this.Buffersize = logConf.Buffersize
	}
	if logConf.FlushInterval > 0 {
		this.FlushInterval = logConf.FlushInterval
	}
	_, err := os.Stat(this.LogDir)
	if os.IsNotExist(err) {
		this.LogDir = ""
	}

	return this.rotateFile(time.Now())
}

func (this *File) Output(msg []byte) error {
	var err error
	now := time.Now()

	if this.nbytes + uint64(len(msg)) >= this.MaxSize {
		if err = this.rotateFile(now); err != nil {
			return err
		}
	}

	n, err := this.Writer.Write(msg)
	this.nbytes += uint64(n)
	return err
}

func (this *File) Flush(){
	this.lockAndFlushAll()
}

func (this *File) rotateFile(now time.Time) error {
	if this.logFile != nil {
		this.Flush()
		this.logFile.Close()
	}

	var err error
	this.logFile, _, err = this.create("log",now)
	this.nbytes = 0
	if err != nil {
		return nil
	}

	this.Writer = bufio.NewWriterSize(this.logFile, this.Buffersize)

	// Write header.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Log file created at: %s\n", now.Format("2006/01/02 15:04:05"))
	fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(&buf, "Log line format: [FEWNITD]mmdd hh:mm:ss.uuuuuu threadid file:line] msg\n")
	n, err := this.logFile.Write(buf.Bytes())
	this.nbytes += uint64(n)
	return err

}

func (this *File) create(tag string, t time.Time) (f *os.File, filename string, err error){
	name, link := this.name(tag,time.Now())

	fname := filepath.Join(this.LogDir, name)
	f, err = os.Create(fname)
	if err == nil {
		symlink := filepath.Join(this.LogDir, link)
		os.Remove(symlink)        // ignore err
		os.Symlink(name, symlink) // ignore err
		return f, fname, nil
	}

	return nil, "", fmt.Errorf("log: cannot create log: %v", err)
}

func (this *File) lockAndFlushAll() {
	this.mu.Lock()
	this.flushAll()
	this.mu.Unlock()
}

func (this *File) flushAll(){
	if this.logFile != nil {
		this.Writer.Flush()
		this.logFile.Sync()
	}
}

// flushDaemon periodically flushes the log file buffers.
func (this *File)flushDaemon() {
	for _ = range time.NewTicker(time.Duration(this.FlushInterval) * time.Second).C {
		this.lockAndFlushAll()
	}
}

func (this *File) name(tag string, t time.Time) (name, link string) {
	if tag != "" {
		tag = "." + tag
	}

	name = fmt.Sprintf("%s%s.%04d%02d%02d-%02d%02d%02d.%d",
		program,
		tag,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		pid)

	return name, program + tag
}

func init() {
	Register("file", NewFileLog)
}
