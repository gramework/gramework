package runtimer

import "unsafe" // #nosec

func PtrToType(ptr unsafe.Pointer) *Type {
	return (*Type)(ptr)
}
