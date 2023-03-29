package service

import (
	"gitea.slauson.io/slausonio/iam-ms/iamError"
	"net/http"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	sioModelGeneric "gitea.slauson.io/slausonio/go-libs/model/generic"
	"gitea.slauson.io/slausonio/iam-ms/client"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type SessionService struct {
	awClient client.AppwriteClient
}

//go:generate mockery --name IamSessionService
type IamSessionService interface {
	CreateEmailSession(r *sioModel.AwEmailSessionRequest, c *gin.Context) *sioModel.AwSession
	DeleteSession(sID string, c *gin.Context) sioModelGeneric.SuccessResponse
}

func NewSessionService() *SessionService {
	return &SessionService{
		awClient: client.NewAwClient(),
	}
}

func (s *SessionService) CreateEmailSession(r *sioModel.AwEmailSessionRequest, c *gin.Context) *sioModel.AwSession {
	response, err := s.awClient.CreateEmailSession(r)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}
	return response
}

func (s *SessionService) DeleteSession(sID string, c *gin.Context) sioModelGeneric.SuccessResponse {
	err := s.awClient.DeleteSession(sID)
	if err != nil {
		log.Error(err)
		c.AbortWithError(http.StatusNotFound, iamError.NoCustomerFound)
		return sioModelGeneric.SuccessResponse{Success: false}
	}

	return sioModelGeneric.SuccessResponse{Success: true}
}
