package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	if _, ok := cache.items[key]; ok {
		newItem := cache.items[key]
		newItem.Value = cacheItem{key: key, value: value}
		cache.queue.MoveToFront(newItem)
		return true
	}

	newItem := cache.queue.PushFront(cacheItem{key: key, value: value})

	if cache.queue.Len() > cache.capacity {
		tail := cache.queue.Back()
		cache.queue.Remove(tail)

		item, ok := tail.Value.(cacheItem)

		if ok {
			delete(cache.items, item.key)
		}
	}

	cache.items[key] = newItem
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	if cache.items[key] == nil {
		return nil, false
	}

	if item, ok := cache.items[key].Value.(cacheItem); ok {
		return item.value, ok
	}

	return nil, false
}

func (cache *lruCache) Clear() {
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.queue = NewList()
}
