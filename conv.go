package gramework

import (
	"unsafe"

	"github.com/gramework/runtimer"
)

// BytesToString effectively converts bytes to string
// nolint: gas
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes effectively converts string to bytes
// nolint: gas
func StringToBytes(s string) []byte {
	strstruct := runtimer.StringStructOf(&s)
	return *(*[]byte)(unsafe.Pointer(&runtimer.SliceType2{
		Array: strstruct.Str,
		Len:   strstruct.Len,
		Cap:   strstruct.Len,
	}))
}
