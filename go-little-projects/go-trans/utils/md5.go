package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func Md5Encode(data []byte) string {
	h := md5.New()
	h.Write([]byte(data))
	strPre := h.Sum(nil)
	return hex.EncodeToString(strPre)
}

func Md5EncodeUpper(data []byte) string {
	return strings.ToUpper(Md5Encode(data))
}
