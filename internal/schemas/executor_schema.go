package schemas

import "taskmanager/internal/models/common"

type AccountInfo struct {
	AccountName     string `json:"accountName" form:"accountName" binding:"required"`
	AccountPassword string `json:"accountPassword" form:"accountPassword" binding:"required"`
}

type ExecutorInfo struct {
	ID             uint              `json:"id,omitempty"`
	CreatedAt      common.CustomTime `json:"createdAt"`
	UpdatedAt      common.CustomTime `json:"updatedAt"`
	LastOperator   string            `json:"lastOperator"`
	HostName       string            `json:"hostName"`
	IPAddr         string            `json:"ipAddr"`
	SSHPort        uint              `json:"sshPort"`
	Status         uint              `json:"status"`
	SecretStatus   uint              `json:"secretStatus"`
	*AccountInfo   `json:",inline"`
	ExecutePath    string `json:"executePath"`
	Remarks        string `json:"remarks"`
	ShouldDelivery bool
}
