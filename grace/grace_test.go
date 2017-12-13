package grace

import (
	"testing"
	"fmt"
	"time"
)

type GracePushFront struct {}
func (this *GracePushFront) ModuleID() string {
	return "GracePushFront"
}
func (this *GracePushFront) Stop() error {
	fmt.Println("GracePushFront stop.")
	return nil
}

type GracePushBack struct {}
func (this *GracePushBack) ModuleID() string {
	return "GracePushBack"
}
func (this *GracePushBack) Stop() error {
	fmt.Println("GracePushBack stop.")
	return nil
}

type GraceInsertAfter struct {}
func (this *GraceInsertAfter) ModuleID() string {
	return "GraceInsertAfter"
}
func (this *GraceInsertAfter) Stop() error {
	fmt.Println("GraceInsertAfter stop.")
	fmt.Println("sleep 2s")
	time.Sleep(2 * time.Second)
	return nil
}

type GraceInsertBefore struct {}
func (this *GraceInsertBefore) ModuleID() string {
	return "GraceInsertBefore"
}
func (this *GraceInsertBefore) Stop() error {
	fmt.Println("GraceInsertBefore stop.")
	return nil
}

var graceTest = NewGrace()

func Test_all(t *testing.T) {
	TestGrace_PushFront(t)
	TestGrace_PushBack(t)
	TestGrace_InsertAfter(t)
	TestGrace_InsertBefore(t)
	TestGrace_Remove(t)

	graceTest.Stop()
}

func TestGrace_PushFront(t *testing.T) {
	err := graceTest.PushFront(&GracePushFront{})
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGrace_PushBack(t *testing.T) {
	err := graceTest.PushBack(&GracePushBack{})
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGrace_InsertAfter(t *testing.T) {
	err := graceTest.InsertAfter("GracePushFront", &GraceInsertAfter{})
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGrace_InsertBefore(t *testing.T) {
	err := graceTest.InsertBefore("GraceInsertAfter",&GraceInsertBefore{})
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestGrace_Remove(t *testing.T) {
	graceTest.Remove("GracePushFront")
}