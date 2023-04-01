package service

import (
	"gitea.slauson.io/slausonio/iam-ms/client/mocks"
	"gitea.slauson.io/slausonio/iam-ms/iamError"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	siotest "gitea.slauson.io/slausonio/go-testing/sio_test"
	"github.com/stretchr/testify/assert"
)

// func TestNewUserService(t *testing.T) {
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
	sessionReq = &sioModel.AwEmailSessionRequest{
		Email:    "test",
		Password: "test",
	}

	mUserSession = &sioModel.AwSession{
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
	gc, _ := siotest.InitGinTestContext()

	awClient.On("CreateEmailSession", mock.AnythingOfType("*sioModel.AwEmailSessionRequest")).Return(mUserSession, nil)
	actual := ss.CreateEmailSession(sessionReq, gc)
	assert.Equalf(t, mUserSession, actual, "actual: %v", actual)
	assert.Emptyf(t, gc.Errors, "gc.Errors: %v", gc.Errors)
}

func TestSessionService_CreateUser_Error(t *testing.T) {
	ss, awClient := initSessionServiceTest(t)
	gc, w := siotest.InitGinTestContext()

	awClient.On("CreateEmailSession", mock.AnythingOfType("*sioModel.AwEmailSessionRequest")).Return(nil, tError)
	actual := ss.CreateEmailSession(sessionReq, gc)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, gc.Errors[0].Err, tError, "error: %v", gc.Errors[0].Err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSessionService_DeleteSession(t *testing.T) {
	ss, awClient := initSessionServiceTest(t)
	gc, _ := siotest.InitGinTestContext()

	awClient.On("DeleteSession", "a").Return(nil)
	actual := ss.DeleteSession("a", gc)
	assert.Truef(t, actual.Success, "actual.Success: %v", actual.Success)
	assert.Emptyf(t, gc.Errors, "gc.Errors: %v", gc.Errors)
}

func TestSessionService_DeleteSession_Error(t *testing.T) {
	ss, awClient := initSessionServiceTest(t)
	gc, w := siotest.InitGinTestContext()

	awClient.On("DeleteSession", "a").Return(tError)
	actual := ss.DeleteSession("a", gc)
	assert.False(t, actual.Success)
	assert.Equalf(t, gc.Errors[0].Err, iamError.NoCustomerFound, "error: %v", gc.Errors[0].Err)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
