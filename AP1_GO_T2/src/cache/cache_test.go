package cache

import (
	"testing"
	"time"
)

func TestCacheCreateWithErr(t *testing.T) {
	cache, err := NewCache[int](0)

	if cache != nil {
		t.Error("cache != nil")
	}

	if err == nil {
		t.Fatal("error got nil")
	}

	expectedError := "capacity must be positive"
	if err.Error() != expectedError {
		t.Errorf("expected: %s, got: %s", expectedError, err.Error())
	}
}

func TestCacheSetBasic(t *testing.T) {
	cache, _ := NewCache[float64](5)

	cache.Set("first", 12.3)
	cache.Set("second", 34.5)
	cache.Set("third", 67.8)

	// third(67.8) -> second(34.5) -> first(12.3)
	// head=third, tail=first

	if cache.head == nil {
		t.Fatal("the head was expected not nil")
	}
	if cache.head.key != "third" {
		t.Errorf("expected: third, got: %s", cache.head.key)
	}
	if cache.head.value != 67.8 {
		t.Errorf("expected: 67.8, got :%v", cache.head.value)
	}

	if cache.tail == nil {
		t.Fatal("the tail was expected not nil")
	}
	if cache.tail.key != "first" {
		t.Errorf("expected: first, got: %s", cache.tail.key)
	}
	if cache.tail.value != 12.3 {
		t.Errorf("expected: 12.3, got: %v", cache.tail.value)
	}

	if cache.head.next == nil {
		t.Fatal("Head should have next node (second)")
	}
	if cache.head.prev != nil {
		t.Error("Head should not have prev node")
	}

	secondNode := cache.head.next
	if secondNode.key != "second" {
		t.Errorf("expected: second, got: %s", secondNode.key)
	}
	if secondNode.value != 34.5 {
		t.Errorf("expected: 34.5, got: %v", secondNode.value)
	}
	if secondNode.prev != cache.head {
		t.Error("Second node should point to head as prev")
	}
	if secondNode.next == nil {
		t.Fatal("Second node should have next node (first)")
	}

	firstNode := secondNode.next
	if firstNode != cache.tail {
		t.Error("First node should be tail")
	}
	if firstNode.key != "first" {
		t.Errorf("expected: first, got: %s", firstNode.key)
	}
	if firstNode.value != 12.3 {
		t.Errorf("expected: 12.3, got: %v", firstNode.value)
	}
	if firstNode.prev != secondNode {
		t.Error("First node should point to second as prev")
	}
	if firstNode.next != nil {
		t.Error("First node should not have next node")
	}

}

func TestCacheSetWithOverflowCapacity(t *testing.T) {
	cache, _ := NewCache[string](2)

	cache.Set("first", "aaa")
	cache.Set("second", "bbb")
	cache.Set("third", "ccc")

	// third(ccc) -> second(bbb)
	// head=third, tail=second

	if _, exists := cache.store["first"]; exists {
		t.Error("Store should not contain evicted element first")
	}
	if _, exists := cache.store["second"]; !exists {
		t.Error("Store should contain second")
	}
	if _, exists := cache.store["third"]; !exists {
		t.Error("Store should contain third")
	}

	if len(cache.store) != 2 {
		t.Errorf("expected: 2 elements, got: %d", len(cache.store))
	}

	if cache.head.key != "third" {
		t.Errorf("expected: third, got: %s", cache.head.key)
	}
	if cache.head.value != "ccc" {
		t.Errorf("expected: ccc, got:%s", cache.head.value)
	}
	if cache.head.prev != nil {
		t.Error("Head should not have prev")
	}

	if cache.tail.key != "second" {
		t.Errorf("expected: second, got: %s", cache.tail.key)
	}
	if cache.tail.value != "bbb" {
		t.Errorf("expected: bbb, got: %s", cache.tail.value)
	}
	if cache.tail.next != nil {
		t.Error("Tail should not have next")
	}

	if cache.head.next != cache.tail {
		t.Error("Head should point to tail as next")
	}
	if cache.tail.prev != cache.head {
		t.Error("Tail should point to head as prev")
	}
}

