package controller

import (
	sioModelGeneric "gitea.slauson.io/slausonio/go-libs/model/generic"
	"net/http"
	"net/http/httptest"
	"testing"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"gitea.slauson.io/slausonio/iam-ms/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

var mUserSession = &sioModel.AwSession{
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

func initControllerForSessionTests(t *testing.T) (*SessionController, *mocks.IamSessionService, *sioUtils.EncryptionUtil) {
	ss := mocks.NewIamSessionService(t)
	eu := sioUtils.NewEncryptionUtil()
	sc := &SessionController{
		s: ss,
	}

	return sc, ss, eu
}

// func TestNewSessionController(t *testing.T) {
// 	c, ms, eu := initControllerForSessionTests(t)
// 	tests := []struct {
// 		name string
// 		want *SessionController
// 	}{
// 		{name}
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewSessionController(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewSessionController() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestSessionController_CreateEmailSession(t *testing.T) {
	sc, ss, eu := initControllerForSessionTests(t)

	tests := []struct {
		name    string
		request *sioModel.AwEmailSessionRequest
		want    *sioModel.AwSession
		status  int
	}{
		{name: "valid", request: &sioModel.AwEmailSessionRequest{Email: "test@test.com", Password: "asdf"}, want: mUserSession, status: http.StatusOK},
		{name: "service error", request: &sioModel.AwEmailSessionRequest{Email: "test@test.com", Password: "asdf"}, want: nil, status: http.StatusOK},
		{name: "Missing Email", request: &sioModel.AwEmailSessionRequest{Email: "", Password: "asdf"}, want: nil, status: http.StatusBadRequest},
		{name: "Missing Pass", request: &sioModel.AwEmailSessionRequest{Email: "test@test.com", Password: ""}, want: nil, status: http.StatusBadRequest},
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

			err := eu.EncryptInterface(tt.request, false)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "POST")
			if tt.status == http.StatusOK {
				ss.On("CreateEmailSession", mock.AnythingOfType("*sioModel.AwEmailSessionRequest"), c).Return(tt.want, nil)
			}
			sc.CreateEmailSession(c)
			if w.Code != tt.status {
				t.Errorf("CreateEmailSession() status = %v, want %v", w.Code, tt.status)
				return
			}
		})
	}
}

func TestSessionController_DeleteSession(t *testing.T) {
	sc, ss, _ := initControllerForSessionTests(t)
	tests := []struct {
		name   string
		want   sioModelGeneric.SuccessResponse
		status int
	}{
		{name: "happy", want: sioModelGeneric.SuccessResponse{Success: true}, status: 200},
		{name: "service error", want: sioModelGeneric.SuccessResponse{Success: false}, status: 200},
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

			c.Params = gin.Params{gin.Param{Key: "sessionId", Value: "a"}}
			ss.On("DeleteSession", "a", c).Return(tt.want)
			sc.DeleteSession(c)
		})
	}
}
