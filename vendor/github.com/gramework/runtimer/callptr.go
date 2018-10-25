package runtimer

import "unsafe" // #nosec

// WIP, don't use those now
func callptr(ptr unsafe.Pointer) {}
func call(argtype *Type, fn, arg unsafe.Pointer, n uint32, retoffset uint32)
