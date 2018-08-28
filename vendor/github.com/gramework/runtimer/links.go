// Package runtimer provides you unsafe way to use runtime internals
package runtimer

import (
	"unsafe" // #nosec required to use go:linkname
)

func Or8(ptr *uint8, val uint8) {
	or8(ptr, val)
}

//go:noescape
//go:linkname or8 runtime.internal.atomic.Or8
func or8(ptr *uint8, val uint8)

//go:linkname PtrSize runtime.internal.sys.PtrSize
const PtrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

func Cas(ptr *uint32, old, new uint32) bool {
	return cas(ptr, old, new)
}

//go:noescape
//go:linkname cas runtime.internal.atomic.Cas
func cas(ptr *uint32, old, new uint32) bool

func Casp1(ptr *unsafe.Pointer, old, new unsafe.Pointer) bool {
	return casp1(ptr, old, new)
}

// NO go:noescape annotation; see atomic_pointer.go.
//go:linkname casp1 runtime.internal.atomic.Casp1
func casp1(ptr *unsafe.Pointer, old, new unsafe.Pointer) bool

func Casuintptr(ptr *uintptr, old, new uintptr) bool {
	return casuintptr(ptr, old, new)
}

//go:noescape
//go:linkname casuintptr runtime.internal.atomic.Casuintptr
func casuintptr(ptr *uintptr, old, new uintptr) bool

func Storeuintptr(ptr *uintptr, new uintptr) {
	storeuintptr(ptr, new)
}

//go:noescape
//go:linkname storeuintptr runtime.internal.atomic.Storeuintptr
func storeuintptr(ptr *uintptr, new uintptr)

func Loaduintptr(ptr *uintptr) uintptr {
	return loaduintptr(ptr)
}

//go:noescape
//go:linkname loaduintptr runtime.internal.atomic.Loaduintptr
func loaduintptr(ptr *uintptr) uintptr

func Loaduint(ptr *uint) uint {
	return loaduint(ptr)
}

//go:noescape
//go:linkname loaduint runtime.internal.atomic.Loaduint
func loaduint(ptr *uint) uint

func Loadint64(ptr *int64) int64 {
	return loadint64(ptr)
}

//go:noescape
//go:linkname loadint64 runtime.internal.atomic.Loadint64
func loadint64(ptr *int64) int64

func Fastrand() uint32 {
	return fastrand()
}

//go:linkname fastrand runtime.fastrand
func fastrand() uint32

func Throw(s string) {
	throw(s)
}

//go:linkname throw runtime.throw
func throw(s string)

func Newarray(typ *Type, n int) unsafe.Pointer {
	return newarray(typ, n)
}

//go:linkname newarray runtime.newarray
func newarray(typ *Type, n int) unsafe.Pointer

func Newobject(typ *Type) unsafe.Pointer {
	return newobject(typ)
}

//go:linkname newobject runtime.newobject
func newobject(typ *Type) unsafe.Pointer

func Typedmemmove(typ *Type, dst, src unsafe.Pointer) {
	typedmemmove(typ, dst, src)
}

//go:linkname typedmemmove runtime.typedmemmove
func typedmemmove(typ *Type, dst, src unsafe.Pointer)

func Typedmemclr(typ *Type, ptr unsafe.Pointer) {
	typedmemclr(typ, ptr)
}

//go:linkname typedmemclr runtime.typedmemclr
func typedmemclr(typ *Type, ptr unsafe.Pointer)

func Lock(l *Mutex) {
	lock(l)
}

//go:linkname lock runtime.lock
func lock(l *Mutex)

func Unlock(l *Mutex) {
	unlock(l)
}

//go:linkname unlock runtime.unlock
func unlock(l *Mutex)

func Msanread(addr unsafe.Pointer, sz uintptr) {
	msanread(addr, sz)
}

//go:linkname msanread runtime.msanread
func msanread(addr unsafe.Pointer, sz uintptr)

func MemclrHasPointers(ptr unsafe.Pointer, n uintptr) {
	memclrHasPointers(ptr, n)
}

//go:linkname memclrHasPointers runtime.memclrHasPointers
func memclrHasPointers(ptr unsafe.Pointer, n uintptr)

//go:linkname Hex runtime.hex
type Hex uint64

func MemclrNoHeapPointers(ptr unsafe.Pointer, n uintptr) {
	memclrNoHeapPointers(ptr, n)
}

//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

func Noescape(p unsafe.Pointer) unsafe.Pointer {
	return noescape(p)
}

