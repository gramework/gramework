// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtimer

import "unsafe" // #nosec

const (
	C0 = uintptr((8-PtrSize)/4*2860486313 + (PtrSize-4)/4*33054211828000289)
	C1 = uintptr((8-PtrSize)/4*3267000013 + (PtrSize-4)/4*23344194077549503)
)

// type algorithms - known to compiler
const (
	AlgNOEQ = iota
	AlgMEM0
	AlgMEM8
	AlgMEM16
	AlgMEM32
	AlgMEM64
	AlgMEM128
	AlgSTRING
	AlgINTER
	AlgNILINTER
	AlgFLOAT32
	AlgFLOAT64
	AlgCPLX64
	AlgCPLX128
	AlgMax
)

// TypeAlg is also copied/used in reflect/type.go.
// keep them in sync.
type TypeAlg struct {
	// function for hashing objects of this type
	// (ptr to object, seed) -> hash
	Hash func(unsafe.Pointer, uintptr) uintptr
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
}

func Memhash0(p unsafe.Pointer, h uintptr) uintptr {
	return memhash0(p, h)
}

//go:linkname memhash0 runtime.memhash0
func memhash0(p unsafe.Pointer, h uintptr) uintptr

func Memhash8(p unsafe.Pointer, h uintptr) uintptr {
	return memhash8(p, h)
}

func Memhash16(p unsafe.Pointer, h uintptr) uintptr {
	return memhash16(p, h)
}

func Memhash32(p unsafe.Pointer, h uintptr) uintptr {
	return memhash32(p, h)
}

func Memhash64(p unsafe.Pointer, h uintptr) uintptr {
	return memhash64(p, h)
}

//go:linkname memhash8 runtime.memhash8
func memhash8(p unsafe.Pointer, h uintptr) uintptr

//go:linkname memhash16 runtime.memhash16
func memhash16(p unsafe.Pointer, h uintptr) uintptr

//go:linkname memhash32 runtime.memhash32
func memhash32(p unsafe.Pointer, h uintptr) uintptr

//go:linkname memhash64 runtime.memhash64
func memhash64(p unsafe.Pointer, h uintptr) uintptr

func Memhash128(p unsafe.Pointer, h uintptr) uintptr {
	return memhash128(p, h)
}

//go:linkname memhash128 runtime.memhash128
func memhash128(p unsafe.Pointer, h uintptr) uintptr

// MemhashVarlen is defined in runtime assembly because it needs access
// to the closure. It appears here to provide an argument
// signature for the assembly routine.
func MemhashVarlen(p unsafe.Pointer, h uintptr) uintptr {
	return memhashVarlen(p, h)
}

//go:linkname memhash128 runtime.memhash_varlen
func memhashVarlen(p unsafe.Pointer, h uintptr) uintptr

var AlgArray = [AlgMax]TypeAlg{
	AlgNOEQ:     {nil, nil},
	AlgMEM0:     {Memhash0, Memequal0},
	AlgMEM8:     {Memhash8, Memequal8},
	AlgMEM16:    {Memhash16, Memequal16},
	AlgMEM32:    {Memhash32, Memequal32},
	AlgMEM64:    {Memhash64, Memequal64},
	AlgMEM128:   {Memhash128, Memequal128},
	AlgSTRING:   {Strhash, Strequal},
	AlgINTER:    {Interhash, Interequal},
	AlgNILINTER: {Nilinterhash, Nilinterequal},
	AlgFLOAT32:  {F32hash, F32equal},
	AlgFLOAT64:  {F64hash, F64equal},
	AlgCPLX64:   {C64hash, C64equal},
	AlgCPLX128:  {C128hash, C128equal},
}

func Aeshash(p unsafe.Pointer, h, s uintptr) uintptr {
	return aeshash(p, h, s)
}

//go:linkname aeshash runtime.aeshash
func aeshash(p unsafe.Pointer, h, s uintptr) uintptr

func Aeshash32(p unsafe.Pointer, h uintptr) uintptr {
	return aeshash32(p, h)
}

//go:linkname aeshash32 runtime.aeshash32
func aeshash32(p unsafe.Pointer, h uintptr) uintptr

func Aeshash64(p unsafe.Pointer, h uintptr) uintptr {
	return aeshash64(p, h)
}

//go:linkname aeshash64 runtime.aeshash64
func aeshash64(p unsafe.Pointer, h uintptr) uintptr

func Aeshashstr(p unsafe.Pointer, h uintptr) uintptr {
	return aeshashstr(p, h)
}

//go:linkname aeshashstr runtime.aeshashstr
func aeshashstr(p unsafe.Pointer, h uintptr) uintptr

func Strhash(p unsafe.Pointer, h uintptr) uintptr {
	return strhash(p, h)
}

//go:linkname strhash runtime.strhash
func strhash(a unsafe.Pointer, h uintptr) uintptr

func F32hash(p unsafe.Pointer, h uintptr) uintptr {
	return f32hash(p, h)
}

//go:linkname f32hash runtime.f32hash
func f32hash(p unsafe.Pointer, h uintptr) uintptr

func F64hash(p unsafe.Pointer, h uintptr) uintptr {
	return f64hash(p, h)
}

//go:linkname f64hash runtime.f64hash
func f64hash(p unsafe.Pointer, h uintptr) uintptr

func C64hash(p unsafe.Pointer, h uintptr) uintptr {
	return c64hash(p, h)
}

//go:linkname c64hash runtime.c64hash
func c64hash(p unsafe.Pointer, h uintptr) uintptr

