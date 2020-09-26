package cache

import (
	"sync"
	"time"
)

//Cache caceh
type Cache struct {
	cache map[string]*entry
	time  int64
	sync.Mutex
}

type entry struct {
	res   result
	time  int64
	ready chan struct{}
}

type result struct {
	value []byte
	err   error
}

// VEFunc value error function
type VEFunc func() ([]byte, error)

// NewCache creates a new Cache instance timout in nanoseconds
func NewCache(timeout int64) *Cache {
	return &Cache{cache: make(map[string]*entry), time: timeout}
}

// Get gets the value of the key or loads it with the function
func (c *Cache) Get(key string, f VEFunc) ([]byte, error) {
	c.Lock()
	e := c.cache[key]
	if e == nil || e.time+c.time < time.Now().UnixNano() {
		e = &entry{ready: make(chan struct{})}
		c.cache[key] = e
		c.Unlock()

		defer func() {
			if r := recover(); r != nil {
				close(e.ready)
			}
		}()
		e.res.value, e.res.err = f()
		e.time = time.Now().UnixNano()
		close(e.ready)
	} else {
		c.Unlock()
		<-e.ready
	}
	return e.res.value, e.res.err
}
