package grace

import (
	"os"
	"syscall"
	"os/signal"
	"fmt"
)

type SignalFunc func(arg interface{})

type Signal struct {
	m map[os.Signal]SignalFunc
}

func NewSignal() *Signal{
	return &Signal{
		m:make(map[os.Signal]SignalFunc),
	}
}

func (sig *Signal) Register(signal os.Signal, f SignalFunc) {
	if _, ok := sig.m[signal]; !ok {
		sig.m[signal] = f
	}
}

func (sig *Signal) Handler(signal os.Signal, arg interface{}) error {
	if _, ok := sig.m[signal]; ok {
		sig.m[signal](arg)
		return nil
	}

	return fmt.Errorf("Signal '%v' is not supported.", signal)
}

func sighup (arg interface{}){

}

func Notoify(){
	sig := NewSignal()
	sig.Register(syscall.SIGHUP,sighup)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)

	for {
		select {
		case s := <- sigChan:
			err := sig.Handler(s, nil)
			if err != nil {
				// 退出
			}
		}
	}
}