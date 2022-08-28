package admin

import "taskmanager/utils/serializer"

type ExecutorService struct {
	Hostname       string          `json:"hostname" binding:"omitempty,required"`
	IPAddr         string          `json:"ipAddr"   binding:"omitempty,required"`
	SHHPort        uint            `json:"sshPort"  binding:"omitempty,required,min=1,max=65535"`
	Account        *AccountService `json:",inline"  binding:"omitempty,required"`
	ExecutePath    string          `json:"executePath" binding:"omitempty,required"`
	ShouldDelivery bool            `json:"shouldDelivery"`
	LastOperator   string          `json:"lastOperator" binding:"omitempty,required"`
}

func (es *ExecutorService) ExecutorAdd() serializer.Response {

	return serializer.Response{}
}
