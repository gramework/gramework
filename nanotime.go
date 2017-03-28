package gramework

import _ "unsafe" // required to use //go:linkname

// Nanotime is monotonic time provider.
func Nanotime() int64 {
	return nanotime()
}

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64
