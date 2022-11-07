package main

import (
	"fmt"
	"time"
)

//type Payload interface {
//	ToBytes() []byte
//}
//type OperationPayload struct {
//	Servers        []uint `json:"servers"`
//	ScriptID       uint   `json:"scriptId"`
//	ScriptContent  string
//	ScriptTypeId   uint
//	ScriptOverTime uint
//	TaskName       string `json:"taskName"`
//	UniqueTag      string `json:"uniqueTag"`
//	TaskOperator   string `json:"taskOperator"`
//	TaskId         string
//}
//
//func (op *OperationPayload) ToBytes() []byte {
//	return utils.Struct2BytesSlice[*OperationPayload](op)
//}
//
//type SliceMock struct {
//	addr uintptr
//	len  int
//	cap  int
//}

func main() {
	fmt.Println(time.Now())
	var cstSh, _ = time.LoadLocation("Asia/Shanghai") //上海
	time.Now().In(cstSh)
	l, _ := time.LoadLocation("")
	time.ParseInLocation()
	fmt.Println(time.Now().In(l))

	//s := time.Now().Format("2006-01-02 15:04:05")

	//t, _ := time.Parse("2006-01-02 15:04:05", s)
	//fmt.Println(s, t)
	//t, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)

	//fmt.Println(t, "\n", time.Now().)

	//time1 := time.Now()
	//time.Sleep(time.Second * 10)
	//time2 := time.Now()
	//fmt.Println(time2.Sub(time1).Seconds())
	//for i := 0; i < 6; i++ {
	//	for j := 0; j <= 4; j++ {
	//		if j == 3 {
	//			break
	//		}
	//		fmt.Println("i ---> j", i, j)
	//	}
	//}
	//op := &OperationPayload{
	//	Servers:       []uint{1, 3},
	//	ScriptID:      1,
	//	ScriptContent: "hello",
	//	ScriptTypeId:  2,
	//}

	//Struct2BytesSlice[Payload](op)

	//Len := unsafe.Sizeof(op)
	//testBytes := &SliceMock{
	//	addr: uintptr(unsafe.Pointer(op)),
	//	cap:  int(Len),
	//	len:  int(Len),
	//}

	//byts := []byte{}
	//byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&byts))
	//byteHeader.Data = uintptr(unsafe.Pointer(op))
	//byteHeader.Cap = int(unsafe.Sizeof(op))
	//byteHeader.Len = int(unsafe.Sizeof(op))
	////
	//data := *(*[]byte)(unsafe.Pointer(byteHeader))
	//fmt.Println("[]byte is : ", data)
	//
	//var op2 *OperationPayload = *(**OperationPayload)(unsafe.Pointer(&data))
	//fmt.Println(op2)
	//r := Struct2BytesSlice[*OperationPayload](op)

	//op1 := &OperationPayload{}
	//
	//err := json.Unmarshal(data, op1)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(op1)

	//r := Struct2BytesSlice(op)
	//fmt.Println(r)

	//type TestStructTobytes struct {
	//	data int64
	//}
	//type SliceMock struct {
	//	addr uintptr
	//	len  int
	//	cap  int
	//}
	//
	//var testStruct = &TestStructTobytes{100}
	//Len := unsafe.Sizeof(*testStruct)
	//testBytes := &SliceMock{
	//	addr: uintptr(unsafe.Pointer(testStruct)),
	//	cap:  int(Len),
	//	len:  int(Len),
	//}
	//data := *(*[]byte)(unsafe.Pointer(testBytes))
	//fmt.Println("[]byte is : ", data)
	//var ptestStruct *TestStructTobytes = *(**TestStructTobytes)(unsafe.Pointer(&data))
	//fmt.Println("ptestStruct.data is : ", ptestStruct.data)
}

//func Struct2BytesSlice[T Payload](value T) T {
//	byts := []byte{}
//	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&byts))
//	byteHeader.Data = uintptr(unsafe.Pointer(&value))
//	byteHeader.Cap = int(unsafe.Sizeof(value))
//	byteHeader.Len = int(unsafe.Sizeof(value))
//	//
//	data := *(*[]byte)(unsafe.Pointer(byteHeader))
//	fmt.Println("[]byte is : ", data)
//
//	var op2 = *(**T)(unsafe.Pointer(&data))
//	fmt.Printf("%#v\n", *op2)
//	return *op2
//	//byts := []byte{}
//	//byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&byts))
//	//byteHeader.Data = uintptr(unsafe.Pointer(&value))
//	//byteHeader.Cap = int(unsafe.Sizeof(value))
//	//byteHeader.Len = int(unsafe.Sizeof(value))
//	//
//	//data := *(*[]byte)(unsafe.Pointer(byteHeader))
//	//fmt.Println("[]byte is : ", data)
//	//
//	//var op2 *OperationPayload = (value)(unsafe.Pointer(&data))
//	//fmt.Println(op2)
//	//v := reflect.ValueOf(value)
//	//if v.Kind() == reflect.Ptr {
//	//	v = v.Elem()
//	//}
//	//if v.Kind() == reflect.Struct {
//	//byts := make([]byte, 0)
//	//if v, ok := value.(T); ok {
//	//	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&v))
//	//	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&byts))
//	//	byteHeader.Data = sliceHeader.Data
//	//	byteHeader.Len = sliceHeader.Len
//	//	byteHeader.Cap = sliceHeader.Len
//	//	return *(*[]byte)(unsafe.Pointer(byteHeader))
//	//}
//
//	//return
//	//}
//}
