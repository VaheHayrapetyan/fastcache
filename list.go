package fastcache

import "fmt"

type node struct {
	key   uint64
	value interface{}
	next  *node
}

type lCache struct {
	mutex     cleverMutex
	store     []*node
	length    int
	cacheSize uint64
	cacheBit  uint64
}

func newLCache(cacheBitCount uint64) (cache ICache) {

	c := &lCache{}

	c.cacheBit = cacheBitCount

	c.cacheSize = 1
	c.cacheSize <<= c.cacheBit
	c.store = make([]*node, c.cacheSize)
	c.length = 0
	cache = c

	return cache
}

func (c *lCache) Set(key uint64, value interface{}) {
	cKey := key & (c.cacheSize - 1)
	cLocalKey := key >> c.cacheBit

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.store[cKey] == nil {
		newNode := &node{
			key:   cLocalKey,
			value: value,
			next:  nil,
		}
		c.store[cKey] = newNode
		c.length++
		return
	}

	for n := c.store[cKey]; n != nil; n = n.next {
		if n.key == cLocalKey {
			n.value = value
			return
		}
	}

	newNode := &node{
		key:   cLocalKey,
		value: value,
		next:  c.store[cKey],
	}
	c.store[cKey] = newNode
	c.length++
}

func (c *lCache) Get(key uint64) (interface{}, bool) {
	cKey := key & (c.cacheSize - 1)
	cLocalKey := key >> c.cacheBit

	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	for n := c.store[cKey]; n != nil; n = n.next {
		if n.key == cLocalKey {
			return n.value, true
		}
	}

	return nil, false
}

func (c *lCache) Delete(key uint64) bool {
	cKey := key & (c.cacheSize - 1)
	cLocalKey := key >> c.cacheBit

	c.mutex.Lock()
	defer c.mutex.Unlock()

	p := c.store[cKey]
	if p == nil {
		return false
	} else if p.key == cLocalKey {
		c.store[cKey] = c.store[cKey].next
		c.length--
		return true
	}

	n := p.next

	for n != nil {
		if n.key == cLocalKey {
			p.next = n.next
			n = nil
			c.length--
			return true
		}
		p = n
		n = n.next
	}

	return false
}

func (c *lCache) Len() int {
	return c.length
}

func (c *lCache) Print() {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	s := "{"

	for cKey, n := range c.store {
		if n == nil {
			continue
		}
		for ; n != nil; n = n.next {
			key := (n.key << c.cacheBit) ^ uint64(cKey)
			s += fmt.Sprintf("%d: %v, ", key, n.value)
		}
	}

	if s != "{" {
		s = s[:(len(s) - 2)]
	}
	s += "}"
	fmt.Print(s)
}

func (c *lCache) Iterator() <-chan interface{} {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	res := make(chan interface{}, c.length)
	defer close(res)

	for _, n := range c.store {
		if n == nil {
			continue
		}

		for ; n != nil; n = n.next {
			res <- n.value
		}
	}

	return res
}

func (c *lCache) Range(f func(key uint64, value interface{}) bool) {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	for cKey, n := range c.store {
		if n == nil {
			continue
		}
		for ; n != nil; n = n.next {
			key := (n.key << c.cacheBit) ^ uint64(cKey)

			if !f(key, n.value) {
				return
			}
		}
	}
}

func (c *lCache) ToMap() map[uint64]interface{} {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	m := make(map[uint64]interface{}, c.length)

	for cKey, n := range c.store {
		if n == nil {
			continue
		}
		for ; n != nil; n = n.next {
			key := (n.key << c.cacheBit) ^ uint64(cKey)

			m[key] = n.value
		}
	}

	return m
}

func (c *lCache) TestPrintAllStructure() {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	for cKey, n := range c.store {
		fmt.Printf("%d: { ", cKey)

		if n == nil {
			fmt.Println("nil }")
			continue
		}

		for ; n != nil; n = n.next {
			fmt.Printf("%d: %v ", n.key, n.value)
		}
		fmt.Println("}")
	}
}
