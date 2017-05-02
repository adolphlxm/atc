package logs

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

type File struct {
	logFile *os.File

	Filename string `json:"filename"`

	Perm string `json:"perm"`
	fPerm os.FileMode

	suffix string
}

func NewFileLog() IAtcLogger {
	file := &File{
		fPerm:   0660,
		suffix: ".log",
	}
	return file
}
type Test struct {
	Filename string `json:"filename"`
	Perm uint32 `json:"perm"`
}
func (f *File) Init(config string) error {

	err := json.Unmarshal([]byte(config), f)
	if err != nil {
		return err
	}

	perm, err := strconv.ParseInt(f.Perm, 8, 64)
	f.fPerm = os.FileMode(perm)

	if len(f.Filename) == 0 {
		return errors.New("config must have filename.")
	}

	f.suffix = filepath.Ext(f.Filename)
	if len(f.suffix) == 0 {
		f.suffix = ".log"
	}
	return err
}

func (f *File) Output(msg string) error {
	var err error

	msg += "\n"
	f.logFile, err = os.OpenFile(f.Filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, f.fPerm)
	defer f.logFile.Close()
	if err == nil {
		os.Chmod(f.Filename, f.fPerm)
	}

	f.logFile.WriteString(msg)

	return err
}

func init() {
	Register("file", NewFileLog)
}
