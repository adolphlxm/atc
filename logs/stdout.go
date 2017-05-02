package logs

import (
	"fmt"
)

type Stdout struct {
}

func NewStdoutLog() IAtcLogger {
	stdout := &Stdout{}
	return stdout
}

func (s *Stdout) Init(config string) error {
	return nil
}

// Output message in stdout.
func (s *Stdout) Output(msg string) error {
	fmt.Println(msg)
	return nil
}

//Register NewStdout
func init() {
	Register("stdout", NewStdoutLog)
}
