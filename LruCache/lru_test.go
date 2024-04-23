package LruCache

import (
	"fmt"
	"testing"
)

type String string

func (d String) Cap_bytes() int64 {
	return int64(len(d))
}

func TestGet(t *testing.T) {
	lru := NewCache(int64(1000), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); ok {
		fmt.Println(v)
	}
	lru.Get("key2")
}

func TestRemoveOneOldData(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap_bytes := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap_bytes), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	lru.Get(k1)
}

func TestDel(t *testing.T) {
	k1, k2 := "key1", "key2"
	v1, v2 := "value1", "value2"
	cap_bytes := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap_bytes), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))

	lru.Del(k1)
	lru.Del(k1)
}

func TestUpdate(t *testing.T) {
	k1, k2 := "key1", "key2"
	v1, v2 := "value1", "value2"
	cap_bytes := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap_bytes), nil)
	lru.Add(k1, String(v1))
	if v, ok := lru.Get(k1); ok {
		fmt.Println(v)
	}
	lru.Add(k1, String("dadasdarqwradasdasdadasdsadasdasdasdasdasdasdasdssdasd"))
	if v, ok := lru.Get(k1); ok {
		fmt.Println(v)
	}
	lru.Add(k1, String("dad"))
	if v, ok := lru.Get(k1); ok {
		fmt.Println(v)
	}
}
