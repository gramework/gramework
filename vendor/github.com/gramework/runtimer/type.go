// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Runtime Type representation.

package runtimer

import (
	"unsafe"
) // #nosec

// Tflag is documented in reflect/Type.go.
//
// Tflag values must be kept in sync with copies in:
//	cmd/compile/internal/gc/reflect.go
//	cmd/link/internal/ld/decodesym.go
//	reflect/Type.go
type Tflag uint8

const (
	TflagUncommon  Tflag = 1 << 0
	TflagExtraStar Tflag = 1 << 1
	TflagNamed     Tflag = 1 << 2
)

// Needs to be in sync with ../cmd/compile/internal/ld/decodesym.go:/^func.commonsize,
// ../cmd/compile/internal/gc/reflect.go:/^func.dcommonType and
// ../reflect/Type.go:/^Type.rType.
//go:linkname Type runtime._type
type Type struct {
	Size       uintptr
	Ptrdata    uintptr // size of memory prefix holding all pointers
	Hash       uint32
	Tflag      Tflag
	Align      uint8
	Fieldalign uint8
	Kind       uint8
	Alg        *TypeAlg
	// gcdata stores the GC Type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	Gcdata    *byte
	Str       NameOff
	PtrToThis TypeOff
}

func (t *Type) String() string {
	s := t.NameOff(t.Str).Name()
	if t.Tflag&TflagExtraStar != 0 {
		return s[1:]
	}
	return s
}

