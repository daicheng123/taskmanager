package payload

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

	//var op2 = *(**T)(unsafe.Pointer(&data))
	//return op2

	//var load = *(**T)(unsafe.Pointer(&data))
	//fmt.Printf("%#v\n", *op2)
	//return *op2
}