func C128hash(p unsafe.Pointer, h uintptr) uintptr {
	return c128hash(p, h)
}

//go:linkname c128hash runtime.c128hash
func c128hash(p unsafe.Pointer, h uintptr) uintptr

func Interhash(p unsafe.Pointer, h uintptr) uintptr {
	return interhash(p, h)
}

//go:linkname interhash runtime.interhash
func interhash(p unsafe.Pointer, h uintptr) uintptr

func Nilinterhash(p unsafe.Pointer, h uintptr) uintptr {
	return nilinterhash(p, h)
}

//go:linkname nilinterhash runtime.nilinterhash
func nilinterhash(p unsafe.Pointer, h uintptr) uintptr

func Memequal(a, b unsafe.Pointer, size uintptr) bool {
	return memequal(a, b, size)
}

//go:linkname memequal runtime.memequal
func memequal(a, b unsafe.Pointer, size uintptr) bool

func Memequal0(p, q unsafe.Pointer) bool {
	return memequal0(p, q)
}

//go:linkname memequal0 runtime.memequal0
func memequal0(p, q unsafe.Pointer) bool

func Memequal8(p, q unsafe.Pointer) bool {
	return memequal8(p, q)
}

//go:linkname memequal8 runtime.memequal8
func memequal8(p, q unsafe.Pointer) bool

func Memequal16(p, q unsafe.Pointer) bool {
	return memequal16(p, q)
}

//go:linkname memequal18 runtime.memequal18
func memequal16(p, q unsafe.Pointer) bool

func Memequal32(p, q unsafe.Pointer) bool {
	return memequal32(p, q)
}

//go:linkname memequal32 runtime.memequal32
func memequal32(p, q unsafe.Pointer) bool

func Memequal64(p, q unsafe.Pointer) bool {
	return memequal64(p, q)
}

//go:linkname memequal64 runtime.memequal64
func memequal64(p, q unsafe.Pointer) bool

func Memequal128(p, q unsafe.Pointer) bool {
	return memequal128(p, q)
}

//go:linkname memequal128 runtime.memequal128
func memequal128(p, q unsafe.Pointer) bool

func F32equal(p, q unsafe.Pointer) bool {
	return f32equal(p, q)
}

//go:linkname f32equal runtime.f32equal
func f32equal(p, q unsafe.Pointer) bool

func F64equal(p, q unsafe.Pointer) bool {
	return f64equal(p, q)
}

//go:linkname f64equal runtime.f64equal
func f64equal(p, q unsafe.Pointer) bool

func C64equal(p, q unsafe.Pointer) bool {
	return c64equal(p, q)
}

//go:linkname c64equal runtime.c64equal
func c64equal(p, q unsafe.Pointer) bool

func C128equal(p, q unsafe.Pointer) bool {
	return c128equal(p, q)
}

//go:linkname c128equal runtime.c128equal
func c128equal(p, q unsafe.Pointer) bool

func Strequal(p, q unsafe.Pointer) bool {
	return strequal(p, q)
}

//go:linkname strequal runtime.strequal
func strequal(p, q unsafe.Pointer) bool

func Interequal(p, q unsafe.Pointer) bool {
	return interequal(p, q)
}

//go:linkname interequal runtime.interequal
func interequal(p, q unsafe.Pointer) bool

func Nilinterequal(p, q unsafe.Pointer) bool {
	return nilinterequal(p, q)
}

//go:linkname nilinterequal runtime.nilinterequal
func nilinterequal(p, q unsafe.Pointer) bool

func Efaceeq(t *Type, x, y unsafe.Pointer) bool {
	return efaceeq(t, x, y)
}

//go:linkname efaceeq runtime.efaceeq
func efaceeq(t *Type, x, y unsafe.Pointer) bool

func Ifaceeq(t *Itab, x, y unsafe.Pointer) bool {
	return ifaceeq(t, x, y)
}

//go:linkname ifaceeq runtime.ifaceeq
func ifaceeq(tab *Itab, x, y unsafe.Pointer) bool

func StringHash(s string, seed uintptr) uintptr {
	return stringHash(s, seed)
}

// Testing adapters for hash quality tests (see hash_test.go)
//go:linkname stringHash runtime.stringHash
func stringHash(s string, seed uintptr) uintptr

func BytesHash(b []byte, seed uintptr) uintptr {
	return bytesHash(b, seed)
}

//go:linkname bytesHash runtime.bytesHash
func bytesHash(b []byte, seed uintptr) uintptr

func Int32Hash(i uint32, seed uintptr) uintptr {
	return int32Hash(i, seed)
}

//go:linkname int32Hash runtime.int32Hash
func int32Hash(i uint32, seed uintptr) uintptr

func Int64Hash(i uint64, seed uintptr) uintptr {
	return int64Hash(i, seed)
}

//go:linkname int64Hash runtime.int64Hash
func int64Hash(i uint64, seed uintptr) uintptr

func EfaceHash(i interface{}, seed uintptr) uintptr {
	return efaceHash(i, seed)
}

//go:linkname efaceHash runtime.efaceHash
func efaceHash(i interface{}, seed uintptr) uintptr

func IfaceHash(i interface {
	F()
}, seed uintptr) uintptr {
	return ifaceHash(i, seed)
}

//go:linkname ifaceHash runtime.ifaceHash
func ifaceHash(i interface {
	F()
}, seed uintptr) uintptr

const HashRandomBytes = PtrSize / 4 * 64

//go:linkname CPUIDECX runtime.cpuid_ecx
var CPUIDECX uint32
