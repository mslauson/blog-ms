package service

import (
	"net/http"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	"gitea.slauson.io/slausonio/iam-ms/client"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	c        *gin.Context
	awClient *client.AwClient
}

type IamUserService interface {
	ListUsers() (*sioModel.AwlistResponse, error)
	GetUserByID(id string) (*sioModel.AwUser, error)
	CreateUser(r *sioModel.AwCreateUserRequest) (*sioModel.AwUser, error)
	UpdateEmail(id string, r *sioModel.UpdateEmailRequest) (*sioModel.AwUser, error)
	UpdatePhone(id string, r *sioModel.UpdatePhoneRequest) (*sioModel.AwUser, error)
	UpdatePassword(id string, r *sioModel.UpdatePasswordRequest) (*sioModel.AwUser, error)
	DeleteUser(id string) error
}

func NewUserService(c *gin.Context) *UserService {
	return &UserService{
		c:        c,
		awClient: client.NewAwClient(),
	}
}

func (s *UserService) ListUsers() (*sioModel.AwlistResponse, error) {
	response, err := s.awClient.ListUsers()
	if err != nil {
		s.c.AbortWithError(http.StatusNotFound, iamConst.NoCustomersFound)
	}
}
