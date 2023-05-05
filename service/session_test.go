package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	siotest "gitea.slauson.io/slausonio/go-testing/sio_test"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
	"gitea.slauson.io/slausonio/iam-ms/client/mocks"
	"gitea.slauson.io/slausonio/iam-ms/constants"
)

// Func TestNewUserService(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		want *UserService
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			assert.Equal(t, tt.want, NewUserService())
// 		})
// 	}
// }

var (
	sessionReq = &siogeneric.AwEmailSessionRequest{
		Email:    "test",
		Password: "test",
	}

	mUserSession = &siogeneric.AwSession{
		ID:                        "blah",
		CreatedAt:                 "blah",
		AwClientCode:              "blah",
		AwClientEngine:            "blah",
		AwClientEngineVersion:     "blah",
		AwClientName:              "blah",
		AwClientType:              "blah",
		AwClientVersion:           "blah",
		CountryCode:               "blah",
		CountryName:               "blah",
		Current:                   true,
		DeviceBrand:               "blah",
		DeviceModel:               "blah",
		DeviceName:                "blah",
		Expire:                    "blah",
		Ip:                        "blah",
		OsCode:                    "blah",
		OsName:                    "blah",
		OsVersion:                 "blah",
		Provider:                  "blah",
		ProviderAccessToken:       "blah",
		ProviderAccessTokenExpiry: "blah",
		ProviderRefreshToken:      "blah",
		ProviderUid:               "blah",
		UserId:                    "blah",
	}
)

func initSessionServiceTest(t *testing.T) (*SessionService, *mocks.AppwriteClient) {
	ac := mocks.NewAppwriteClient(t)
	ss := &SessionService{
		awClient: ac,
	}
	return ss, ac
}

func TestSessionService_CreateUser(t *testing.T) {
	ss, awClient := initSessionServiceTest(t)

	awClient.On("CreateEmailSession", mock.AnythingOfType("*siogeneric.AwEmailSessionRequest")).Return(mUserSession, nil)
	actual, err := ss.CreateEmailSession(sessionReq)
	assert.Equalf(t, mUserSession, actual, "actual: %v", actual)
	assert.Emptyf(t, err, "err: %v", err)
}

func TestSessionService_CreateUser_Error(t *testing.T) {
	ss, awClient := initSessionServiceTest(t)

	awClient.On("CreateEmailSession", mock.AnythingOfType("*siogeneric.AwEmailSessionRequest")).Return(nil, siotest.TError)
	actual, err := ss.CreateEmailSession(sessionReq)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, err.Error(), siotest.TUnauthorizedError.Error(), "error: %v", err.Error())
}

func TestSessionService_DeleteSession(t *testing.T) {
	ss, awClient := initSessionServiceTest(t)

	awClient.On("DeleteSession", "a", "a").Return(nil)
	actual, err := ss.DeleteSession("a", "a")
	assert.Truef(t, actual.Success, "actual.Success: %v", actual.Success)
	assert.Emptyf(t, err, "err: %v", err)
}

func TestSessionService_DeleteSession_Error(t *testing.T) {
	ss, awClient := initSessionServiceTest(t)

	awClient.On("DeleteSession", "a", "a").Return(siotest.TError)
	actual, err := ss.DeleteSession("a", "a")
	assert.False(t, actual.Success)
	assert.Equalf(t, err.Error(), sioerror.NewSioNotFoundError(constants.NoUserFound).Error(), "error: %v", err.Error())
}
