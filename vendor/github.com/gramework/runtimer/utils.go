package runtimer

import (
	"unsafe" // #nosec
)

func PtrToString(ptr unsafe.Pointer) string {
	return *(*string)(ptr)
}

func PtrToStringPtr(ptr unsafe.Pointer) *string {
	return (*string)(ptr)
}

func PtrPtrToStringPtr(ptr *unsafe.Pointer) *string {
	return (*string)(*ptr)
}
