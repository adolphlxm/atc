package grace

import (
	"container/list"
	"errors"
)

var err_push = errors.New("grace link table write error.")
var err_insert = errors.New("No linked list elements are found.")

type TT interface {
	ModuleID() string
	Stop() error
}

type Grace struct {
	list list.List
}

func NewGrace() *Grace {
	return &Grace{}
}

// PushFront inserts a new interface TT with value v at the front of list l.
func (this *Grace) PushFront(quit TT) error {
	e := this.list.PushFront(quit)
	if e == nil {
		return err_push
	}
	return nil
}

// PushBack inserts a new interface TT with value v at the back of list l.
func (this *Grace) PushBack(quit TT) error {
	e := this.list.PushBack(quit)
	if e == nil {
		return err_push
	}
	return nil
}

// InsertAfter inserts a new interface TT with value v immediately after mark.
// If mark is not an TT of l, the list is not modified.
func (this *Grace) InsertAfter(moduleID string, quit TT) error {
	if e := this.findElement(moduleID); e != nil {
		ee := this.list.InsertAfter(quit, e)
		if ee == nil {
			return err_push
		}

		return nil
	}

	return err_insert
}

// InsertBefore inserts a new interface TT with value v immediately before mark.
// If mark is not an TT of l, the list is not modified.
func (this *Grace) InsertBefore(moduleID string, quit TT) error {
	if e := this.findElement(moduleID); e != nil {
		ee := this.list.InsertBefore(quit, e)

		if ee == nil {
			return err_push
		}

		return nil
	}

	return err_insert
}

func (this *Grace) Remove(moduleID string){
	if e := this.findElement(moduleID); e != nil {
		this.list.Remove(e)
	}
}

// MoveAfter moves interface TT to its new position after mark.
// If e or mark is not an TT of l, or e == mark, the list is not modified.
func (this *Grace) MoveAfter(moduleID1, moduleID2 string){
	var e1, e2 *list.Element
	for e := this.list.Front(); e != nil; e = e.Next() {
		nt := e.Value.(TT)
		if nt.ModuleID() == moduleID1 {
			e1 = e
		}
		if nt.ModuleID() == moduleID2 {
			e2 = e
		}
	}

	this.list.MoveAfter(e1, e2)
}

// MoveBefore moves interface TT to its new position before mark.
// If e or mark is not an TT of l, or e == mark, the list is not modified.
func (this *Grace) MoveBefore(moduleID1, moduleID2 string){
	var e1, e2 *list.Element
	for e := this.list.Front(); e != nil; e = e.Next() {
		nt := e.Value.(TT)
		if nt.ModuleID() == moduleID1 {
			e1 = e
		}
		if nt.ModuleID() == moduleID2 {
			e2 = e
		}
	}

	this.list.MoveBefore(e1, e2)
}

func (this *Grace) findElement(moduleID string) *list.Element{
	for e := this.list.Front(); e != nil; e = e.Next() {
		nt := e.Value.(TT)
		if nt.ModuleID() == moduleID {
			return e
		}
	}
	return nil
}

func (this *Grace) Stop() error {
	for i, n := 0, this.list.Len(); i < n; i++ {
		e := this.list.Back()
		if e == nil {
			break
		}
		nt := e.Value.(TT)

		if err := nt.Stop(); err != nil {
			// TODO 如果关闭发生错误，则继续循环关闭
			continue
		}

		// Delete
		this.list.Remove(e)
	}

	return nil
}