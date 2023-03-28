package service

import (
	"net/http"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	sioModelGeneric "gitea.slauson.io/slausonio/go-libs/model/generic"
	"gitea.slauson.io/slausonio/iam-ms/client"
	"gitea.slauson.io/slausonio/iam-ms/iamError"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	c        *gin.Context
	awClient *client.AwClient
}

type IamUserService interface {
	ListUsers() *sioModel.AwlistResponse
	GetUserByID(id string) *sioModel.AwUser
	CreateUser(r *sioModel.AwCreateUserRequest) *sioModel.AwUser
	UpdateEmail(id string, r *sioModel.UpdateEmailRequest) *sioModel.AwUser
	UpdatePhone(id string, r *sioModel.UpdatePhoneRequest) *sioModel.AwUser
	UpdatePassword(id string, r *sioModel.UpdatePasswordRequest) *sioModel.AwUser
	DeleteUser(id string) sioModelGeneric.SuccessResponse
}

func NewUserService(c *gin.Context) *UserService {
	return &UserService{
		c:        c,
		awClient: client.NewAwClient(),
	}
}

func (s *UserService) ListUsers() *sioModel.AwlistResponse {
	response, err := s.awClient.ListUsers()
	if err != nil {
		s.c.AbortWithError(http.StatusNotFound, iamError.NoCustomersFound)
	}

	return response
}

func (s *UserService) GetUserByID(id string) *sioModel.AwUser {
	response, err := s.awClient.GetUserByID(id)
	if err != nil {
		s.c.AbortWithError(http.StatusNotFound, iamError.NoCustomerFound)
	}

	return response
}

func (s *UserService) CreateUser(r *sioModel.AwCreateUserRequest) *sioModel.AwUser {
	response, err := s.awClient.CreateUser(r)
	if err != nil {
		s.c.AbortWithError(http.StatusBadRequest, err)
	}

	return response
}

func (s *UserService) UpdateEmail(id string, r *sioModel.UpdateEmailRequest) *sioModel.AwUser {
	response, err := s.awClient.UpdateEmail(id, r)
	if err != nil {
		// TODO: check for 404 and return 404
		s.c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}
	return response
}

func (s *UserService) UpdatePhone(id string, r *sioModel.UpdatePhoneRequest) *sioModel.AwUser {
	response, err := s.awClient.UpdatePhone(id, r)
	if err != nil {
		// TODO: check for 404 and return 404
		s.c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}
	return response
}

func (s *UserService) UpdatePassword(id string, r *sioModel.UpdatePasswordRequest) *sioModel.AwUser {
	response, err := s.awClient.UpdatePassword(id, r)
	if err != nil {
		// TODO: check for 404 and return 404
		s.c.AbortWithError(http.StatusBadRequest, err)
		return nil
	}
	return response
}

func (s *UserService) DeleteUser(id string) sioModelGeneric.SuccessResponse {
	err := s.awClient.DeleteUser(id)
	if err != nil {
		s.c.AbortWithError(http.StatusNotFound, iamError.NoCustomerFound)
		return sioModelGeneric.SuccessResponse{Success: true}
	}

	return sioModelGeneric.SuccessResponse{Success: true}
}
