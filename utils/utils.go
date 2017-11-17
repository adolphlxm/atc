package utils

import (
	"math/rand"
	"regexp"
	"time"
)

const (
	RegMobile = "^1[34578]\\d{9}$"
	RegEmail  = "^[A-Za-z\\d]+([-_.][A-Za-z\\d]+)*@([A-Za-z\\d]+[-.])+[A-Za-z\\d]{2,4}$"
)

// 生成随机字符串
func RandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}

// 手机号验证
func VerifyMobile(mobile string) bool {
	return regexp.MustCompile(RegMobile).MatchString(mobile)
}

// 邮箱验证
func VerifyEmail(email string) bool {
	return regexp.MustCompile(RegEmail).MatchString(email)
}
