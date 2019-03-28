// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"runtime"
	"sync"
	"time"
)

type (
	cache struct {
		v  map[string]*msc
		mu sync.RWMutex
	}

	// method-specific cache
	msc struct {
		v  map[string]*cacheRecord
		mu sync.RWMutex
	}

	cacheRecord struct {
		n              *node
		tsr            bool
		values         map[string]string
		lastAccessTime int64
	}
)

const cacheRecordTTLDelta = 20 * 1000000000

func (c *cache) getOrInitMSC(method string) *msc {
	c.mu.Lock()
	if v, ok := c.v[method]; ok {
		c.mu.Unlock()
		return v
	}
	ms := &msc{
		v: make(map[string]*cacheRecord),
	}
	c.v[method] = ms
	c.mu.Unlock()
	return ms
}

func (c *cache) getMSC(method string) *msc {
	c.mu.RLock()
	if v, ok := c.v[method]; ok {
		c.mu.RUnlock()
		return v
	}
	c.mu.RUnlock()
	return nil
}

func (c *cache) Put(path string, n *node, tsr bool, method string) {
	msc := c.getOrInitMSC(method)
	msc.mu.Lock()
	msc.v[path] = &cacheRecord{
		n:              n,
		tsr:            tsr,
		lastAccessTime: Nanotime(),
	}
	msc.mu.Unlock()
}

func (c *cache) PutWild(path string, n *node, tsr bool, values map[string]string, method string) {
	msc := c.getOrInitMSC(method)
	msc.mu.Lock()
	msc.v[path] = &cacheRecord{
		n:              n,
		tsr:            tsr,
		values:         values,
		lastAccessTime: Nanotime(),
	}
	msc.mu.Unlock()
}

func (c *cache) Get(path string, method string) (n *cacheRecord, ok bool) {
	msc := c.getMSC(method)
	if msc == nil {
		return nil, false
	}
	msc.mu.RLock()
	n, ok = msc.v[path]
	if ok {
		n.lastAccessTime = Nanotime()
	}
	msc.mu.RUnlock()
	return
}

func (c *cache) maintain() {
	for {
		runtime.Gosched()
		time.Sleep(30 * time.Second)
		skipIter := true
		c.mu.RLock()
		for _, v := range c.v {
			v.mu.RLock()
			mscLen := len(v.v)
			v.mu.RUnlock()
			if mscLen > 256 {
				skipIter = false
				break
			}
		}
		if skipIter {
			c.mu.RUnlock()
			continue
		}
		for _, msc := range c.v {
			if len(msc.v) <= 256 {
				continue
			}
			msc.mu.Lock()
			for path := range msc.v {
				if Nanotime()-cacheRecordTTLDelta > msc.v[path].lastAccessTime {
					msc.v[path].n.hits = 0
					delete(c.v, path)
				}
			}
			msc.mu.Unlock()
		}
		c.mu.RUnlock()
	}
}
