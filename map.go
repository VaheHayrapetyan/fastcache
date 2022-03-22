package cache

import (
	"fmt"
)

type mCache struct {
	mutex CleverMutex
	store map[uint64]interface{}
	cacheSize uint64
}

func newMCache(cacheSize uint64) (cache ICache, err error) {

	c := &mCache{}

	c.cacheSize = cacheSize

	c.store = make(map[uint64]interface{}, c.cacheSize)
	cache = c
	return cache, nil
}

func (c *mCache) Set(key uint64, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.store[key] = value
}

func (c *mCache) Get(key uint64) (value interface{}, ok bool) {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	value, ok = c.store[key]
	return
}

func (c *mCache) Delete(key uint64) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.store[key]; ok {
		delete(c.store, key)
		return true
	}

	return false
}

func (c *mCache) Len() int {
	return len(c.store)
}

func (c *mCache) Print() {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	fmt.Print(c.store)
}

func (c *mCache) Iterator() <- chan interface{} {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	res := make(chan interface{}, len(c.store))
	defer close(res)

	for k, _ := range c.store {
		res <- c.store[k]
	}

	return res
}

func (c *mCache) Range(f func(key uint64, value interface{}) bool) {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	for k, _ := range c.store {
		if !f(k, c.store[k]) {
			return
		}
	}
}

func (c *mCache) ToMap() (m map[uint64]interface{}) {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()
	m = c.store
	return
}

func (c *mCache) TestPrintAllStructure() {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()
	fmt.Println(c.store)
}