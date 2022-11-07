package schemas

import (
	"reflect"
	"unsafe"
)

type Payloader interface {
	ToBytes() []byte
}

func Struct2BytesSlice[T Payloader](value T) []byte {
	byts := []byte{}
	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&byts))
	byteHeader.Data = uintptr(unsafe.Pointer(&value))
	byteHeader.Cap = int(unsafe.Sizeof(value))
	byteHeader.Len = int(unsafe.Sizeof(value))
	return *(*[]byte)(unsafe.Pointer(byteHeader))
}
