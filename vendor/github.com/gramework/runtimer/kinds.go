package runtimer

const (
	KindBool = 1 + iota
	KindInt
	KindInt8
	KindInt16
	KindInt32
	KindInt64
	KindUint
	KindUint8
	KindUint16
	KindUint32
	KindUint64
	KindUintptr
	KindFloat32
	KindFloat64
	KindComplex64
	KindComplex128
	KindArray
	KindChan
	KindFunc
	KindInterface
	KindMap
	KindPtr
	KindSlice
	KindString
	KindStruct
	KindUnsafePointer

	KindDirectIface = 1 << 5
	KindGCProg      = 1 << 6
	KindNoPointers  = 1 << 7
	KindMask        = (1 << 5) - 1
)

// IsDirectIface reports whether t is stored directly in an interface value.
func IsDirectIface(t *Type) bool {
	return t.Kind&KindDirectIface != 0
}
