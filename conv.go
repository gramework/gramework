// Copyright 2017-present Kirill Danshin and Gramework contributors
// Copyright 2019-present Highload LTD (UK CN: 11893420)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"unsafe"
)

// BytesToString effectively converts bytes to string
// nolint: gas
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes effectively converts string to bytes
// nolint: gas
func StringToBytes(s string) []byte {
	strstruct := stringStructOf(&s)
	return *(*[]byte)(unsafe.Pointer(&sliceType2{
		Array: strstruct.Str,
		Len:   strstruct.Len,
		Cap:   strstruct.Len,
	}))
}

type sliceType2 struct {
	Array unsafe.Pointer
	Len   int
	Cap   int
}

type stringStruct struct {
	Str unsafe.Pointer
	Len int
}


func stringStructOf(sp *string) *stringStruct {
	return (*stringStruct)(unsafe.Pointer(sp))
}