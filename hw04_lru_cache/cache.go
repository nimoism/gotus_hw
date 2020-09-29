package hw04_lru_cache //nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	*sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.Lock()
	defer l.Unlock()

	if item, ok := l.items[key]; ok {
		item.Value.(*cacheItem).value = value
		l.queue.MoveToFront(item)
		return true
	}
	listItem := l.queue.PushFront(&cacheItem{key: key, value: value})
	l.items[key] = listItem
	if l.queue.Len() > l.capacity {
		itemToRemove := l.queue.Back()
		delete(l.items, itemToRemove.Value.(*cacheItem).key)
		l.queue.Remove(itemToRemove)
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.Lock()
	defer l.Unlock()

	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.Lock()
	defer l.Unlock()

	l.queue = NewList()
	l.items = map[Key]*ListItem{}
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		Mutex:    &sync.Mutex{},
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
