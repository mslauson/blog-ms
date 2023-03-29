package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	sioModelGeneric "gitea.slauson.io/slausonio/go-libs/model/generic"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	"gitea.slauson.io/slausonio/iam-ms/service/mocks"
)

var mAwUser = sioModel.AwUser{
	Email: "t@t.com",
}

var (
	mAwUserPtr = &mAwUser
	mUserList  = &sioModel.AwlistResponse{
		Total: 1,
		Users: []sioModel.AwUser{mAwUser},
	}
)

func initController(t *testing.T) (*UserController, *mocks.IamUserService, *sioUtils.EncryptionUtil) {
	ms := mocks.NewIamUserService(t)
	eu := sioUtils.NewEncryptionUtil()
	controller := &UserController{
		s: ms,
	}

	return controller, ms, eu
}

func MockJson(c *gin.Context, content interface{}, method string) {
	c.Request.Method = method // or PUT
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func TestListUsers(t *testing.T) {
	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)

	uc, ms, _ := initController(t)
	tests := []struct {
		name   string
		want   *sioModel.AwlistResponse
		status int
	}{
		{name: "happy", want: mUserList, status: 200},
		{name: "service error", want: mUserList, status: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.On("ListUsers", c).Return(tt.want)
			uc.ListUsers(c)

			if w.Code != tt.status {
				t.Errorf("ListUsers() = %v, want %v", w.Code, tt.status)
				return
			}
		})
	}
}

func TestUserController_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		request *sioModel.AwCreateUserRequest
		status  int
		result  *sioModel.AwUser
	}{
		{name: "happy", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusOK, result: mAwUserPtr},
		{name: "service failure", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusOK, result: nil},
		{name: "Bad Request - No Email", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Request - bad Email", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "a(", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Request - bad too long", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "No Password", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: ""}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Number", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Upper", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "attesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Special", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "Mattesting1"}, status: http.StatusBadRequest, result: nil},
		{name: "Short Password", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "Ma"}, status: http.StatusBadRequest, result: nil},
		{name: "No Name", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "long Name", request: &sioModel.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "No Phone", request: &sioModel.AwCreateUserRequest{Phone: "", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Phone short", request: &sioModel.AwCreateUserRequest{Phone: "131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Phone long", request: &sioModel.AwCreateUserRequest{Phone: "212121123123122131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
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

			uc, ms, eu := initController(t)

			err := eu.EncryptInterface(tt.request, false)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "POST")
			if tt.status == http.StatusOK {
				ms.On("CreateUser", mock.AnythingOfType("*sioModel.AwCreateUserRequest"), c).Return(tt.result)
			}
			uc.CreateUser(c)
			if w.Code != tt.status {
				t.Errorf("CreateUser() = %v, want %v", w.Code, tt.status)
				return
			}
		})
	}
}

func TestUserController_UpdatePassword(t *testing.T) {
	tests := []struct {
		name    string
		request *sioModel.UpdatePasswordRequest
		status  int
		result  *sioModel.AwUser
	}{
		{name: "happy", request: &sioModel.UpdatePasswordRequest{Password: "Mm112a23!"}, status: http.StatusOK, result: mAwUserPtr},
		{name: "service failure", request: &sioModel.UpdatePasswordRequest{Password: "Mm112a23!"}, status: http.StatusOK, result: nil},
		{name: "No Password", request: &sioModel.UpdatePasswordRequest{Password: ""}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Number", request: &sioModel.UpdatePasswordRequest{Password: "Mmsdfsdfafsdf!"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Upper", request: &sioModel.UpdatePasswordRequest{Password: "aaam112a23!"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Special", request: &sioModel.UpdatePasswordRequest{Password: "Mm112a2333"}, status: http.StatusBadRequest, result: nil},
		{name: "Short Password", request: &sioModel.UpdatePasswordRequest{Password: "Mm1^"}, status: http.StatusBadRequest, result: nil},
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

			c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

			uc, ms, eu := initController(t)

			err := eu.EncryptInterface(tt.request, false)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "PUT")
			if tt.status == http.StatusOK {
				ms.On("UpdatePassword", "a", mock.AnythingOfType("*sioModel.UpdatePasswordRequest"), c).Return(tt.result)
			}
			uc.UpdatePassword(c)
			if w.Code != tt.status {
				t.Errorf("UpdatePassword() = %v, want %v", w.Code, tt.status)
				return
			}
		})
	}
}

