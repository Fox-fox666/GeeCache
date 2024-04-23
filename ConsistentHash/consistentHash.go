package ConsistentHash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

type NodeMap struct {
	hash      Hash           //hash算法
	replices  int            //节点扩展倍数
	keys      []int          //所有节点分布的key，从小至大排序
	keyToNode map[int]string //每一个key对应的节点
}

func NewNodeMap(hash Hash, replices int) *NodeMap {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}
	return &NodeMap{
		hash:      hash,
		replices:  replices,
		keys:      make([]int, 0),
		keyToNode: make(map[int]string),
	}
}

func (m *NodeMap) Add(nodenames ...string) {
	for _, nodename := range nodenames {
		for j := 0; j < m.replices; j++ {
			key := int(m.hash([]byte(strconv.Itoa(j) + nodename)))
			m.keys = append(m.keys, key)
			m.keyToNode[key] = nodename
		}
	}
	sort.Ints(m.keys)
}

func (m *NodeMap) Get(s string) string {
	if len(s) == 0 {
		return ""
	}

	key := int(m.hash([]byte(s)))
	pos := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= key
	})
	pos = pos % len(m.keys)
	return m.keyToNode[m.keys[pos]]
}
