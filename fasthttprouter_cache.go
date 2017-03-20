package gramework

import (
	"runtime"
	"sync"
	"time"
)

type cache struct {
	v  map[string]*cacheRecord
	mu sync.Mutex
}

type cacheRecord struct {
	n   *node
	tsr bool
}

func (c *cache) Put(path string, n *node, tsr bool) {
	c.mu.Lock()
	c.v[path] = &cacheRecord{
		n:   n,
		tsr: tsr,
	}
	c.mu.Unlock()
}

func (c *cache) Get(path string) (n *cacheRecord, ok bool) {
	c.mu.Lock()
	n, ok = c.v[path]
	c.mu.Unlock()
	return
}

func (c *cache) maintain() {
	for {
		runtime.Gosched()
		time.Sleep(10 * time.Second)
		for path := range c.v {
			c.mu.Lock()
			c.v[path].n.hits = 0
			c.mu.Unlock()
		}
		c.v = make(map[string]*cacheRecord, 0)
	}
}
