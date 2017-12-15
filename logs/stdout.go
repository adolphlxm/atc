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

func (s *Stdout) Init(config interface{}) error {
	return nil
}

// Output message in stdout.
func (s *Stdout) Output(msg []byte) error {
	fmt.Println(string(msg))
	return nil
}

func (s *Stdout) Flush(){

}

//Register NewStdout
func init() {
	Register("stdout", NewStdoutLog)
}
