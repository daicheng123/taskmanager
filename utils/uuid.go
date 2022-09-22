package utils

import "github.com/google/uuid"

func NewUuid() string {
	u, _ := uuid.NewUUID()
	return u.String()
}
