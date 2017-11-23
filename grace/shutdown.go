package grace
//
//import (
//	"time"
//	"sync/atomic"
//)
//
//type TT interface {
//	// TODO 监听退出信号
//	// TODO 退出命令
//}
//
//type Grace struct {
//
//}
//
//func (this *Grace)
//
//type T int
//
//func Shutdown(ch <-chan T) bool {
//	select {
//	case <-ch:
//		return true
//	default:
//	}
//
//	return false
//}
//var shutdownPollInterval = 500 * time.Millisecond
//
//func (p *TSimpleServer) ShuttingDown() bool {
//	return atomic.LoadInt32(&p.inShutdown) != 0
//}
//
//
//for {
//if srv.closeIdleConns() {
//return lnerr
//}
//select {
//case <-ctx.Done():
//return ctx.Err()
//case <-ticker.C:
//}
//}
