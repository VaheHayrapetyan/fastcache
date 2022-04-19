package fastcache

import "sync"

type cleverMutex struct {
	readMutex  sync.Mutex
	writeMutex sync.Mutex
}

func (cm *cleverMutex) ReadLock() {
	cm.readMutex.Lock()
}

func (cm *cleverMutex) WriteLock() {
	cm.writeMutex.Lock()
}

func (cm *cleverMutex) ReadUnlock() {
	cm.readMutex.Unlock()
}

func (cm *cleverMutex) WriteUnlock() {
	cm.writeMutex.Unlock()
}

func (cm *cleverMutex) Lock() {
	cm.readMutex.Lock()
	cm.writeMutex.Lock()
}

func (cm *cleverMutex) Unlock() {
	cm.readMutex.Unlock()
	cm.writeMutex.Unlock()
}
