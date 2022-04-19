# fastcache

Hello. Here you will get acquainted with the fastcache library. It will allow you to use a key-value structure to store your data.        
It's like a map. The key here can be an arbitrary value of type int, and the value can be of any type.  
Most importantly, fastcache is twice as fast as golang map.

Functions:

    func NewCache(bitCount uint64) ICache
    func (c ICache) Set(key uint64, value interface{})
	func (c ICache) Get(key uint64) (interface{}, bool)
	func (c ICache) Delete(key uint64) bool
	func (c ICache) Len() int
	func (c ICache) Iterator() <-chan interface{}
	func (c ICache) Range(f func(key uint64, value interface{}) bool)
	func (c ICache) ToMap() map[uint64]interface{}
	func (c ICache) Print()
	func (c ICache) TestPrintAllStructure()

Now let's start from the beginning․     
You must use NewCache function with which you will create the ICache interface. You can then store the key-value data, get it, delete it, and perform other actions.

    func NewCache(bitCount uint64) ICache

Let's try to understand the structure of our cache. If you are familiar with the structure of the map, you can easily understand it․    
First, a vertical array is created, the size of which is determined by the bitCount value passed to the NewCache function: the size of this array or cache capacity is 2^bitCount.  

Now let's understand how key-value data is stored. This also depends on the value of the bitCount.  
1․ If the bitCount less than 20, the key-values are stored in a linked list created below the vertical array according to the index: each node in this list stores one key-value.           
2․ If the bitCount is greater than or equal to 20, the key-values are stored in a slice created below the vertical array according to the index։ the keys are stored at even indexes, and the value corresponding to the key written at that even index is stored at the odd index following that even index. And with some bit operations, it is determined under which index of the main slice to store the key-value.    

Each key-value pair is stored in a vertical array in a slice or list at a specific index. The value of this index is determined by the bit operations performed on the key and the capacity of the vertical array.  
Each key is stored in the cache not in its real value, and the value of each of them is changed in accordance with some bit operations, after which the result is stored in the cache as a key. 

After using the NewCache function, you will have access to the ICache interface:

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

Let's take a look at the purpose of each of these functions.

    func (c ICache) Set(key uint64, value interface{})

The Set function adds a new key-value pair to the cache. If the key already exists in cache, the corresponding value will be updated.

    func (c ICache) Get(key uint64) (interface{}, bool)

The Get function returns the value corresponding to the passed key and true if the key exists in the cache and returns nil, false if the key does not exist.

	func (c ICache) Delete(key uint64) bool

The Delete function deletes the key-value pair corresponding to the passed key and returns true if the key exists in the cache. If the key does not exist, the function returns false.

	func (c ICache) Len() int

The Len function returns the length of the cache, which corresponds to the number of key-value pairs.

	func (c ICache) Iterator() <-chan interface{}

The Iterator function returns a read-only channel containing all cache values.

	func (c ICache) Range(f func(key uint64, value interface{}) bool)

The Range function ranges over the cache, passes each key-value pair to the function f that is passed to the Range function, and exits if the f function returns false. The function f takes (key, value), does what you want, and returns a boolean value.

	func (c ICache) ToMap() map[uint64]interface{}

The ToMap function converts the cache to a map[uint64]interface{} and returns it.

	func (c ICache) Print()

The Print function prints all the cache key-value pairs in the order they are actually stored in the structure.

	func (c ICache) TestPrintAllStructure()

The TestPrintAllStructure function prints the cache structure where the keys correspond not to the values actually passed as keys, but to the values actually stored in the cache.

The map.go file is for testing purposes only...

That's all !!! Thank you for your attention.