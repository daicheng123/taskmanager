package utils

import "reflect"

//Size 判断切片/map 长度
func Size(s interface{}) int {
	//t := reflect.TypeOf(s)
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