func TestUserController_UpdateEmail(t *testing.T) {
	tests := []struct {
		name    string
		request *sioModel.UpdateEmailRequest
		status  int
		result  *sioModel.AwUser
	}{
		{name: "happy", request: &sioModel.UpdateEmailRequest{Email: "fake@fake.com"}, status: http.StatusOK, result: mAwUserPtr},
		{name: "service failure", request: &sioModel.UpdateEmailRequest{Email: "fake@fake.com"}, status: http.StatusOK, result: nil},
		{name: "No Email", request: &sioModel.UpdateEmailRequest{Email: ""}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Email ", request: &sioModel.UpdateEmailRequest{Email: "fakefake.com"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Email too long", request: &sioModel.UpdateEmailRequest{Email: "\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@.com"}, status: http.StatusBadRequest, result: nil},
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

			c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

			uc, ms, eu := initController(t)

			err := eu.EncryptInterface(tt.request, false)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "PUT")
			if tt.status == http.StatusOK {
				ms.On("UpdateEmail", "a", mock.AnythingOfType("*sioModel.UpdateEmailRequest"), c).Return(tt.result)
			}
			uc.UpdateEmail(c)
			if w.Code != tt.status {
				t.Errorf("UpdateEmail() = %v, want %v", w.Code, tt.status)
				return
			}
		})
	}
}

func TestUserController_UpdatePhone(t *testing.T) {
	tests := []struct {
		name    string
		request *sioModel.UpdatePhoneRequest
		status  int
		result  *sioModel.AwUser
	}{
		{name: "happy", request: &sioModel.UpdatePhoneRequest{Number: "1239323939"}, status: http.StatusOK, result: mAwUserPtr},
		{name: "service failure", request: &sioModel.UpdatePhoneRequest{Number: "1239323939"}, status: http.StatusOK, result: nil},
		{name: "No Phone", request: &sioModel.UpdatePhoneRequest{Number: ""}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Phone too short ", request: &sioModel.UpdatePhoneRequest{Number: "23423"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Phone too long", request: &sioModel.UpdatePhoneRequest{Number: "2934830298420394809238402893"}, status: http.StatusBadRequest, result: nil},
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

			c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

			uc, ms, eu := initController(t)

			err := eu.EncryptInterface(tt.request, false)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "PUT")
			if tt.status == http.StatusOK {
				ms.On("UpdatePhone", "a", mock.AnythingOfType("*sioModel.UpdatePhoneRequest"), c).Return(tt.result)
			}
			uc.UpdatePhone(c)
			if w.Code != tt.status {
				t.Errorf("UpdatePhone() = %v, want %v", w.Code, tt.status)
				return
			}
		})
	}
}

func TestGetUserById(t *testing.T) {
	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)

	uc, ms, _ := initController(t)
	tests := []struct {
		name   string
		want   *sioModel.AwUser
		status int
	}{
		{name: "happy", want: mAwUserPtr, status: 200},
		{name: "service error", want: nil, status: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

			ms.On("GetUserByID", "a", c).Return(tt.want)
			uc.GetUserById(c)

			if w.Code != tt.status {
				t.Errorf("GetUserById()() = %v, want %v", w.Code, tt.status)
				return
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)

	uc, ms, _ := initController(t)
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
			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

			ms.On("DeleteUser", "a", c).Return(tt.want)
			uc.DeleteUser(c)

			if w.Code != tt.status {
				t.Errorf("DeleteUser( )= %v, want %v", w.Code, tt.status)
				return
			}
		})
	}
}