func TestCacheTTLExpiration(t *testing.T) {
	// need to set tickerPause for 1 second
	// need to set cacheTTL for 10 seconds

	cache, _ := NewCache[int](3)

	cache.Set("first", 100)
	time.Sleep(5 * time.Second)
	cache.Set("second", 200)
	time.Sleep(7 * time.Second)

	// first: 5 + 7 = 12 seconds > TTL 10 seconds - delete
	// second: 7 seconds < TTL 10 seconds - must stay

	if _, exists := cache.store["first"]; exists {
		t.Error("first must not be present")
	}

	if _, exists := cache.store["second"]; !exists {
		t.Error("second should still be present")
	}

	if cache.head == nil || cache.tail == nil {
		t.Fatal("Head and tail should not be nil")
	}
	if cache.head.key != "second" {
		t.Errorf("expected: second, got: %s", cache.head.key)
	}
	if cache.tail.key != "second" {
		t.Errorf("expected: second, got: %s", cache.tail.key)
	}

	if cache.head != cache.tail {
		t.Error("Head and tail should point to the same node")
	}

	if cache.head.prev != nil {
		t.Error("Head should not have prev")
	}
	if cache.head.next != nil {
		t.Error("Head should not have next")
	}
	if cache.tail.prev != nil {
		t.Error("Tail should not have prev")
	}
	if cache.tail.next != nil {
		t.Error("Tail should not have next")
	}
}

func TestCacheGetNonExistentKey(t *testing.T) {
	cache, _ := NewCache[int](3)

	cache.Set("first", 100)
	cache.Set("second", 200)
	cache.Set("third", 300)

	value, found := cache.Get("nonexistent")
	if found {
		t.Error("Expected false for non-existent key")
	}
	if value != 0 {
		t.Errorf("expected: zero value, got %d", value)
	}
}

func TestCacheGetUpdatesOrder(t *testing.T) {
	cache, _ := NewCache[int](3)

	cache.Set("first", 100)
	cache.Set("second", 200)
	cache.Set("third", 300)
	// third(300) -> second(200) -> first(100)
	// head=third, tail=first

	value, found := cache.Get("first")
	if !found || value != 100 {
		t.Error("Failed to get first")
	}

	// first(100) -> third(300) -> second(200)
	// head=first, tail=second

	if cache.head.key != "first" {
		t.Errorf("expected:: first, got: %s", cache.head.key)
	}
	if cache.head.next.key != "third" {
		t.Errorf("expected:: third, got: %s", cache.head.next.key)
	}
	if cache.head.next.next.key != "second" {
		t.Errorf("expected: second, got: %s", cache.head.next.next.key)
	}
	if cache.tail.key != "second" {
		t.Errorf("expected:: second, got: %s", cache.tail.key)
	}

	value, found = cache.Get("second")
	if !found || value != 200 {
		t.Error("Failed to get second")
	}

	// second(200) -> first(100) -> third(300)
	// head=second, tail=third

	if cache.head.key != "second" {
		t.Errorf("expected: second, got: %s", cache.head.key)
	}
	if cache.head.next.key != "first" {
		t.Errorf("expected: first, got: %s", cache.head.next.key)
	}
	if cache.head.next.next.key != "third" {
		t.Errorf("expected: third, got: %s", cache.head.next.next.key)
	}
	if cache.tail.key != "third" {
		t.Errorf("expected: third, got: %s", cache.tail.key)
	}

	value, found = cache.Get("first")
	if !found || value != 100 {
		t.Error("Failed to get: first again")
	}

	// first(100) -> second(200) -> third(300)
	// head=first, tail=third

	if cache.head.key != "first" {
		t.Errorf("expected: first, got: %s", cache.head.key)
	}
	if cache.head.next.key != "second" {
		t.Errorf("expected: second, got: %s", cache.head.next.key)
	}
	if cache.head.next.next.key != "third" {
		t.Errorf("expected: third, got: %s", cache.head.next.next.key)
	}
	if cache.tail.key != "third" {
		t.Errorf("expected: third, got: %s", cache.tail.key)
	}

	if cache.head.prev != nil {
		t.Error("Head should not have prev")
	}
	if cache.head.next.prev != cache.head {
		t.Error("Second node should point to head as prev")
	}
	if cache.tail.next != nil {
		t.Error("Tail should not have next")
	}
	if cache.tail.prev != cache.head.next {
		t.Error("Tail should point to second node as prev")
	}
}

