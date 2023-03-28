package controller

import (
	"gitea.slauson.io/slausonio/iam-ms/service"
	"github.com/gin-gonic/gin"
)

type SessionController struct {
	s service.IamSessionService
}

type IamSessionController interface {
	CreateEmailSession(c *gin.Context)
	DeleteEmailSession(c *gin.Context)
}

func NewSessionController(c *gin.Context) *SessionController {
	return &SessionController{
		s: service.NewSessionService(c),
	}
}
