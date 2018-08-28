// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package client

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
