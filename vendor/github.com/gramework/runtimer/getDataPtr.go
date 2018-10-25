package runtimer

import "unsafe" // #nosec

func GetEfaceDataPtr(eface interface{}) unsafe.Pointer {
	return ((*[2]unsafe.Pointer)(unsafe.Pointer(&eface))[1])
}

func EfaceDataPtr(eface interface{}) *unsafe.Pointer {
	return &((*[2]unsafe.Pointer)(unsafe.Pointer(&eface))[1])
}
