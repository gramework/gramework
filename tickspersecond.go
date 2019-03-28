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
	_ "unsafe" // required to use //go:linkname
)

// TicksPerSecond reports cpu ticks per second counter
func TicksPerSecond() int64 {
	return tickspersecond()
}

//go:noescape
//go:linkname tickspersecond runtime.tickspersecond
func tickspersecond() int64
