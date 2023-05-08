package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
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

// @Summary Create Email Session
// POST
// @Tags session
// @Accept  json
// @Produce  json
// @Param sessionRequest body siogeneric.AwEmailSessionRequest true "Session Request"
// @Success 200 {object} siogeneric.AwSession
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/iam/v1/session [post]
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
}

// @Summary Delete Session
// PUT
// @Tags session
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} siogeneric.SuccessResponse
// @Failure 400 {object} siogeneric.ErrorResponse
// @Failure 401 {object} siogeneric.ErrorResponse
// @Failure 404 {object} siogeneric.ErrorResponse
// @Failure 500 {object} siogeneric.ErrorResponse
// @Router /api/iam/v1/session/:id/:sessionId [delete]
func (sc *SessionController) DeleteSession(c *gin.Context) {
	ID := c.Param("id")
	sessionID := c.Param("sessionId")
	response, err := sc.s.DeleteSession(ID, sessionID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
