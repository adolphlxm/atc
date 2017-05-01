package logs

import (
	"fmt"
)

type Stdout struct {
}

func NewStdoutHandler() IAtcLogger {
	stdout := &Stdout{}
	return stdout
}

// Output message in stdout.
func (s *Stdout) Output(msg string) error {
	fmt.Println(msg)
	return nil
}

//Register NewStdout
func init() {
	Register("stdout", NewStdoutHandler)
}
