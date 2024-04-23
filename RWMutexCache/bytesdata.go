package RWMutexCache

// 只读的缓存数据
type BytesData struct {
	data []byte
}

func NewBytesData(data []byte) BytesData {
	return BytesData{data: data}
}

func (b BytesData) ToString() string {
	return string(b.data)
}

func (b BytesData) Cap_bytes() int64 {
	return int64(len(b.data))
}

func (b BytesData) ToSilce() []byte {
	dst := make([]byte, len(b.data))
	copy(dst, b.data)
	return dst
}
