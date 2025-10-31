package cache

import "time"

// starts background TTL cleanup goroutine
func (c *Cache[T]) startCleanup() {
	ticker := time.NewTicker(tickerPause)
	defer ticker.Stop()
	for {
		<-ticker.C
		c.cleanupExpired()
	}
}

func (c *Cache[T]) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	expiredKeys := make([]string, 0)

	for key, node := range c.store {
		if now.Sub(node.storageTime) > cacheTTL {
			expiredKeys = append(expiredKeys, key)
		}
	}

	for _, key := range expiredKeys {
		c.removeNode(c.store[key])
		delete(c.store, key)
	}
}
