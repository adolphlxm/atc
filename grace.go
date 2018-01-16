package atc

import (
	"time"
	"context"

	"github.com/adolphlxm/atc/grace"
	"github.com/adolphlxm/atc/logs"
)

func GracePushFront(quit grace.TT) error {
	lazyInit()
	return graceNodeTree.PushFront(quit)
}

func GracePushBack(quit grace.TT) error {
	lazyInit()
	return graceNodeTree.PushBack(quit)
}

func GraceInsertAfter(moduleID string, quit grace.TT) error {
	lazyInit()
	return graceNodeTree.InsertAfter(moduleID, quit)
}

func GraceInsertBefore(moduleID string, quit grace.TT) error {
	lazyInit()
	return graceNodeTree.InsertBefore(moduleID, quit)
}

func GraceRemove(moduleID string) {
	lazyInit()
	graceNodeTree.Remove(moduleID)
}

func GraceMoveAfter(moduleID1, moduleID2 string) {
	lazyInit()
	graceNodeTree.MoveAfter(moduleID1, moduleID2)
}

func GraceMoveBefore(moduleID1, moduleID2 string) {
	lazyInit()
	graceNodeTree.MoveBefore(moduleID1, moduleID2)
}

// lazyInit lazily initializes a zero Grace list value.
func lazyInit() {
	if graceNodeTree == nil {
		graceNodeTree = grace.NewGrace()
	}
}

type httpShutDown struct {}
func (this *httpShutDown) ModuleID() string {
	return "http"
}
func (this *httpShutDown) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(Aconfig.HTTPQTimeout)*time.Second)
	err := HttpAPP.Server.Shutdown(ctx)
	logs.Tracef("shutdown: http, biggest waiting for %ds...", Aconfig.HTTPQTimeout)
	time.Sleep(1 * time.Millisecond)
	return err
}

type thriftShutDown struct {}
func (this *thriftShutDown) ModuleID() string {
	return "thrift"
}
func (this *thriftShutDown) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(Aconfig.ThriftQTimeout)*time.Second)
	err := ThriftRPC.Shutdown(ctx)
	logs.Tracef("shutdown: thrift, biggest waiting for %ds...", Aconfig.ThriftQTimeout)
	return err
}