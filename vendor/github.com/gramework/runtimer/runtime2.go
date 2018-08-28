package runtimer

// layout of Itab known to compilers
// allocated in non-garbage-collected memory
// Needs to be in sync with
// ../cmd/compile/internal/gc/reflect.go:/^func.dumptypestructs.
type Itab struct {
	Inter  *InterfaceType
	Type   *Type
	Link   *Itab
	Hash   uint32 // copy of _type.hash. Used for type switches.
	Bad    bool   // type does not implement interface
	Inhash bool   // has this itab been added to hash?
	Unused [2]byte
	Fun    [1]uintptr // variable sized
}
