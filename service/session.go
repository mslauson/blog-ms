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
	CreateEmailSession(r *sioModel.AwEmailSessionRequest) (*sioModel.AwSession, error)
	DeleteSession(sID string) (sioModelGeneric.SuccessResponse, error)
}

func NewSessionService(c *gin.Context) *SessionService {
	return &SessionService{
		c:        c,
		awClient: client.NewAwClient(),
	}
}

func (s *SessionService) CreateEmailSession(r *sioModel.AwEmailSessionRequest) (*sioModel.AwSession, error) {
	response, err := s.awClient.CreateEmailSession(r)
	if err != nil {
		log.Error(err)
		s.c.AbortWithError(http.StatusBadRequest, err)
		return nil, err
	}
	return response, nil
}

func (s *SessionService) DeleteSession(sID string) (sioModelGeneric.SuccessResponse, error) {
	err := s.awClient.DeleteSession(sID)
	if err != nil {
		log.Error(err)
		s.c.AbortWithError(http.StatusBadRequest, err)
		return sioModelGeneric.SuccessResponse{Success: true}, err
	}

	return sioModelGeneric.SuccessResponse{Success: true}, nil
}
