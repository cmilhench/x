package cache

import (
	"runtime"
	"sync"
	"time"
)

type Item struct {
	Object     interface{}
	Expiration int64
}

type Cache struct {
	items map[string]Item
	lock  sync.RWMutex
	stop  chan struct{}
}

func New(interval time.Duration) *Cache {
	c := &Cache{
		items: make(map[string]Item),
		stop:  make(chan struct{}),
	}
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				clean(c)
			case <-c.stop:
				ticker.Stop()
				return
			}
		}
	}()
	runtime.SetFinalizer(c, func(c *Cache) {
		c.stop <- struct{}{}
	})
	return c
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}
	c.lock.Lock()
	defer func() {
		c.lock.Unlock()
	}()
	c.items[key] = Item{
		Object:     value,
		Expiration: exp,
	}
}

func (c *Cache) Get(key string) (value interface{}, found bool) {
	c.lock.RLock()
	defer func() {
		c.lock.RUnlock()
	}()
	item, found := c.items[key]

	if !found {
		//log.Debugf("  - %s not found in cache of %d items %p", key, len(c.items), c)
		return
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			found = false
			//log.Debugf("  - %s not found in cache of %d items %p", key, len(c.items), c)
			return
		}
	}
	value = item.Object
	//log.Debugf("  - %s found in cache of %d items %p", key, len(c.items), c)
	return
}

func (c *Cache) Count() int {
	c.lock.RLock()
	defer func() {
		c.lock.RUnlock()
	}()
	return len(c.items)
}

func (c *Cache) Delete(key string) {
	c.lock.Lock()
	defer func() {
		c.lock.Unlock()
	}()
	delete(c.items, key)
	//log.Debugf("  - %s removed from cache of %d items %p", key, len(c.items), c)
}

func (c *Cache) Flush() {
	c.lock.Lock()
	defer func() {
		c.lock.Unlock()
	}()
	c.items = map[string]Item{}
	//log.Debugf("  - Everything removed from cache of %d items %p", len(c.items), c)
}

// -- Housekeeping

func clean(c *Cache) {
	c.lock.Lock()
	defer func() {
		c.lock.Unlock()
	}()
	now := time.Now().UnixNano()
	for key, item := range c.items {
		if item.Expiration > 0 && now > item.Expiration {
			delete(c.items, key)
		}
	}
}
