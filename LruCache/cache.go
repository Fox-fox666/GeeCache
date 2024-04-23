package LruCache

import (
	"container/list"
	"fmt"
	"log"
)

type Lru_Cache struct {
	maxbytes  int64                             //缓存最大容量
	nbytes    int64                             //当前存了多少容量
	cache     map[string]*cacheValue            //kv缓存
	queue     *list.List                        //LRU维护的队列,队列元素为key
	OnEvicted func(key string, value Cachedata) //回调函数，当有数据被删除或淘汰时触发
}

type Cachedata interface {
	Cap_bytes() int64 //返回该数据占的字节数
}

type cacheValue struct {
	data Cachedata     //缓存数据
	pos  *list.Element //该key在队列中的位置，方便维护队列
}

func NewCache(maxbytes int64, callback func(string, Cachedata)) *Lru_Cache {
	return &Lru_Cache{
		maxbytes:  maxbytes,
		nbytes:    0,
		cache:     make(map[string]*cacheValue),
		queue:     &list.List{},
		OnEvicted: callback,
	}
}

// 添加或者修改一个缓存
func (c *Lru_Cache) Add(key string, data Cachedata) {
	cap_bytes := data.Cap_bytes() + int64(len(key))
	if cap_bytes > c.maxbytes {
		log.Println(fmt.Sprintf("数据过大，存入失败:%s", key))
		return
	}
	//判断缓存中有没有这个数据，有就修改
	if value, ok := c.cache[key]; ok {
		c.cache[key].data = data                              //修改缓存
		c.queue.MoveToBack(value.pos)                         //放到队尾
		c.nbytes += data.Cap_bytes() - value.data.Cap_bytes() //修改目前内存大小
	} else {
		pos := c.queue.PushBack(key) //新增节点的队列节点位置放到队尾
		c.cache[key] = &cacheValue{  //加入缓存
			data: data,
			pos:  pos,
		}
		c.nbytes += cap_bytes
	}

	for c.maxbytes != 0 && c.nbytes > c.maxbytes {
		c.RemoveOneOldData()
	}
}

// 删除指定缓存
func (c *Lru_Cache) Del(key string) {
	if cv, ok := c.cache[key]; ok {
		c.nbytes -= cv.data.Cap_bytes() + int64(len(key))
		c.queue.Remove(cv.pos)
		delete(c.cache, key)
		if c.OnEvicted != nil {
			c.OnEvicted(key, cv.data)
		}
	} else {
		log.Println(fmt.Sprintf("Del失败，无缓存:%s", key))
	}
}

// 查找指定缓存
func (c *Lru_Cache) Get(key string) (Cachedata, bool) {
	//取值
	if cv, ok := c.cache[key]; ok {
		//维护队列，将该key移到队列尾部
		c.queue.MoveToBack(cv.pos)
		return cv.data, ok
	} else {
		log.Println(fmt.Sprintf("Get失败，无缓存:%s", key))
		return nil, false
	}
}

// 淘汰一个最近最久未使用缓存
func (c *Lru_Cache) RemoveOneOldData() {
	ele := c.queue.Front()
	if ele != nil {
		key := ele.Value.(string)
		c.Del(key)
	} else {
		log.Println("RemoveOneOldData失败,缓存为空")
	}
}
