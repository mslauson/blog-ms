package controller

import (
	"fmt"
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
	enc := sioUtils.NewEncryptionUtil()
	request := new(siogeneric.AwEmailSessionRequest)

	err := c.ShouldBindJSON(request)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = enc.DecryptInterface(request, false)
	if err != nil {
		c.AbortWithError(
			http.StatusInternalServerError,
			fmt.Errorf("decryption failed - unable to proceed"),
		)
		return
	}
	response := sc.s.CreateEmailSession(request, c)
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
