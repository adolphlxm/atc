package mgo

import (
	"testing"
	"fmt"
)

const addrs = "mongodb://127.0.0.1:27017"

var M *MgoDB

func initMgo(){
	var err error
	M, err = NewMgoDB(addrs)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type Data struct {
	A int
	B string
}

func TestMgoDB_Insert(t *testing.T) {
	initMgo()

	data := &Data{A:1,B:"atc"}
	err := M.Session().DB("test1").C("test1").Insert(data)
	if err != nil {
		t.Error(err)
	}
}