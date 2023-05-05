package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitea.slauson.io/slausonio/go-types/siogeneric"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"gitea.slauson.io/slausonio/iam-ms/service/mocks"
)

var mUserSession = &siogeneric.AwSession{
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

func initControllerForSessionTests(
	t *testing.T,
) (*SessionController, *mocks.IamSessionService, *sioUtils.EncryptionUtil) {
	ss := mocks.NewIamSessionService(t)
	eu := sioUtils.NewEncryptionUtil()
	sc := &SessionController{
		s: ss,
	}

	return sc, ss, eu
}

func TestNewSessionController(t *testing.T) {
	sc := NewSessionController()
	assert.NotNil(t, sc)
}

func TestSessionController_CreateEmailSession(t *testing.T) {
	tests := []struct {
		name    string
		request *siogeneric.AwEmailSessionRequest
		want    *siogeneric.AwSession
		status  int
	}{
		{
			name:    "valid",
			request: &siogeneric.AwEmailSessionRequest{Email: "test@test.com", Password: "asdf"},
			want:    mUserSession,
			status:  http.StatusOK,
		},
		{
			name:    "Missing Email",
			request: &siogeneric.AwEmailSessionRequest{Email: "", Password: "asdf"},
			want:    nil,
			status:  http.StatusBadRequest,
		},
		{
			name:    "Missing Pass",
			request: &siogeneric.AwEmailSessionRequest{Email: "test@test.com", Password: ""},
			want:    nil,
			status:  http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				w    = httptest.NewRecorder()
				c, _ = gin.CreateTestContext(w)
			)
			c.Request = &http.Request{
				Header: make(http.Header),
			}

			sc, ss, eu := initControllerForSessionTests(t)

			err := eu.EncryptInterface(tt.request)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "POST")
			if tt.want != nil {
				ss.On("CreateEmailSession", mock.AnythingOfType("*siogeneric.AwEmailSessionRequest")).
					Return(tt.want, nil)
			}
			sc.CreateEmailSession(c)
			if tt.want != nil {
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestUserController_CreateEmailSessionServiceFailure(t *testing.T) {
	request := &siogeneric.AwEmailSessionRequest{Email: "test@test.com", Password: "asdf"}
	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	sc, ss, eu := initControllerForSessionTests(t)

	err := eu.EncryptInterface(request)
	if err != nil {
		t.Error(err)
		return
	}

	MockJson(c, request, "POST")
	ss.On("CreateEmailSession", mock.AnythingOfType("*siogeneric.AwEmailSessionRequest")).
		Return(nil, errors.New("error"))
	sc.CreateEmailSession(c)

	assert.Truef(t, c.Errors != nil, "c.Errors shouldn't be nil")
}

func TestDeleteSession(t *testing.T) {
	uc, ms, _ := initControllerForSessionTests(t)

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}, gin.Param{Key: "sessionId", Value: "a"}}
	ms.On("DeleteSession", "a", "a").Return(siogeneric.SuccessResponse{Success: true}, nil)
	uc.DeleteSession(c)

	assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
}

func TestDeleteSessionError(t *testing.T) {
	uc, ms, _ := initControllerForSessionTests(t)

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}, gin.Param{Key: "sessionId", Value: "a"}}
	ms.On("DeleteSession", "a", "a").
		Return(siogeneric.SuccessResponse{Success: false}, errors.New("asdf"))
	uc.DeleteSession(c)

	assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
}
