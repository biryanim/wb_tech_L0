package lru_cache

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestLRU_SetExistingElementToFullCache(t *testing.T) {
	lru := New(3)
	lru.Set("someKey1", 8)
	lru.Set("someKey2", "23")
	emptyMap := make(map[string]int)
	lru.Set("someKey3", emptyMap)

	lru.Set("someKey1", 10)

	frontItem := lru.queue.Front().Value.(*Item)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey1", frontItem.Key)
	assert.Equal(t, 10, frontItem.Value)
	assert.Equal(t, "someKey2", backItem.Key)
	assert.Equal(t, "23", backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_SetExistingElementToNotFullCache(t *testing.T) {
	lru := New(3)
	lru.Set("someKey1", 8)
	lru.Set("someKey2", "23")

	lru.Set("someKey1", 10)

	frontItem := lru.queue.Front().Value.(*Item)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey1", frontItem.Key)
	assert.Equal(t, 10, frontItem.Value)
	assert.Equal(t, "someKey2", backItem.Key)
	assert.Equal(t, "23", backItem.Value)
	assert.Equal(t, 2, lru.queue.Len())
}

func TestLRU_SetNewElementToFullCache(t *testing.T) {
	lru := New(3)
	lru.Set("someKey1", 8)
	lru.Set("someKey2", "23")
	lru.Set("someKey3", Item{"key", 7})

	lru.Set("someKey4", 99)

	frontItem := lru.queue.Front().Value.(*Item)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey4", frontItem.Key)
	assert.Equal(t, 99, frontItem.Value)
	assert.Equal(t, "someKey2", backItem.Key)
	assert.Equal(t, "23", backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
	assert.Nil(t, lru.Get("someKey1"))
}

func TestLRU_SetNewElementToNotFullCache(t *testing.T) {
	lru := New(3)
	lru.Set("someKey1", 8)
	lru.Set("someKey2", 3)

	frontItem := lru.queue.Front().Value.(*Item)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Equal(t, "someKey2", frontItem.Key)
	assert.Equal(t, 3, frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 2, lru.queue.Len())
}

func TestLRU_SetNewElementAsync(t *testing.T) {
	wg := sync.WaitGroup{}
	lru := New(3)
	wg.Add(3)

	go func() {
		lru.Set("someKey1", 8)
		wg.Done()
	}()
	go func() {
		lru.Set("someKey2", "AAA")
		wg.Done()
	}()
	go func() {
		lru.Set("someKey3", Item{"key", 7})
		wg.Done()
	}()

	wg.Wait()

	assert.Equal(t, 8, lru.Get("someKey1"))
	assert.Equal(t, "AAA", lru.Get("someKey2"))
	assert.Equal(t, Item{"key", 7}, lru.Get("someKey3"))
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_GetHasElement(t *testing.T) {
	lru := New(3)
	lru.Set("someKey1", 8)
	lru.Set("someKey2", 3)
	lru.Set("someKey3", 0)

	item := lru.Get("someKey2")

	frontItem := lru.queue.Front().Value.(*Item)
	backItem := lru.queue.Back().Value.(*Item)

	assert.Equal(t, 3, item)
	assert.Equal(t, "someKey2", frontItem.Key)
	assert.Equal(t, 3, frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_GetHasNotElement(t *testing.T) {
	lru := New(3)
	lru.Set("someKey1", 8)
	lru.Set("someKey2", 3)
	lru.Set("someKey3", 0)

	item := lru.Get("someKey")

	frontItem := lru.queue.Front().Value.(*Item)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Nil(t, item)
	assert.Equal(t, "someKey3", frontItem.Key)
	assert.Equal(t, 0, frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}

func TestLRU_RemoveHasElement(t *testing.T) {
	lru := New(3)
	lru.Set("someKey1", 8)
	lru.Set("someKey2", 3)
	lru.Set("someKey3", 0)

	result := lru.Remove("someKey2")

	frontItem := lru.queue.Front().Value.(*Item)
	backItem := lru.queue.Back().Value.(*Item)
	assert.Nil(t, lru.Get("someKey2"))
	assert.True(t, result)
	assert.Equal(t, "someKey3", frontItem.Key)
	assert.Equal(t, 0, frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 2, lru.queue.Len())
}

func TestLRU_RemoveHasNotElement(t *testing.T) {
	lru := New(3)
	lru.Set("someKey1", 8)
	lru.Set("someKey2", 3)
	lru.Set("someKey3", 0)

	result := lru.Remove("someKey")

	frontItem := lru.queue.Front().Value.(*Item)
	backItem := lru.queue.Back().Value.(*Item)
	assert.True(t, result)
	assert.Equal(t, "someKey3", frontItem.Key)
	assert.Equal(t, 0, frontItem.Value)
	assert.Equal(t, "someKey1", backItem.Key)
	assert.Equal(t, 8, backItem.Value)
	assert.Equal(t, 3, lru.queue.Len())
}
