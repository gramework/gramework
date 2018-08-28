// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtimer

import "unsafe" // #nosec

type StringStruct struct {
	Str unsafe.Pointer
	Len int
}

// Variant with *byte pointer type for DWARF debugging.
type StringStructDWARF struct {
	Str *byte
	Len int
}

func StringStructOf(sp *string) *StringStruct {
	return (*StringStruct)(unsafe.Pointer(sp))
}

func Contains(s, t string) bool {
	return Index(s, t) >= 0
}

func Index(s, t string) int {
	if len(t) == 0 {
		return 0
	}
	for i := 0; i < len(s); i++ {
		if s[i] == t[0] && HasPrefix(s[i:], t) {
			return i
		}
	}
	return -1
}

func HasPrefix(s, t string) bool {
	return len(s) >= len(t) && s[:len(t)] == t
}

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
)

// Atoi parses an int from a string s.
// The bool result reports whether s is a number
// representable by a value of type int.
func Atoi(s string) (int, bool) {
	if s == "" {
		return 0, false
	}

	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	un := uint(0)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		if un > maxUint/10 {
			// overflow
			return 0, false
		}
		un *= 10
		un1 := un + uint(c) - '0'
		if un1 < un {
			// overflow
			return 0, false
		}
		un = un1
	}

	if !neg && un > uint(maxInt) {
		return 0, false
	}
	if neg && un > uint(maxInt)+1 {
		return 0, false
	}

	n := int(un)
	if neg {
		n = -n
	}

	return n, true
}

// Atoi32 is like Atoi but for integers
// that fit into an int32.
func Atoi32(s string) (int32, bool) {
	if n, ok := Atoi(s); n == int(int32(n)) {
		return int32(n), ok
	}
	return 0, false
}

func Findnull(s *byte) int {
	return findnull(s)
}

//go:nosplit
//go:linkname findnull runtime.findnull
func findnull(s *byte) int

func Findnullw(s *uint16) int {
	return findnullw(s)
}

//go:linkname findnullw runtime.findnullw
func findnullw(s *uint16) int

func Gostringnocopy(str *byte) string {
	return gostringnocopy(str)
}

//go:nosplit
//go:linkname gostringnocopy runtime.gostringnocopy
func gostringnocopy(str *byte) string

const (
	PageShift uint = 13

	// Public64bit = 1 on 64-bit systems, 0 on 32-bit systems
	Public64bit uint = 1 << (^uintptr(0) >> 63) / 2

	MHeapMapTotalBits = (Public64bit*GoosWindows)*35 + (Public64bit*(1-GoosWindows)*(1-GoosDarwin*GoarchArm64))*39 + GoosDarwin*GoarchArm64*31 + (1-Public64bit)*(32-(GoarchMips+GoarchMipsle))
	MHeapMapBits      = MHeapMapTotalBits - PageShift

	// MaxMem is the maximum heap arena size minus 1.
	//
	// On 32-bit, this is also the maximum heap pointer value,
	// since the arena starts at address 0.
	MaxMem = 1<<MHeapMapTotalBits - 1
)

func gostringw(strw *uint16) string {
	var buf [8]byte
	str := (*[MaxMem/2/2 - 1]uint16)(unsafe.Pointer(strw))
	n1 := 0
	for i := 0; str[i] != 0; i++ {
		n1 += Encoderune(buf[:], rune(str[i]))
	}
	s, b := Rawstring(n1 + 4)
	n2 := 0
	for i := 0; str[i] != 0; i++ {
		// check for race
		if n2 >= n1 {
			break
		}
		n2 += Encoderune(b[n2:], rune(str[i]))
	}
	b[n2] = 0 // for luck
	return s[:n2]
}
