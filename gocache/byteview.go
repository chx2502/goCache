package gocache

type ByteView struct {
	b []byte
}

func (bv ByteView) Len() int {
	return len(bv.b)
}
// 返回切片 b 的拷贝
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// 将数据以 string 的形式返回
func (bv ByteView) String() string {
	return string(bv.b)
}
