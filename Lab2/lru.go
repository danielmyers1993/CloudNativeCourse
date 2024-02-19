package cache

import "errors"

type Cacher[K comparable, V any] interface {
	Get(key K) (value V, err error)
	Put(key K, value V) (err error)
}

// Concrete LRU cache
type lruCache[K comparable, V any] struct {
	size      int
	remaining int
	cache     map[K]V
	queue     []K
}

// Constructor
func NewCacher[K comparable, V any](size int) Cacher[K, V] {
	return &lruCache[K, V]{size: size, remaining: size, cache: make(map[K]V), queue: make([]K, 0)}
}

func (c *lruCache[K, V]) Get(key K) (value V, err error) {
	// Check if the key exists in the cache
	val, ok := c.cache[key]
	if !ok {
		return value, errors.New("key not found")
	}

	// Move the key to the end of the queue (mark as recently used)
	c.deleteFromQueue(key)
	c.queue = append(c.queue, key)

	return val, nil
}

func (c *lruCache[K, V]) Put(key K, value V) (err error) {
	// Your code here ...
	// Hint - Check if key already exists
	// Hint - Check capacity, and evict if needed
	// Hint - Add new key-value pair

	// Check if the key already exists in the cache
	_, exists := c.cache[key]

	// If the key does not exist and the cache is full, evict the least recently used item
	if !exists && c.remaining == 0 {
		// Evict the least recently used item (head of the queue)
		evictedKey := c.queue[0]
		delete(c.cache, evictedKey)
		c.deleteFromQueue(evictedKey)
		c.remaining++
	}

	// If the key exists, move it to the end of the queue
	if exists {
		c.deleteFromQueue(key)
	}

	// Add or update the key-value pair in the cache
	c.cache[key] = value
	c.queue = append(c.queue, key)

	// Decrement the remaining capacity of the cache
	if !exists {
		c.remaining--
	}

	return nil
}

// Helper method to delete all occurrences of a key from the queue
func (c *lruCache[K, V]) deleteFromQueue(key K) {
	newQueue := make([]K, 0, c.size)
	for _, k := range c.queue {
		if k != key {
			newQueue = append(newQueue, k)
		}
	}
	c.queue = newQueue
}
