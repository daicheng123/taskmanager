package utils

import (
	"reflect"
	"unsafe"
)

//Size 判断切片/map 长度
func Size(s interface{}) int {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr || !v.IsNil() {
		v = v.Elem()
	}
	k := v.Kind()
	if k != reflect.Map && k != reflect.Slice {
		return -1
	}
	return v.Len()
}

// IsZero 校验是否为0值
func IsZero(v interface{}) bool {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value.IsZero()
}

func Struct2BytesSlice(value interface{}) []byte {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		byts := make([]byte, 0)
		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(value.(*struct{})))
		byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&byts))
		byteHeader.Data = sliceHeader.Data
		byteHeader.Len = sliceHeader.Len
		byteHeader.Cap = sliceHeader.Len
		return *(*[]byte)(unsafe.Pointer(byteHeader))
		//return
	}
	return nil
}
