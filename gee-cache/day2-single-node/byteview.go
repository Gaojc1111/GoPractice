package day2_single_node

// 只读数据结构 ByteView 用来表示缓存值
type Byteview struct {
	b []byte
}

// 被缓存对象必须实现 Value 接口
func (v Byteview) Len() int {
	return len(v.b)
}

func (v Byteview) ByteSlice() []byte {
	return cloneBytes(v.b) //返回一个拷贝，防止缓存值被外部程序修改。
	// b是只读的，所以不能更改，不能直接返回
}

func (v Byteview) String() string {
	return string(v.b)
}
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
