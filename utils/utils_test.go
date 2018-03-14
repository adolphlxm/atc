package utils

import (
	"fmt"
	"testing"
	"regexp"
)

func TestRandString(t *testing.T) {
	fmt.Println(RandString(10))
}

func TestVerifyMobile(t *testing.T) {
	if b := VerifyMobile("18691635352"); !b {
		t.Error(b)
	}
}

func TestVerifyEmail(t *testing.T) {
	if b := VerifyEmail("1235355@qq.com"); !b {
		t.Error(b)
	}
}