func TestCacheSetUpdatesExistingKey(t *testing.T) {
	cache, _ := NewCache[bool](3)

	cache.Set("first", true)
	cache.Set("second", false)
	cache.Set("third", true)

	// third(true) -> second(false) -> first(true)
	// head=third, tail=first

	cache.Set("first", false)

	// first(false) -> third(true) -> second(false)
	// head=first, tail=second

	if cache.head.key != "first" {
		t.Errorf("expected: first , got: %s", cache.head.key)
	}
	if cache.head.value != false {
		t.Errorf("expected: false, got %v", cache.head.value)
	}

	if cache.head.next.key != "third" {
		t.Errorf("expected: third, got: %s", cache.head.next.key)
	}
	if cache.head.next.next.key != "second" {
		t.Errorf("expected: second, got: %s", cache.head.next.next.key)
	}
	if cache.tail.key != "second" {
		t.Errorf("expected: second, got: %s", cache.tail.key)
	}

	if cache.head.prev != nil {
		t.Error("Head should not have prev")
	}
	if cache.head.next.prev != cache.head {
		t.Error("Third node should point to head as prev")
	}
	if cache.tail.next != nil {
		t.Error("Tail should not have next")
	}
	if cache.tail.prev != cache.head.next {
		t.Error("Tail should point to third node as prev")
	}

	if len(cache.store) != 3 {
		t.Errorf("expected: 3, got: %d", len(cache.store))
	}
}

func TestCacheEvictionAfterMultipleGets(t *testing.T) {
	cache, _ := NewCache[int](3)

	cache.Set("a", 100)
	cache.Set("b", 200)
	cache.Set("c", 300)

	cache.Get("a") // A -> C -> B
	cache.Get("b") // B -> A -> C
	cache.Get("a") // A -> B -> C
	cache.Get("c") // C -> A -> B

	if cache.head.key != "c" {
		t.Errorf("expected: C, got: %s", cache.head.key)
	}
	if cache.head.next.key != "a" {
		t.Errorf("expected: A, got: %s", cache.head.next.key)
	}
	if cache.tail.key != "b" {
		t.Errorf("expected: B, got: %s", cache.tail.key)
	}

	cache.Set("d", 400)

	_, found := cache.Get("b")
	if found {
		t.Error("B should be evicted")
	}

	if cache.head.key != "d" {
		t.Errorf("expected: D, got: %s", cache.head.key)
	}
	if cache.head.next.key != "c" {
		t.Errorf("expected: C, got: %s", cache.head.next.key)
	}
	if cache.tail.key != "a" {
		t.Errorf("expected: A, got: %s", cache.tail.key)
	}

	val, _ := cache.Get("d")
	if val != 400 {
		t.Errorf("expected: 400, got: %d", val)
	}
	val, _ = cache.Get("c")
	if val != 300 {
		t.Errorf("expected: 300, got: %d", val)
	}
	val, _ = cache.Get("a")
	if val != 100 {
		t.Errorf("expected: 100, got: %d", val)
	}

	if cache.head.prev != nil {
		t.Error("Head should not have prev")
	}
	if cache.head.next.prev != cache.head {
		t.Error("C should point to D as prev")
	}
	if cache.tail.next != nil {
		t.Error("Tail should not have next")
	}
	if cache.tail.prev != cache.head.next {
		t.Error("A should point to C as prev")
	}

	if len(cache.store) != 3 {
		t.Errorf("expected: 3, got %d", len(cache.store))
	}
}

func TestCacheClearAndReuse(t *testing.T) {
	cache, _ := NewCache[int](3)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Clear()

	if cache.head != nil || cache.tail != nil || len(cache.store) != 0 {
		t.Error("Cache should be empty after Clear")
	}

	cache.Set("c", 3)
	val, found := cache.Get("c")
	if !found || val != 3 {
		t.Error("Should work after Clear")
	}
}

func TestCacheCapacityOne(t *testing.T) {
	cache, _ := NewCache[int](1)

	cache.Set("first", 1)
	cache.Set("second", 2)

	if _, found := cache.Get("first"); found {
		t.Error("First should be evicted")
	}
	if cache.head != cache.tail {
		t.Error("Head and tail should be same for capacity 1")
	}
}

func TestCacheTTLWithActiveUsage(t *testing.T) {
	cache, _ := NewCache[int](3)

	cache.Set("active", 100)

	for i := range 5 {
		time.Sleep(2 * time.Second)
		cache.Set("active", 100+i)
	}

	time.Sleep(5 * time.Second)

	if _, exists := cache.store["active"]; !exists {
		t.Error("Active element should not expire")
	}
}

func TestCacheTTLWithActiveUsage2(t *testing.T) {
	cache, _ := NewCache[int](3)

	cache.Set("active", 100)

	for range 5 {
		time.Sleep(2 * time.Second)
		cache.Get("active")
	}

	time.Sleep(5 * time.Second)

	if _, exists := cache.store["active"]; !exists {
		t.Error("Active element should not expire")
	}
}

func TestCacheEmptyOperations(t *testing.T) {
	cache, _ := NewCache[int](3)

	cache.Clear()
}
