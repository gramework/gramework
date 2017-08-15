package gramework

import (
	"unsafe"

	"github.com/gramework/runtimer"
)

type sh struct {
	d   unsafe.Pointer
	len int
}

// BytesToString effectively converts bytes to string
func BytesToString(b []byte) string {
	bh := (*sh)(unsafe.Pointer(&b))
	return *(*string)(unsafe.Pointer(&runtimer.StringStruct{Str: bh.d, Len: bh.len}))
}