func (t *Type) Uncommon() *UncommonType {
	if t.Tflag&TflagUncommon == 0 {
		return nil
	}
	switch t.Kind & KindMask {
	case KindStruct:
		type u struct {
			StructType
			u UncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u
	case KindPtr:
		type u struct {
			PtrType
			u UncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u
	case KindFunc:
		type u struct {
			FuncType
			u UncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u
	case KindSlice:
		type u struct {
			SliceType
			u UncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u
	case KindArray:
		type u struct {
			ArrayType
			u UncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u
	case KindChan:
		type u struct {
			ChanType
			u UncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u
	case KindMap:
		type u struct {
			MapType
			u UncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u
	case KindInterface:
		type u struct {
			InterfaceType
			u UncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u
	default:
		type u struct {
			Type
			u UncommonType
		}
		return &(*u)(unsafe.Pointer(t)).u
	}
}

func (t *Type) Name() string {
	if t.Tflag&TflagNamed == 0 {
		return ""
	}
	s := t.String()
	i := len(s) - 1
	for i >= 0 {
		if s[i] == '.' {
			break
		}
		i--
	}
	return s[i+1:]
}

// Mutex - Mutual exclusion locks.  In the uncontended case,
// as fast as spin locks (just a few user-level instructions),
// but on the contention path they sleep in the kernel.
// A zeroed Mutex is unlocked (no need to initialize each lock).
type Mutex struct {
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	Key uintptr
}

func ResolveNameOff(ptrInModule unsafe.Pointer, off NameOff) Name {
	return resolveNameOff(ptrInModule, off)
}

//go:linkname resolveNameOff runtime.resolveNameOff
func resolveNameOff(ptrInModule unsafe.Pointer, off NameOff) Name

func (t *Type) NameOff(off NameOff) Name {
	return resolveNameOff(unsafe.Pointer(t), off)
}

func ResolveTypeOff(ptrInModule unsafe.Pointer, off TypeOff) *Type {
	return resolveTypeOff(ptrInModule, off)
}

//go:linkname resolveTypeOff runtime.resolveTypeOff
func resolveTypeOff(ptrInModule unsafe.Pointer, off TypeOff) *Type

func (t *Type) TypeOff(off TypeOff) *Type {
	return ResolveTypeOff(unsafe.Pointer(t), off)
}

// func (t *Type) TextOff(off textOff) unsafe.Pointer {
// 	return t.textOff(off)
// }

// //go:linkname Type.textOff runtime._type.textOff
// func (t *Type) textOff(off textOff) unsafe.Pointer

func (t *FuncType) in() []*Type {
	// See FuncType in reflect/Type.go for details on data layout.
	uadd := uintptr(unsafe.Sizeof(FuncType{}))
	if t.typ.Tflag&TflagUncommon != 0 {
		uadd += unsafe.Sizeof(UncommonType{})
	}
	return (*[1 << 20]*Type)(Add(unsafe.Pointer(t), uadd))[:t.inCount]
}

func (t *FuncType) out() []*Type {
	// See FuncType in reflect/Type.go for details on data layout.
	uadd := uintptr(unsafe.Sizeof(FuncType{}))
	if t.typ.Tflag&TflagUncommon != 0 {
		uadd += unsafe.Sizeof(UncommonType{})
	}
	outCount := t.outCount & (1<<15 - 1)
	return (*[1 << 20]*Type)(Add(unsafe.Pointer(t), uadd))[t.inCount : t.inCount+outCount]
}

func (t *FuncType) Dotdotdot() bool {
	return t.outCount&(1<<15) != 0
}

type NameOff int32
type TypeOff int32
type textOff int32

type method struct {
	Name NameOff
	mtyp TypeOff
	ifn  textOff
	tfn  textOff
}

type UncommonType struct {
	Pkgpath NameOff
	Mcount  uint16 // number of methods
	_       uint16 // unused
	Moff    uint32 // offset from this UncommonType to [mcount]method
	_       uint32 // unused
}

type imethod struct {
	Name NameOff
	ityp TypeOff
}

type InterfaceType struct {
	typ     Type
	pkgpath Name
	mhdr    []imethod
}

type MapType struct {
	Typ           Type
	Key           *Type
	Elem          *Type
	Bucket        *Type  // internal Type representing a hash bucket
	Hmap          *Type  // internal Type representing a hmap
	Keysize       uint8  // size of key slot
	Indirectkey   bool   // store ptr to key instead of key itself
	Valuesize     uint8  // size of value slot
	Indirectvalue bool   // store ptr to value instead of value itself
	Bucketsize    uint16 // size of bucket
	Reflexivekey  bool   // true if k==k for all keys
	Needkeyupdate bool   // true if we need to update key on an overwrite
}

type ArrayType struct {
	Typ   Type
	Elem  *Type
	Slice *Type
	Len   uintptr
}

type ChanType struct {
	Typ  Type
	Elem *Type
	Dir  uintptr
}

type SliceType struct {
	Typ  Type
	Elem *Type
}

type SliceType2 struct {
	Array unsafe.Pointer
	Len   int
	Cap   int
}

type FuncType struct {
	typ      Type
	inCount  uint16
	outCount uint16
}

type PtrType struct {
	typ  Type
	elem *Type
}

type Structfield struct {
	Name       Name
	Typ        *Type
	OffsetAnon uintptr
}

func (f *Structfield) Offset() uintptr {
	return f.OffsetAnon >> 1
}

type StructType struct {
	Typ     Type
	PkgPath Name
	Fields  []Structfield
}

// Name is an encoded Type Name with optional extra data.
// See reflect/Type.go for details.
type Name struct {
	Bytes *byte
}

func (n Name) Data(off int) *byte {
	return (*byte)(Add(unsafe.Pointer(n.Bytes), uintptr(off)))
}

func (n Name) IsExported() bool {
	return (*n.Bytes)&(1<<0) != 0
}

func (n Name) NameLen() int {
	return int(uint16(*n.Data(1))<<8 | uint16(*n.Data(2)))
}

func (n Name) TagLen() int {
	if *n.Data(0)&(1<<1) == 0 {
		return 0
	}
	off := 3 + n.NameLen()
	return int(uint16(*n.Data(off))<<8 | uint16(*n.Data(off + 1)))
}

func (n Name) Name() (s string) {
	if n.Bytes == nil {
		return ""
	}
	nl := n.NameLen()
	if nl == 0 {
		return ""
	}
	hdr := (*StringStruct)(unsafe.Pointer(&s))
	hdr.Str = unsafe.Pointer(n.Data(3))
	hdr.Len = nl
	return s
}

func (n Name) Tag() (s string) {
	tl := n.TagLen()
	if tl == 0 {
		return ""
	}
	nl := n.NameLen()
	hdr := (*StringStruct)(unsafe.Pointer(&s))
	hdr.Str = unsafe.Pointer(n.Data(3 + nl + 2))
	hdr.Len = tl
	return s
}

func (n Name) PkgPath() string {
	if n.Bytes == nil || *n.Data(0)&(1<<2) == 0 {
		return ""
	}
	off := 3 + n.NameLen()
	if tl := n.TagLen(); tl > 0 {
		off += 2 + tl
	}
	var NameOff NameOff
	copy((*[4]byte)(unsafe.Pointer(&NameOff))[:], (*[4]byte)(unsafe.Pointer(n.Data(off)))[:])
	pkgPathName := resolveNameOff(unsafe.Pointer(n.Bytes), NameOff)
	return pkgPathName.Name()
}

// TypesEqual reports whether two Types are equal.
//
// Everywhere in the runtime and reflect packages, it is assumed that
// there is exactly one *Type per Go Type, so that pointer equality
// can be used to test if Types are equal. There is one place that
// breaks this assumption: buildmode=shared. In this case a Type can
// appear as two different pieces of memory. This is hidden from the
// runtime and reflect package by the per-module Typemap built in
// Typelinksinit. It uses TypesEqual to map Types from later modules
// back into earlier ones.
//
// Only Typelinksinit needs this function.
func TypesEqual(t, v *Type) bool {
	return TypesEqual(t, v)
}

//go:linkname typesEqual runtime.typesEqual
func typesEqual(t, v *Type) bool
