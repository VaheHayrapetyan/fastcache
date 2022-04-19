package fastcache

import (
	"fmt"
)

type sCache struct {
	mutex     cleverMutex
	store     [][]interface{}
	length    int
	cacheSize uint64
	cacheBit  uint64
}

func newSCache(cacheBitCount uint64) (cache ICache) {

	c := &sCache{}

	c.cacheBit = cacheBitCount

	c.cacheSize = 1
	c.cacheSize <<= c.cacheBit
	c.store = make([][]interface{}, c.cacheSize)
	c.length = 0
	cache = c

	return cache
}

func (c *sCache) Set(key uint64, value interface{}) {
	cKey := key & (c.cacheSize - 1)
	cLocalKey := key >> c.cacheBit

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.store[cKey] == nil {
		c.store[cKey] = make([]interface{}, 2)
		c.store[cKey][uint64(0)] = cLocalKey
		c.store[cKey][uint64(1)] = value
		c.length++
	} else {
		lenLocal := uint64(len(c.store[cKey]))
		for i := uint64(0); i < lenLocal; i += 2 {
			if c.store[cKey][i] == cLocalKey {
				c.store[cKey][i+1] = value
				return
			}
		}

		c.store[cKey] = append(c.store[cKey], cLocalKey)
		c.store[cKey] = append(c.store[cKey], value)
		c.length++
	}
}

func (c *sCache) Get(key uint64) (interface{}, bool) {
	cKey := key & (c.cacheSize - 1)
	cLocalKey := key >> c.cacheBit

	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	if c.store[cKey] == nil {
		return nil, false
	}

	lenLocal := uint64(len(c.store[cKey]))

	for i := uint64(0); i < lenLocal; i += 2 {
		if c.store[cKey][i] == cLocalKey {
			return c.store[cKey][i+1], true
		}
	}

	return nil, false
}

func (c *sCache) Delete(key uint64) bool {
	cKey := key & (c.cacheSize - 1)
	cLocalKey := key >> c.cacheBit

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.store[cKey] == nil {
		return false
	}

	lenLocal := uint64(len(c.store[cKey]))
	if lenLocal == 2 && c.store[cKey][0] == cLocalKey {
		c.store[cKey] = nil
		c.length--
		return true
	}

	for i := uint64(0); i < lenLocal; i += 2 {
		if c.store[cKey][i] == cLocalKey {
			c.store[cKey] = append(c.store[cKey][:i], c.store[cKey][(i+2):]...)
			c.length--
			return true
		}
	}

	return false
}

func (c *sCache) Len() int {
	return c.length
}

func (c *sCache) Print() {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	s := "{"

	for cKey, vArr := range c.store {
		if vArr == nil {
			continue
		}
		lenVArr := uint64(len(vArr))
		for i := uint64(0); i < lenVArr; i += 2 {
			cLocalKey := vArr[i]
			value := vArr[i+1]
			key := (cLocalKey.(uint64) << c.cacheBit) ^ uint64(cKey)
			s += fmt.Sprintf("%v: %v, ", key, value)
		}
	}
	if s != "{" {
		s = s[:(len(s) - 2)]
	}
	s += "}"
	fmt.Print(s)
}

func (c *sCache) Iterator() <-chan interface{} {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	res := make(chan interface{}, c.length)
	defer close(res)

	for _, vArr := range c.store {
		if vArr == nil {
			continue
		}
		lenVArr := uint64(len(vArr))
		for i := uint64(0); i < lenVArr; i += 2 {
			value := vArr[i+1]
			res <- value
		}
	}

	return res
}

func (c *sCache) Range(f func(key uint64, value interface{}) bool) {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	for cKey, vArr := range c.store {
		if vArr == nil {
			continue
		}
		lenVArr := uint64(len(vArr))
		for i := uint64(0); i < lenVArr; i += 2 {
			cLocalKey := vArr[i]
			value := vArr[i+1]
			key := (cLocalKey.(uint64) << c.cacheBit) ^ uint64(cKey)
			if !f(key, value) {
				return
			}
		}
	}
}

func (c *sCache) ToMap() map[uint64]interface{} {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	m := make(map[uint64]interface{}, c.length)

	for cKey, vArr := range c.store {
		if vArr == nil {
			continue
		}
		lenVArr := uint64(len(vArr))
		for i := uint64(0); i < lenVArr; i += 2 {
			cLocalKey := vArr[i]
			value := vArr[i+1]
			key := (cLocalKey.(uint64) << c.cacheBit) ^ uint64(cKey)
			m[key] = value
		}
	}
	return m
}

func (c *sCache) TestPrintAllStructure() {
	c.mutex.WriteLock()
	defer c.mutex.WriteUnlock()

	for cKey, vArr := range c.store {
		fmt.Printf("%d: { ", cKey)

		if vArr == nil {
			fmt.Println("nil }")
			continue
		}

		lenVArr := uint64(len(vArr))
		for i := uint64(0); i < lenVArr; i += 2 {
			fmt.Printf("%v: %v ", vArr[i], vArr[i+1])
		}

		fmt.Println("}")
	}
}
