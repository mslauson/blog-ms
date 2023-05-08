package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
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

func initController(
	t *testing.T,
) (*UserController, *mocks.IamUserService, *sioUtils.EncryptionUtil) {
	ms := mocks.NewIamUserService(t)
	eu := sioUtils.NewEncryptionUtil()
	controller := &UserController{
		s: ms,
	}

	return controller, ms, eu
}

func MockJson(c *gin.Context, content any, method string) {
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

func TestNewUserController(t *testing.T) {
	uc := NewUserController()
	assert.NotNil(t, uc)
}

func TestListUsers(t *testing.T) {
	uc, ms, _ := initController(t)

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	ms.On("ListUsers").Return(mUserList, nil)
	uc.ListUsers(c)

	assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
}

func TestListUsersError(t *testing.T) {
	uc, ms, _ := initController(t)

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	ms.On("ListUsers").Return(mUserList, errors.New("asdf"))
	uc.ListUsers(c)

	assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
}

func TestUserController_CreateUser(t *testing.T) {
	tests := []struct {
		name      string
		request   *siogeneric.AwCreateUserRequest
		result    *siogeneric.AwUser
		bindError error
	}{
		{
			name: "happy",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "t@t.com",
				Name:     "b",
				Password: "MattTesting&*^1",
			},
			result:    mAwUserPtr,
			bindError: nil,
		},
		{
			name: "Bad Request - No UserID",
			request: &siogeneric.AwCreateUserRequest{
				Phone:    "2121212131",
				Email:    "t@t.com",
				Name:     "b",
				Password: "MattTesting&*^1",
			},
			result:    nil,
			bindError: errors.New("asdf"),
		},
		{
			name: "Bad Request - No Email",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "",
				Name:     "b",
				Password: "MattTesting&*^1",
			},
			result:    nil,
			bindError: errors.New("asdf"),
		},
		{
			name: "Bad Request - bad Email",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "a(",
				Name:     "b",
				Password: "MattTesting&*^1",
			},
			result: nil,
		},
		{
			name: "Bad Request - bad too long",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@.com",
				Name:     "b",
				Password: "MattTesting&*^1",
			},
			result: nil,
		},
		{
			name: "No Password",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "t@t.com",
				Name:     "b",
				Password: "",
			},
			result:    nil,
			bindError: errors.New("asdf"),
		},
		{
			name: "Bad Password Missing Number",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "t@t.com",
				Name:     "b",
				Password: "MattTesting&*^",
			},
			result: nil,
		},
		{
			name: "Bad Password Missing Upper",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "t@t.com",
				Name:     "b",
				Password: "attesting&*^1",
			},
			result: nil,
		},
		{
			name: "Bad Password Missing Special",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "t@t.com",
				Name:     "b",
				Password: "Mattesting1",
			},
			result: nil,
		},
		{
			name: "Short Password",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "t@t.com",
				Name:     "b",
				Password: "Ma",
			},
			result: nil,
		},
		{
			name: "No Name",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "t@t.com",
				Name:     "",
				Password: "MattTesting&*^1",
			},
			result:    nil,
			bindError: errors.New("asdf"),
		},
		{
			name: "long Name",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "2121212131",
				Email:    "t@t.com",
				Name:     "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				Password: "MattTesting&*^1",
			},
			result: nil,
		},
		{
			name: "No Phone",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "",
				Email:    "t@t.com",
				Name:     "b",
				Password: "MattTesting&*^1",
			},
			result:    nil,
			bindError: errors.New("asdf"),
		},
		{
			name: "Bad Phone short",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "131",
				Email:    "t@t.com",
				Name:     "b",
				Password: "MattTesting&*^1",
			},
			result: nil,
		},
		{
			name: "Bad Phone long",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "abc",
				Phone:    "212121123123122131",
				Email:    "t@t.com",
				Name:     "b",
				Password: "MattTesting&*^1",
			},
			result: nil,
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

			uc, ms, eu := initController(t)

			err := eu.EncryptInterface(tt.request)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "POST")

			if tt.result != nil {
				ms.On("CreateUser", mock.AnythingOfType("*siogeneric.AwCreateUserRequest")).
					Return(tt.result, nil)
			}
			uc.CreateUser(c)
			if tt.result != nil {
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestUserController_CreateUserServiceFailure(t *testing.T) {
	request := &siogeneric.AwCreateUserRequest{
		UserID:   "abc",
		Phone:    "2121212131",
		Email:    "t@t.com",
		Name:     "b",
		Password: "MattTesting&*^1",
	}
	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	uc, ms, eu := initController(t)

	err := eu.EncryptInterface(request)
	if err != nil {
		t.Error(err)
		return
	}

	MockJson(c, request, "POST")

	ms.On("CreateUser", mock.AnythingOfType("*siogeneric.AwCreateUserRequest")).
		Return(mAwUserPtr, errors.New("error"))
	uc.CreateUser(c)

	assert.Truef(t, c.Errors != nil, "c.Errors shouldn't be nil")
}

func TestUserController_UpdatePassword(t *testing.T) {
	tests := []struct {
		name      string
		request   *siogeneric.UpdatePasswordRequest
		status    int
		result    *siogeneric.AwUser
		bindError error
	}{
		{
			name:      "happy",
			request:   &siogeneric.UpdatePasswordRequest{Password: "Mm112a23!"},
			status:    http.StatusOK,
			result:    mAwUserPtr,
			bindError: nil,
		},
		{
			name:      "No Password",
			request:   &siogeneric.UpdatePasswordRequest{Password: ""},
			status:    http.StatusBadRequest,
			result:    nil,
			bindError: errors.New("asdf"),
		},
		{
			name:    "Bad Password Missing Number",
			request: &siogeneric.UpdatePasswordRequest{Password: "Mmsdfsdfafsdf!"},
			status:  http.StatusBadRequest,
			result:  nil,
		},
		{
			name:    "Bad Password Missing Upper",
			request: &siogeneric.UpdatePasswordRequest{Password: "aaam112a23!"},
			status:  http.StatusBadRequest,
			result:  nil,
		},
		{
			name:    "Bad Password Missing Special",
			request: &siogeneric.UpdatePasswordRequest{Password: "Mm112a2333"},
			status:  http.StatusBadRequest,
			result:  nil,
		},
		{
			name:    "Short Password",
			request: &siogeneric.UpdatePasswordRequest{Password: "Mm1^"},
			status:  http.StatusBadRequest,
			result:  nil,
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

			c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

			uc, ms, eu := initController(t)

			err := eu.EncryptInterface(tt.request)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "PUT")

			if tt.result != nil {
				ms.On("UpdatePassword", "a", mock.AnythingOfType("*siogeneric.UpdatePasswordRequest")).
					Return(tt.result, nil)
			}
			uc.UpdatePassword(c)
			if tt.result != nil {
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestUserController_UpdatePasswordServiceFailure(t *testing.T) {
	request := &siogeneric.UpdatePasswordRequest{Password: "Mm112a23!"}
	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

	uc, ms, eu := initController(t)

	err := eu.EncryptInterface(request)
	if err != nil {
		t.Error(err)
		return
	}

	MockJson(c, request, "PUT")

	ms.On("UpdatePassword", "a", mock.AnythingOfType("*siogeneric.UpdatePasswordRequest")).
		Return(mAwUserPtr, errors.New("error"))
	uc.UpdatePassword(c)

	assert.Truef(t, c.Errors != nil, "c.Errors shouldn't be nil")
}

func TestUserController_UpdateEmail(t *testing.T) {
	tests := []struct {
		name      string
		request   *siogeneric.UpdateEmailRequest
		result    *siogeneric.AwUser
		bindError error
	}{
		{
			name:    "happy",
			request: &siogeneric.UpdateEmailRequest{Email: "fake@fake.com"},
			result:  mAwUserPtr,
		},
		{
			name:      "No Email",
			request:   &siogeneric.UpdateEmailRequest{Email: ""},
			result:    nil,
			bindError: errors.New("asdf"),
		},
		{
			name:    "Bad Email ",
			request: &siogeneric.UpdateEmailRequest{Email: "fakefake.com"},
			result:  nil,
		},
		{
			name: "Bad Email too long",
			request: &siogeneric.UpdateEmailRequest{
				Email: "\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@.com",
			},
			result: nil,
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

			c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

			uc, ms, eu := initController(t)

			err := eu.EncryptInterface(tt.request)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "PUT")

			if tt.result != nil {
				ms.On("UpdateEmail", "a", mock.AnythingOfType("*siogeneric.UpdateEmailRequest")).
					Return(tt.result, nil)
			}
			uc.UpdateEmail(c)
			if tt.result != nil {
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestUserController_UpdateEmailServiceFailure(t *testing.T) {
	request := &siogeneric.UpdateEmailRequest{Email: "a@abc.com"}

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

	uc, ms, eu := initController(t)

	err := eu.EncryptInterface(request)
	if err != nil {
		t.Error(err)
		return
	}

	MockJson(c, request, "PUT")

	ms.On("UpdateEmail", "a", mock.AnythingOfType("*siogeneric.UpdateEmailRequest")).
		Return(mAwUserPtr, errors.New("error"))
	uc.UpdateEmail(c)

	assert.Truef(t, c.Errors != nil, "c.Errors shouldn't be nil")
}

func TestUserController_UpdatePhone(t *testing.T) {
	tests := []struct {
		name    string
		request *siogeneric.UpdatePhoneRequest
		result  *siogeneric.AwUser
		bindErr error
	}{
		{
			name:    "happy",
			request: &siogeneric.UpdatePhoneRequest{Number: "1239323939"},
			result:  mAwUserPtr,
		},
		{
			name:    "No Phone",
			request: &siogeneric.UpdatePhoneRequest{Number: ""},
			result:  nil,
			bindErr: errors.New("asdf"),
		},
		{
			name:    "Bad Phone too short ",
			request: &siogeneric.UpdatePhoneRequest{Number: "23423"},
			result:  nil,
		},
		{
			name:    "Bad Phone too long",
			request: &siogeneric.UpdatePhoneRequest{Number: "2934830298420394809238402893"},
			result:  nil,
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

			c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

			uc, ms, eu := initController(t)

			err := eu.EncryptInterface(tt.request)
			if err != nil {
				t.Error(err)
				return
			}

			MockJson(c, tt.request, "PUT")
			if tt.result != nil {
				ms.On("UpdatePhone", "a", mock.AnythingOfType("*siogeneric.UpdatePhoneRequest")).
					Return(tt.result, nil)
			}

			uc.UpdatePhone(c)
			if tt.result != nil {
				assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
			} else {
				assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
			}
		})
	}
}

func TestUserController_UpdatePhoneServiceFailure(t *testing.T) {
	request := &siogeneric.UpdatePhoneRequest{Number: "3647586976"}

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}

	uc, ms, eu := initController(t)

	err := eu.EncryptInterface(request)
	if err != nil {
		t.Error(err)
		return
	}

	MockJson(c, request, "PUT")

	ms.On("UpdatePhone", "a", mock.AnythingOfType("*siogeneric.UpdatePhoneRequest")).
		Return(mAwUserPtr, errors.New("error"))
	uc.UpdatePhone(c)

	assert.Truef(t, c.Errors != nil, "c.Errors shouldn't be nil")
}

func TestGetUserById(t *testing.T) {
	uc, ms, _ := initController(t)

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	ms.On("GetUserByID", "a").Return(mAwUserPtr, nil)
	uc.GetUserById(c)

	assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
}

func TestGetUserByIdError(t *testing.T) {
	uc, ms, _ := initController(t)

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	ms.On("GetUserByID", "a").Return(nil, errors.New("asdf"))
	uc.GetUserById(c)

	assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
}

func TestDeleteUser(t *testing.T) {
	uc, ms, _ := initController(t)

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	ms.On("DeleteUser", "a").Return(siogeneric.SuccessResponse{Success: true}, nil)
	uc.DeleteUser(c)

	assert.Truef(t, c.Errors == nil, "c.Errors should be nil")
}

func TestDeleteUserError(t *testing.T) {
	uc, ms, _ := initController(t)

	var (
		w    = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(w)
	)
	c.Request = &http.Request{
		Header: make(http.Header),
	}

	c.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	ms.On("DeleteUser", "a").Return(siogeneric.SuccessResponse{Success: false}, errors.New("asdf"))
	uc.DeleteUser(c)

	assert.Truef(t, c.Errors != nil, "c.Errors shouldnt be nil")
}
