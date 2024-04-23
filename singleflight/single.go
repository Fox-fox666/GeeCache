package singleflight

import (
	"sync"
)

type call struct {
	w   *sync.WaitGroup
	val any
	err error
}

type CallMap struct {
	mu      sync.Mutex
	callmap map[string]*call
}

func (c *CallMap) Do(key string, fn func() (any, error)) (any, error) {
	c.mu.Lock()
	if v, ok := c.callmap[key]; ok {
		c.mu.Unlock()
		v.w.Wait()
		return v.val, v.err
	}
	cc := new(call)
	cc.w.Add(1)
	c.callmap[key] = cc
	c.mu.Unlock()

	cc.val, cc.err = fn()
	cc.w.Done()

	c.mu.Lock()
	delete(c.callmap, key)
	c.mu.Unlock()

	return cc.val, cc.err
}
