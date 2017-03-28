package gramework

import _ "unsafe" // required to use //go:linkname

// Nanotime is monotonic time provider.
//
//go:noescape
//go:linkname Nanotime runtime.nanotime
func Nanotime() int64
