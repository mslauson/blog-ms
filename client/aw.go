package client

import (
	sioRest "gitea.slauson.io/slausonio/go-utils/rest"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
)

type AwClient struct {
	h    sioRest.RestHelpers
	host string
}

type AppwriteClient interface {
	ListUsers() (*sioModel.AwlistResponse, error)
	GetUserByID(id string) (*sioModel.AwUser, error)
	CreateUser(r *sioModel.AwCreateUserRequest) (*sioModel.AwUser, error)
	UpdateEmail(id string, r *sioModel.UpdateEmailRequest) (*sioModel.AwUser, error)
	UpdatePhone(id string, r *sioModel.UpdatePhoneRequest) (*sioModel.AwUser, error)
	UpdatePassword(id string, r *sioModel.UpdatePasswordRequest) (*sioModel.AwUser, error)

	CreateSession(r *sioModel.AwEmailSessionRequest) (*sioModel.AwSession, error)
	DeleteSession(sID string) error
}
