package payload

// OperationPayloadReq 脚本执行载荷
type OperationPayloadReq struct {
	Servers        []uint `json:"servers"`
	ScriptID       uint   `json:"scriptId"`
	ScriptContent  string
	ScriptTypeId   uint
	ScriptOverTime int
	TaskName       string `json:"taskName"`
	UniqueTag      string `json:"uniqueTag"`
	TaskOperator   string `json:"taskOperator"`
	TaskId         string
}

func (op *OperationPayloadReq) ToBytes() []byte {
	return Struct2BytesSlice[*OperationPayloadReq](op)
}
