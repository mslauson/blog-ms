package controller

import (
	"bytes"
	"encoding/json"
	"gitea.slauson.io/slausonio/go-types/siogeneric"

	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"gitea.slauson.io/slausonio/iam-ms/service/mocks"
)

var (
	mAwUser = siogeneric.AwUser{
		Email: "t@t.com",
	}
	mAwUserPtr = &mAwUser
	mUserList  = &siogeneric.AwlistResponse{
		Total: 1,
		Users: []siogeneric.AwUser{mAwUser},
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
		want   *siogeneric.AwlistResponse
		status int
	}{
		{name: "happy", want: mUserList, status: 200},
		{name: "service error", want: mUserList, status: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.On("ListUsers").Return(tt.want)
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
		request *siogeneric.AwCreateUserRequest
		status  int
		result  *siogeneric.AwUser
	}{
		{name: "happy", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusOK, result: mAwUserPtr},
		{name: "service failure", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusOK, result: nil},
		{name: "Bad Request - No Email", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Request - bad Email", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "a(", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Request - bad too long", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "No Password", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: ""}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Number", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Upper", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "attesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Special", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "Mattesting1"}, status: http.StatusBadRequest, result: nil},
		{name: "Short Password", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "b", Password: "Ma"}, status: http.StatusBadRequest, result: nil},
		{name: "No Name", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "long Name", request: &siogeneric.AwCreateUserRequest{Phone: "2121212131", Email: "t@t.com", Name: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "No Phone", request: &siogeneric.AwCreateUserRequest{Phone: "", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Phone short", request: &siogeneric.AwCreateUserRequest{Phone: "131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Phone long", request: &siogeneric.AwCreateUserRequest{Phone: "212121123123122131", Email: "t@t.com", Name: "b", Password: "MattTesting&*^1"}, status: http.StatusBadRequest, result: nil},
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
				ms.On("CreateUser", mock.AnythingOfType("*siogeneric.AwCreateUserRequest")).Return(tt.result)
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
		request *siogeneric.UpdatePasswordRequest
		status  int
		result  *siogeneric.AwUser
	}{
		{name: "happy", request: &siogeneric.UpdatePasswordRequest{Password: "Mm112a23!"}, status: http.StatusOK, result: mAwUserPtr},
		{name: "service failure", request: &siogeneric.UpdatePasswordRequest{Password: "Mm112a23!"}, status: http.StatusOK, result: nil},
		{name: "No Password", request: &siogeneric.UpdatePasswordRequest{Password: ""}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Number", request: &siogeneric.UpdatePasswordRequest{Password: "Mmsdfsdfafsdf!"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Upper", request: &siogeneric.UpdatePasswordRequest{Password: "aaam112a23!"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Password Missing Special", request: &siogeneric.UpdatePasswordRequest{Password: "Mm112a2333"}, status: http.StatusBadRequest, result: nil},
		{name: "Short Password", request: &siogeneric.UpdatePasswordRequest{Password: "Mm1^"}, status: http.StatusBadRequest, result: nil},
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
				ms.On("UpdatePassword", "a", mock.AnythingOfType("*siogeneric.UpdatePasswordRequest")).Return(tt.result)
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
		request *siogeneric.UpdateEmailRequest
		status  int
		result  *siogeneric.AwUser
	}{
		{name: "happy", request: &siogeneric.UpdateEmailRequest{Email: "fake@fake.com"}, status: http.StatusOK, result: mAwUserPtr},
		{name: "service failure", request: &siogeneric.UpdateEmailRequest{Email: "fake@fake.com"}, status: http.StatusOK, result: nil},
		{name: "No Email", request: &siogeneric.UpdateEmailRequest{Email: ""}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Email ", request: &siogeneric.UpdateEmailRequest{Email: "fakefake.com"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Email too long", request: &siogeneric.UpdateEmailRequest{Email: "\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@.com"}, status: http.StatusBadRequest, result: nil},
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
				ms.On("UpdateEmail", "a", mock.AnythingOfType("*siogeneric.UpdateEmailRequest")).Return(tt.result)
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
		request *siogeneric.UpdatePhoneRequest
		status  int
		result  *siogeneric.AwUser
	}{
		{name: "happy", request: &siogeneric.UpdatePhoneRequest{Number: "1239323939"}, status: http.StatusOK, result: mAwUserPtr},
		{name: "service failure", request: &siogeneric.UpdatePhoneRequest{Number: "1239323939"}, status: http.StatusOK, result: nil},
		{name: "No Phone", request: &siogeneric.UpdatePhoneRequest{Number: ""}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Phone too short ", request: &siogeneric.UpdatePhoneRequest{Number: "23423"}, status: http.StatusBadRequest, result: nil},
		{name: "Bad Phone too long", request: &siogeneric.UpdatePhoneRequest{Number: "2934830298420394809238402893"}, status: http.StatusBadRequest, result: nil},
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
				ms.On("UpdatePhone", "a", mock.AnythingOfType("*siogeneric.UpdatePhoneRequest")).Return(tt.result)
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
		want   *siogeneric.AwUser
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

			ms.On("GetUserByID", "a").Return(tt.want)
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
		want   siogeneric.SuccessResponse
		status int
	}{
		{name: "happy", want: siogeneric.SuccessResponse{Success: true}, status: 200},
		{name: "service error", want: siogeneric.SuccessResponse{Success: false}, status: 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.Request = &http.Request{
				Header: make(http.Header),
			}

			c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

			ms.On("DeleteUser", "a").Return(tt.want)
			uc.DeleteUser(c)

			if w.Code != tt.status {
				t.Errorf("DeleteUser( )= %v, want %v", w.Code, tt.status)
				return
			}
		})
	}
}
