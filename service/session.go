package service

import (
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
	"gitea.slauson.io/slausonio/iam-ms/client"
	"gitea.slauson.io/slausonio/iam-ms/constants"
)

type SessionService struct {
	awClient client.AppwriteClient
}

//go:generate mockery --name IamSessionService
type IamSessionService interface {
	CreateEmailSession(
		r *siogeneric.AwEmailSessionRequest,
	) (*siogeneric.AwSession, error)
	DeleteSession(ID, sID string) (siogeneric.SuccessResponse, error)
}

func NewSessionService() *SessionService {
	return &SessionService{
		awClient: client.NewAwClient(),
	}
}

func (s *SessionService) CreateEmailSession(
	r *siogeneric.AwEmailSessionRequest,
) (*siogeneric.AwSession, error) {
	response, err := s.awClient.CreateEmailSession(r)
	if err != nil {
		return nil, sioerror.NewSioUnauthorizedError(err.Error())
	}
	return response, nil
}

func (s *SessionService) DeleteSession(
	ID string,
	sID string,
) (siogeneric.SuccessResponse, error) {
	err := s.awClient.DeleteSession(ID, sID)
	if err != nil {
		return siogeneric.SuccessResponse{
				Success: false,
			}, sioerror.NewSioNotFoundError(
				constants.NoUserFound,
			)
	}

	return siogeneric.SuccessResponse{Success: true}, nil
}
