package service

import (
	"net/http"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	sioModelGeneric "gitea.slauson.io/slausonio/go-libs/model/generic"
	"gitea.slauson.io/slausonio/iam-ms/client"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type SessionService struct {
	c        *gin.Context
	awClient *client.AwClient
}

type IamSessionService interface {
	CreateEmailSession(r *sioModel.AwEmailSessionRequest) *sioModel.AwSession
	DeleteSession(sID string) sioModelGeneric.SuccessResponse
}

func NewSessionService(c *gin.Context) *SessionService {
	return &SessionService{
		c:        c,
		awClient: client.NewAwClient(),
	}
}

func (s *SessionService) CreateEmailSession(r *sioModel.AwEmailSessionRequest) *sioModel.AwSession {
	response, err := s.awClient.CreateEmailSession(r)
	if err != nil {
		log.Error(err)
		s.c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}
	return response
}

func (s *SessionService) DeleteSession(sID string) sioModelGeneric.SuccessResponse           {
	err := s.awClient.DeleteSession(sID)
	if err != nil {
		log.Error(err)
		s.c.AbortWithError(http.StatusBadRequest, err)
		return sioModelGeneric.SuccessResponse{Success: true}
	}

	return sioModelGeneric.SuccessResponse{Success: true}
}
