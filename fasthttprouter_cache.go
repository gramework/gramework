package gramework

import (
	"runtime"
	"sync"
	"time"
)

const cacheRecordTTLDelta = 20 * 1000000000

type cache struct {
	v  map[string]*cacheRecord
	mu sync.RWMutex
}

type cacheRecord struct {
	n              *node
	tsr            bool
	values         map[string]string
	lastAccessTime int64
}

func (c *cache) Put(path string, n *node, tsr bool) {
	c.mu.Lock()
	c.v[path] = &cacheRecord{
		n:              n,
		tsr:            tsr,
		lastAccessTime: Nanotime(),
	}
	c.mu.Unlock()
}

func (c *cache) PutWild(path string, n *node, tsr bool, values map[string]string) {
	c.mu.Lock()
	c.v[path] = &cacheRecord{
		n:              n,
		tsr:            tsr,
		values:         values,
		lastAccessTime: Nanotime(),
	}
	c.mu.Unlock()
}

func (c *cache) Get(path string) (n *cacheRecord, ok bool) {
	c.mu.RLock()
	n, ok = c.v[path]
	if ok {
		n.lastAccessTime = Nanotime()
	}
	c.mu.RUnlock()
	return
}

func (c *cache) maintain() {
	for {
		runtime.Gosched()
		time.Sleep(30 * time.Second)
		if len(c.v) < 256 {
			continue
		}
		c.mu.Lock()
		for path := range c.v {
			if Nanotime()-cacheRecordTTLDelta > c.v[path].lastAccessTime {
				c.v[path].n.hits = 0
				delete(c.v, path)
			}
		}
		c.mu.Unlock()
	}
}
