package cache

// moves node to most recently used position
func (c *Cache[T]) moveToFront(node *node[T]) {
	if c.head == node {
		return
	}

	node.prev.next = node.next

	if node.next != nil {
		node.next.prev = node.prev
	}

	if c.tail == node {
		c.tail = node.prev
	}

	node.prev = nil
	node.next = c.head

	if c.head != nil {
		c.head.prev = node
	}
	c.head = node
}

// evicts least recently used item when at capacity
func (c *Cache[T]) evictLRU() {
	oldTailKey := c.tail.key
	c.removeNode(c.tail)
	delete(c.store, oldTailKey)
}
