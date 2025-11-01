package lrucache

import (
	"container/list"
	"sync"

	"github.com/biryanim/wb_tech_L0/internal/client/cache"
)

var _ cache.Client = (*Cache)(nil)

// Item represents a single cache entry with a key-value pair.
type Item struct {
	Key   string
	Value interface{}
}

// Cache implements an LRU (Least Recently Used) cache with thread-safe operations.
type Cache struct {
	capacity int
	queue    *list.List
	mutex    *sync.RWMutex
	items    map[string]*list.Element
}

// New creates and returns a new Cache instance with the specified capacity.
func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		queue:    list.New(),
		mutex:    new(sync.RWMutex),
		items:    make(map[string]*list.Element),
	}
}

// Set stores or updates a value in the cache with the given key.
func (c *Cache) Set(key string, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exists := c.items[key]; exists {
		c.queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		return true
	}

	if c.queue.Len() == c.capacity {
		c.clear()
	}

	item := &Item{
		Key:   key,
		Value: value,
	}

	element := c.queue.PushFront(item)
	c.items[item.Key] = element

	return true
}

// Get retrieves a value from the cache by key and moves it to the front as most recently used.
func (c *Cache) Get(key string) interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	element, exists := c.items[key]
	if exists == false {
		return nil
	}

	c.queue.MoveToFront(element)
	return element.Value.(*Item).Value
}

// Remove deletes a value from the cache by key and returns true if the key existed.
func (c *Cache) Remove(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if val, found := c.items[key]; found {
		c.deleteItem(val)
	}

	return true
}

func (c *Cache) clear() {
	if element := c.queue.Back(); element != nil {
		c.deleteItem(element)
	}
}

func (c *Cache) deleteItem(element *list.Element) {
	item := c.queue.Remove(element).(*Item)
	delete(c.items, item.Key)
}
