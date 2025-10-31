package cache

import (
	"errors"
	"sync"
	"time"
)

var (
	tickerPause = 2 * time.Minute //set less for tests (1 second)
	cacheTTL    = 1 * time.Hour   //set less for tests (10 seconds)
)

type node[T any] struct {
	key         string    // key for map lookup and deletion
	value       T         // cached data
	prev        *node[T]  // previous node in LRU list (newer element)
	next        *node[T]  // next node in LRU list (older element)
	storageTime time.Time // timestamp
}

type Cache[T any] struct {
	store    map[string]*node[T] // key to node mapping for direct access
	head     *node[T]            // most recently used node
	tail     *node[T]            // least recently used node
	capacity int                 // maximum number of items in cache
	mu       sync.RWMutex
}

func NewCache[T any](capacity int) (*Cache[T], error) {
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}

	cache := &Cache[T]{
		store:    make(map[string]*node[T]),
		head:     nil,
		tail:     nil,
		capacity: capacity,
	}
	go cache.startCleanup() // TTL cleanup goroutine
	return cache, nil
}

func (c *Cache[T]) Set(key string, item T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	existingNode, exists := c.store[key]
	if exists {
		existingNode.value = item
		c.moveToFront(existingNode)
		existingNode.storageTime = time.Now()
		return
	}

	newNode := &node[T]{
		key:         key,
		value:       item,
		prev:        nil,    // will become the new head
		next:        c.head, // link to current head
		storageTime: time.Now(),
	}

	c.store[key] = newNode

	if c.head != nil {
		c.head.prev = newNode
	}
	c.head = newNode

	if c.tail == nil {
		c.tail = newNode
	}

	if len(c.store) > c.capacity {
		c.evictLRU()
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	node, exists := c.store[key]
	c.mu.RUnlock()

	if !exists {
		var zero T
		return zero, false
	}

	c.mu.Lock()
	c.moveToFront(node)
	value := node.value
	node.storageTime = time.Now()
	c.mu.Unlock()
	return value, true
}

func (c *Cache[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store = make(map[string]*node[T])
	c.head = nil
	c.tail = nil
}
