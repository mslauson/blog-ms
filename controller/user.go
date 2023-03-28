package controller

import (
	"gitea.slauson.io/slausonio/iam-ms/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	s service.IamUserService
}

type IamUserController interface {
	ListUsers(c *gin.Context)
	GetUserById(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdatePassword(c *gin.Context)
	UpdateEmail(c *gin.Context)
	UpdatePhone(c *gin.Context)
	DeleteUser(c *gin.Context)
}

func NewUserController(s service.IamUserService) *UserController {
	return &UserController{
		s: service.NewUserService(),
	}
}