//go:linkname noescape runtime.noescape
func noescape(p unsafe.Pointer) unsafe.Pointer

func Xaddint64(ptr *int64, delta int64) int64 {
	return xaddint64(ptr, delta)
}

//go:noescape
//go:linkname xaddint64 runtime.internal.atomic.Xaddint64
func xaddint64(ptr *int64, delta int64) int64

func Encoderune(p []byte, r rune) int {
	return encoderune(p, r)
}

//go:linkname encoderune runtime.encoderune
func encoderune(p []byte, r rune) int

func Rawstring(size int) (s string, b []byte) {
	return rawstring(size)
}

//go:linkname rawstring runtime.rawstring
func rawstring(size int) (s string, b []byte)

const (
	//go:linkname GOOS runtime.internal.sys.GOOS
	GOOS = `unknown`

	//go:linkname GoosAndroid runtime.internal.sys.GoosAndroid
	GoosAndroid uint = 0
	//go:linkname GoosDarwin runtime.internal.sys.GoosDarwin
	GoosDarwin uint = 0
	//go:linkname GoosDragonfly runtime.internal.sys.GoosDragonfly
	GoosDragonfly uint = 0
	//go:linkname GoosFreebsd runtime.internal.sys.GoosFreebsd
	GoosFreebsd uint = 0
	//go:linkname GoosLinux runtime.internal.sys.GoosLinux
	GoosLinux uint = 0
	//go:linkname GoosNacl runtime.internal.sys.GoosNacl
	GoosNacl uint = 0
	//go:linkname GoosNetbsd runtime.internal.sys.GoosNetbsd
	GoosNetbsd uint = 0
	//go:linkname GoosOpenbsd runtime.internal.sys.GoosOpenbsd
	GoosOpenbsd uint = 0
	//go:linkname GoosPlan9 runtime.internal.sys.GoosPlan9
	GoosPlan9 uint = 0
	//go:linkname GoosSolaris runtime.internal.sys.GoosSolaris
	GoosSolaris uint = 0
	//go:linkname GoosWindows runtime.internal.sys.GoosWindows
	GoosWindows uint = 0

	//go:linkname Goarch386 runtime.internal.sys.Goarch386
	Goarch386 uint = 0
	//go:linkname GoarchAmd64 runtime.internal.sys.GoarchAmd64
	GoarchAmd64 uint = 0
	//go:linkname GoarchAmd64p32 runtime.internal.sys.GoarchAmd64p32
	GoarchAmd64p32 uint = 0
	//go:linkname GoarchArm runtime.internal.sys.GoarchArm
	GoarchArm uint = 0
	//go:linkname GoarchArmbe runtime.internal.sys.GoarchArmbe
	GoarchArmbe uint = 0
	//go:linkname GoarchArm64 runtime.internal.sys.GoarchArm64
	GoarchArm64 uint = 0
	//go:linkname GoarchArm64be runtime.internal.sys.GoarchArm64be
	GoarchArm64be uint = 0
	//go:linkname GoarchPpc64 runtime.internal.sys.GoarchPpc64
	GoarchPpc64 uint = 0
	//go:linkname GoarchPpc64le runtime.internal.sys.GoarchPpc64le
	GoarchPpc64le uint = 0
	//go:linkname GoarchMips runtime.internal.sys.GoarchMips
	GoarchMips uint = 0
	//go:linkname GoarchMipsle runtime.internal.sys.GoarchMipsle
	GoarchMipsle uint = 0
	//go:linkname GoarchMips64 runtime.internal.sys.GoarchMips64
	GoarchMips64 uint = 0
	//go:linkname GoarchMips64le runtime.internal.sys.GoarchMips64le
	GoarchMips64le uint = 0
	//go:linkname GoarchMips64p32 runtime.internal.sys.GoarchMips64p32
	GoarchMips64p32 uint = 0
	//go:linkname GoarchMips64p32le runtime.internal.sys.GoarchMips64p32le
	GoarchMips64p32le uint = 0
	//go:linkname GoarchPpc runtime.internal.sys.GoarchPpc
	GoarchPpc uint = 0
	//go:linkname GoarchS390 runtime.internal.sys.GoarchS390
	GoarchS390 uint = 0
	//go:linkname GoarchS390x runtime.internal.sys.GoarchS390x
	GoarchS390x uint = 0
	//go:linkname GoarchSparc runtime.internal.sys.GoarchSparc
	GoarchSparc uint = 0
	//go:linkname GoarchSparc64 runtime.internal.sys.GoarchSparc64
	GoarchSparc64 uint = 0
)
