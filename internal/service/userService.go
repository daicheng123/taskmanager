package service

import (
	"github.com/gin-gonic/gin"
	"taskmanager/internal/mapper"
	"taskmanager/internal/models"
)

//type UserService interface {
//	Login() error
//	LogOut()
//}

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us *UserService) Valid(ctx *gin.Context) {

}

func (us *UserService) GetUserByEmail(email string) (user *models.UserModel, err error) {
	return mapper.GetUserMapper().FindByEmail(email)
}
