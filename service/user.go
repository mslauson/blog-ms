package service

import (
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioerror"

	"gitea.slauson.io/slausonio/iam-ms/client"
	"gitea.slauson.io/slausonio/iam-ms/constants"
)

type UserService struct {
	awClient client.AppwriteClient
}

//go:generate mockery --name IamUserService
type IamUserService interface {
	ListUsers() (*siogeneric.AwlistResponse, error)
	GetUserByID(id string) (*siogeneric.AwUser, error)
	CreateUser(r *siogeneric.AwCreateUserRequest) (*siogeneric.AwUser, error)
	UpdateEmail(id string, r *siogeneric.UpdateEmailRequest) (*siogeneric.AwUser, error)
	UpdatePhone(id string, r *siogeneric.UpdatePhoneRequest) (*siogeneric.AwUser, error)
	UpdatePassword(
		id string,
		r *siogeneric.UpdatePasswordRequest,
	) (*siogeneric.AwUser, error)
	DeleteUser(id string) (siogeneric.SuccessResponse, error)
}

func NewUserService() *UserService {
	return &UserService{
		awClient: client.NewAwClient(),
	}
}

func (s *UserService) ListUsers() (*siogeneric.AwlistResponse, error) {
	response, err := s.awClient.ListUsers()
	if err != nil {
		return nil, sioerror.NewSioNotFoundError(constants.NoCustomersFound)
	}

	return response, nil
}

func (s *UserService) GetUserByID(id string) (*siogeneric.AwUser, error) {
	response, err := s.awClient.GetUserByID(id)
	if err != nil {
		return nil, sioerror.NewSioNotFoundError(constants.NoCustomerFound)
	}

	return response, nil
}

func (s *UserService) CreateUser(
	r *siogeneric.AwCreateUserRequest,
) (*siogeneric.AwUser, error) {
	response, err := s.awClient.CreateUser(r)
	if err != nil {
		return nil, sioerror.NewSioBadRequestError(err.Error())
	}

	return response, nil
}

func (s *UserService) UpdateEmail(
	id string,
	r *siogeneric.UpdateEmailRequest,
) (*siogeneric.AwUser, error) {
	response, err := s.awClient.UpdateEmail(id, r)
	if err != nil {
		return nil, sioerror.NewSioBadRequestError(err.Error())
	}
	return response, nil
}

func (s *UserService) UpdatePhone(
	id string,
	r *siogeneric.UpdatePhoneRequest,
) (*siogeneric.AwUser, error) {
	response, err := s.awClient.UpdatePhone(id, r)
	if err != nil {
		return nil, sioerror.NewSioBadRequestError(err.Error())
	}
	return response, nil
}

func (s *UserService) UpdatePassword(
	id string,
	r *siogeneric.UpdatePasswordRequest,
) (*siogeneric.AwUser, error) {
	response, err := s.awClient.UpdatePassword(id, r)
	if err != nil {
		return nil, sioerror.NewSioBadRequestError(err.Error())
	}
	return response, nil
}

func (s *UserService) DeleteUser(id string) (siogeneric.SuccessResponse, error) {
	err := s.awClient.DeleteUser(id)
	if err != nil {
		return siogeneric.SuccessResponse{Success: false}, sioerror.NewSioNotFoundError(err.Error())
	}

	return siogeneric.SuccessResponse{Success: true}, nil
}
