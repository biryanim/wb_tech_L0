package lru_cache

import (
	"container/list"
	"github.com/biryanim/wb_tech_L0/internal/client/cache"
	"sync"
)

var _ cache.Client = (*Cache)(nil)

type Item struct {
	Key   string
	Value interface{}
}

type Cache struct {
	capacity int
	queue    *list.List
	mutex    *sync.RWMutex
	items    map[string]*list.Element
}

func New(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		queue:    list.New(),
		mutex:    new(sync.RWMutex),
		items:    make(map[string]*list.Element),
	}
}

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
