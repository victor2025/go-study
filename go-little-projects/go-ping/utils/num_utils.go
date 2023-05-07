package utils

func Bytes2Uint64(data []byte) uint64 {
	if len(data) == 0 {
		return 0
	}
	var res uint64
	for i := 0; i < len(data)-1; i++ {
		res += uint64(data[i]) << 8
	}
	res += uint64(data[len(data)-1])
	return res
}
