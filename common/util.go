package common

import (
	"reflect"
	"unsafe"
)

func UnsafeBytesToString(bytes []byte) string {
	hdr := &reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&bytes[0])),
		Len:  len(bytes),
	}
	return *(*string)(unsafe.Pointer(hdr))
}
