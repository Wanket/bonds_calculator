package util

import (
	"reflect"
	"unsafe"
)

func StringToBytes(s string) []byte {
	var bytes []byte

	(*reflect.SliceHeader)(unsafe.Pointer(&bytes)).Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	(*reflect.SliceHeader)(unsafe.Pointer(&bytes)).Cap = len(s)
	(*reflect.SliceHeader)(unsafe.Pointer(&bytes)).Len = len(s)

	return bytes
}

func ByteArrayToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
