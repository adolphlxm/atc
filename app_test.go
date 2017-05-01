package atc

//import (
//	"fmt"
//	"testing"
//)
//
//func HandlerAppRun1(ctx *Context) {
//	fmt.Println("appRun1")
//	ctx.Write([]byte(ctx.Query("a")))
//}
//
//func HandlerAppRun2(ctx *Context) {
//	fmt.Println("appRun2")
//	ctx.Write([]byte(ctx.Query("a")))
//}
//
//func TestAppRun(t *testing.T) {
//	//appRun 1
//	go func() {
//		atcRun1 := NewApp("127.0.0.1", 8081)
//		atcRun1.Handler.AddRouter("/test1", HandlerAppRun1)
//		atcRun1.Run()
//	}()
//
//	//appRun 2
//	atcRun2 := NewApp("127.0.0.1", 8082)
//	atcRun2.Handler.AddRouter("/test2", HandlerAppRun2)
//	atcRun2.Run()
//}
