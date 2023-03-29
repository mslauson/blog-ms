package controller

import (
	"net/http"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	"gitea.slauson.io/slausonio/iam-ms/service"
	"github.com/gin-gonic/gin"
)

type SessionController struct {
	s service.IamSessionService
}

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
	var request sioModel.AwEmailSessionRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	response := sc.s.CreateEmailSession(&request, c)
	if response != nil {
		c.JSON(http.StatusOK, response)
	}
}

func (sc *SessionController) DeleteSession(c *gin.Context) {
	sessionId := c.Param("sessionId")
	response := sc.s.DeleteSession(sessionId, c)

	if response.Success {
		c.JSON(http.StatusOK, response)
	}
}
