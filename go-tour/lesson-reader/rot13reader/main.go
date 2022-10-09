package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (reader rot13Reader) Read(b []byte) (int, error) {
	// 从自身reader中读取
	n, e := reader.r.Read(b)
	// 先判断是否还有数据
	if e == io.EOF {
		return n, e
	}
	// 若有数据，则转换之后返回
	for i := 0; i < n; i++ {
		if b[i]+13 <= 'z' {
			b[i] = b[i] + 13
		} else {
			b[i] = 'a' + b[i] + 13 - 'z' - 1
		}
	}
	return n, e
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
