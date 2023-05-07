package utils

import (
	"testing"
)

func TestBytes2Uint64(t *testing.T) {
	var n1 byte = 0
	var n2 byte = 40
	bytes := []byte{n1, n2}
	expect := uint64(n1)<<8 + uint64(n2)
	t.Logf("expect: %d", expect)
	if res := Bytes2Uint64(bytes); res != expect {
		t.Fatalf("expect: %v, get: %v\n", 1<<8+1, res)
	}
}
