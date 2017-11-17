package grace

import (
	"fmt"
	"os"
	"os/signal"
)

type SignalFunc func() Signal

type Signal interface {
	RegisterSignal(chan<- os.Signal) error
}

var adapters = make(map[string]SignalFunc)

//type Signal struct {
//	m map[os.Signal]SignalFunc
//}

//func NewSignal() *Signal {
//	return &Signal{
//		m: make(map[os.Signal]SignalFunc),
//	}
//}

//func (sig *Signal) Register(signal os.Signal, f SignalFunc) {
//	if _, ok := sig.m[signal]; !ok {
//		sig.m[signal] = f
//	}
//}
func Register(adapterName string, adapter SignalFunc) {

	if adapter == nil {
		panic("ATC signal: Register handler is nil")
	}
	if _, found := adapters[adapterName]; found {
		panic("ATC signal: Register failed for handler " + adapterName)
	}
	adapters[adapterName] = adapter
}

func NewSignal(adapterName string) (Signal, error) {
	if handler, ok := adapters[adapterName]; ok {
		return handler(), nil
	} else {
		return nil, fmt.Errorf("ATC signal: unknown adapter name %s failed.", adapterName)
	}
}

func Notify() {
	// Signal
	//	1. TERM,INT 立即终止
	//	2. QUIT 平滑终止
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan)
}

//func (sig *Signal) Handler(signal os.Signal, arg interface{}) error {
//	if _, ok := sig.m[signal]; ok {
//		sig.m[signal](arg)
//		return nil
//	}
//
//	return fmt.Errorf("Signal '%v' is not supported.", signal)
//}

//func Notoify() {
//	sig := NewSignal()
//	sig.Register(syscall.SIGHUP, sighup)
//
//	sigChan := make(chan os.Signal)
//	signal.Notify(sigChan)
//
//	for {
//		select {
//		case s := <-sigChan:
//			err := sig.Handler(s, nil)
//			if err != nil {
//				// 退出
//			}
//		}
//	}
//}
