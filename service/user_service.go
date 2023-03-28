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

	return response, nil
}

func (s *UserService) GetUserByID(id string) (*sioModel.AwUser, error) {
	response, err := s.awClient.GetUserByID(id)
	if err != nil {
		s.c.AbortWithError(http.StatusNotFound, iamConst.NoCustomerFound)
	}

	return response, nil
}

func (s *UserService) CreateUser(r *sioModel.AwCreateUserRequest) (*sioModel.AwUser, error) {
	response, err := s.awClient.CreateUser(r)
	if err != nil {
		s.c.AbortWithError(http.StatusBadRequest, err)
	}

	return response, nil
}

func (s *UserService) UpdateEmail(id string, r *sioModel.UpdateEmailRequest) (*sioModel.AwUser, error) {
	response, err := s.awClient.UpdateEmail(id, r)
	if err != nil {
		// TODO: check for 404 and return 404
		s.c.AbortWithError(http.StatusBadRequest, err)
		return nil, err
	}
	return response, nil
}

func (s *UserService) UpdatePhone(id string, r *sioModel.UpdatePhoneRequest) (*sioModel.AwUser, error) {
	response, err := s.awClient.UpdatePhone(id, r)
	if err != nil {
		// TODO: check for 404 and return 404
		s.c.AbortWithError(http.StatusBadRequest, err)
		return nil, err
	}
	return response, nil
}

func (s *UserService) UpdatePassword(id string, r *sioModel.UpdatePasswordRequest) (*sioModel.AwUser, error) {
	response, err := s.awClient.UpdatePassword(id, r)
	if err != nil {
		// TODO: check for 404 and return 404
		s.c.AbortWithError(http.StatusBadRequest, err)
		return nil, err
	}
	return response, nil
}

func (s *UserService) DeleteUser(id string) (sioModel.GenericSuccessResposne, error) {
	err := s.awClient.DeleteUser(id)
	if err != nil {
		s.c.AbortWithError(http.StatusNotFound, iamConst.NoCustomerFound)
	}

	return sioModel.GenericSuccessResposne{true}
}
