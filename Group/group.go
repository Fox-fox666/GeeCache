package Group

import (
	"GeeCache/RWMutexCache"
	cachepb_pb "GeeCache/cachepb"
	"GeeCache/singleflight"
	"log"
	"sync"
)

//                           是
//接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
////                 |否                        是
////                 |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
////                                |  否
////                                |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶

type Group struct {
	name      string                    //命名
	getter    Getter                    //回调函数，当缓存中没有数据时，从数据源中获取数据放入缓存
	maincache RWMutexCache.RWMutexCache //并发缓存
	nodes     NodePicker
	loder     *singleflight.CallMap
}

var (
	groups = make(map[string]*Group) //管理各个缓存
	m      sync.RWMutex
)

func NewGroup(name string, cachebBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("无回调函数")
	}
	g := &Group{
		name:   name,
		getter: getter,
		maincache: RWMutexCache.RWMutexCache{
			MaxBytes: cachebBytes,
		},
	}
	m.Lock()
	defer m.Unlock()
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	m.RLock()
	defer m.RUnlock()
	return groups[name]
}

func (g *Group) Get(key string) (RWMutexCache.BytesData, error) {
	m.RLock()
	defer m.RUnlock()
	if v, ok := g.maincache.Get(key); ok { //缓存命中
		log.Println("Geecache hit")
		return v, nil
	}
	return g.load(key) //不命中，从别处加载
}

// 加载数据到缓存（本地或者别的节点）
func (g *Group) load(key string) (RWMutexCache.BytesData, error) {
	bytes, err := g.loder.Do(key, func() (any, error) {
		if g.nodes != nil {
			if nodegetter, ok := g.nodes.PickNode(key); ok {
				response, err := nodegetter.Get(&cachepb_pb.Request{
					Groupname: g.name,
					Key:       key,
				})
				if err != nil {
					return RWMutexCache.BytesData{}, err
				}
				return RWMutexCache.NewBytesData(response.GetData()), nil
			}
		}
		return g.getLocally(key)
	})
	if err != nil {
		return RWMutexCache.BytesData{}, err
	}
	return RWMutexCache.NewBytesData(bytes.([]byte)), nil
}

// 从本地，调用回调函数，拿到数据，并存入缓存
func (g *Group) getLocally(key string) (RWMutexCache.BytesData, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return RWMutexCache.BytesData{}, err
	}
	data := RWMutexCache.NewBytesData(bytes)
	g.Add(key, data)
	return data, nil
}

func (g *Group) Add(key string, value RWMutexCache.BytesData) {
	g.maincache.Add(key, value)
}

func (g *Group) RegisterNodes(node *HttpNode) {
	g.nodes = node
}
