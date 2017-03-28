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
	n      *node
	tsr    bool
	values map[string]string
}

func (c *cache) Put(path string, n *node, tsr bool) {
	c.mu.Lock()
	c.v[path] = &cacheRecord{
		n:   n,
		tsr: tsr,
	}
	c.mu.Unlock()
}

func (c *cache) PutWild(path string, n *node, tsr bool, values map[string]string) {
	c.mu.Lock()
	c.v[path] = &cacheRecord{
		n:      n,
		tsr:    tsr,
		values: values,
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
		c.mu.Lock()
		for path := range c.v {
			c.v[path].n.hits = 0
		}
		c.v = make(map[string]*cacheRecord, 0)
		c.mu.Unlock()
	}
}
