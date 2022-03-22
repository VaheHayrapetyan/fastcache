package cache

import "sync"

type CleverMutex struct {
	readMutex sync.Mutex
	writeMutex sync.Mutex
}

func (cm *CleverMutex) ReadLock() {
	cm.readMutex.Lock()
}

func (cm *CleverMutex) WriteLock() {
	cm.writeMutex.Lock()
}

func (cm *CleverMutex) ReadUnlock() {
	cm.readMutex.Unlock()
}

func (cm *CleverMutex) WriteUnlock() {
	cm.writeMutex.Unlock()
}

func (cm *CleverMutex) Lock() {
	cm.readMutex.Lock()
	cm.writeMutex.Lock()
}

func (cm *CleverMutex) Unlock() {
	cm.readMutex.Unlock()
	cm.writeMutex.Unlock()
}