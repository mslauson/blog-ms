package controller

import (
	"net/http"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/gin-gonic/gin"

	"gitea.slauson.io/slausonio/iam-ms/service"
)

type SessionController struct {
	s service.IamSessionService
}

//go:generate mockery --name IamSessionController
type IamSessionController interface {
	CreateEmailSession(c *gin.Context)
	DeleteSession(c *gin.Context)
}

func NewSessionController() *SessionController {
	return &SessionController{
		s: service.NewSessionService(),
	}
}

func (sc *SessionController) CreateEmailSession(c *gin.Context) {
	request := new(siogeneric.AwEmailSessionRequest)
	err := sioUtils.DecryptAndHandle(request, c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response, err := sc.s.CreateEmailSession(request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, response)
	return
}

func (sc *SessionController) DeleteSession(c *gin.Context) {
	sessionId := c.Param("sessionId")
	response, err := sc.s.DeleteSession(sessionId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
	return
}
