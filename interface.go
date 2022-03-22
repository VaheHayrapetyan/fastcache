package fastcache

type ICache interface {
	Set(key uint64, value interface{})
	Get(key uint64) (interface{}, bool)
	Delete(key uint64) bool
	Len() int
	Iterator() <-chan interface{}
	Range(f func(key uint64, value interface{}) bool)
	ToMap() map[uint64]interface{}
	Print()
	TestPrintAllStructure()
}

func NewCache(bitCount uint64) (cache ICache, err error) {
	if bitCount >= 20 {
		return newSCache(bitCount)
	} else {
		return newLCache(bitCount)
	}
}
