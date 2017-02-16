package apiClient

import (
	"sync/atomic"
)

func newRangeBalancer() *rangeBalancer {
	var curr, total int64
	return &rangeBalancer{
		curr:  &curr,
		total: &total,
	}
}

func (i *rangeBalancer) next() int64 {
	if atomic.LoadInt64(i.curr) == atomic.LoadInt64(i.total) {
		return atomic.SwapInt64(i.curr, 0)
	}
	return atomic.AddInt64(i.curr, 1)
}
