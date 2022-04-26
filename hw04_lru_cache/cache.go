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

func (l *lruCache) Set(key Key, value interface{}) bool {
	isKeyExist := false
	if elem, ok := l.items[key]; ok {
		elem.Value.(item).setValue(value)
		l.queue.MoveToFront(l.items[key])
		isKeyExist = true
	} else {
		elem := newCacheItem(key, value)
		if l.capacity == l.queue.Len() {
			lastElem := l.queue.Back()
			lastKey := lastElem.Value.(item).getKey()
			l.queue.Remove(lastElem)
			delete(l.items, lastKey)
			l.items[key] = l.queue.PushFront(elem)
		} else {
			l.items[key] = l.queue.PushFront(elem)
		}
	}
	return isKeyExist
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if elem, ok := l.items[key]; ok {
		l.queue.MoveToFront(elem)
		return elem.Value.(item).getValue(), true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

type item interface {
	getValue() interface{}
	getKey() Key
	setValue(v interface{})
}

func (c *cacheItem) getValue() interface{} {
	return c.value
}

func (c *cacheItem) getKey() Key {
	return c.key
}

func (c *cacheItem) setValue(v interface{}) {
	c.value = v
}

func newCacheItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{key: key, value: value}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
