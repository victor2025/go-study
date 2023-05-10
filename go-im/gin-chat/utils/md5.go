package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	strPre := h.Sum(nil)
	return hex.EncodeToString(strPre)
}

func Md5EncodeUpper(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// md5+salt生成密码
func MakePassword(plainPwd, salt string) string {
	return Md5Encode(plainPwd + salt)
}

// 判断密码是否正确
func ValidatePassword(plainPwd, salt, password string) bool {
	return Md5Encode(plainPwd+salt) == password
}

func GenSalt() string {
	return fmt.Sprintf("%06d", rand.Uint32())
}
