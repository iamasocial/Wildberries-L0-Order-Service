package cache

import (
	"L0/internal/entities"
	"container/list"
	"sync"
)

type Cache interface {
	Get(orderUID string) (*entities.Order, bool)
	Set(order *entities.Order)
	GetSize() int
}

type entry struct {
	key   string
	value *entities.Order
}

type lruCache struct {
	mu           sync.RWMutex
	capacity     int
	orders       map[string]*list.Element
	evictionList *list.List
}

func NewLRUCache(capacity int) Cache {
	return &lruCache{
		capacity:     capacity,
		orders:       make(map[string]*list.Element),
		evictionList: list.New(),
	}
}

func (c *lruCache) Get(orderUID string) (*entities.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, ok := c.orders[orderUID]; ok {
		c.evictionList.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

func (c *lruCache) Set(order *entities.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.orders[order.OrderUID]; ok {
		c.evictionList.MoveToFront(elem)
		elem.Value.(*entry).value = order
		return
	}

	elem := c.evictionList.PushFront(&entry{key: order.OrderUID, value: order})
	c.orders[order.OrderUID] = elem

	if c.evictionList.Len() > c.capacity {
		backElem := c.evictionList.Back()
		c.evictionList.Remove(backElem)
		delete(c.orders, backElem.Value.(*entry).key)
	}
}

func (c *lruCache) GetSize() int {
	return c.capacity
}
