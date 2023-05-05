package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
)

func initForTests(t *testing.T) (*AwClient, *sioUtils.MockSioRestHelpers) {
	h := sioUtils.NewMockSioRestHelpers(t)
	ac := &AwClient{
		h: h,
		defaultHeaders: map[string][]string{
			"Content-Type":       {"application/json"},
			"X-Appwrite-Project": {"fake"},
		},
		host: "http://localhost:8080/v1",
		key:  "test",
	}
	return ac, h
}

var (
	mAwUser = siogeneric.AwUser{
		Email: "t@t.com",
	}
	mUserList = &siogeneric.AwlistResponse{
		Total: 1,
		Users: []siogeneric.AwUser{mAwUser},
	}

	sessionReq = &siogeneric.AwEmailSessionRequest{
		Email:    "test",
		Password: "test",
	}
	mCr = &siogeneric.AwCreateUserRequest{
		Email:    "t@t.com",
		Password: "test_password",
		Name:     "test_name",
		Phone:    "test_phone",
	}
)

func TestAwClient_ListUsers(t *testing.T) {
	tests := []struct {
		name     string
		happy    bool
		result   *siogeneric.AwlistResponse
		execErr  error
		ParseErr error
		code     int
	}{
		{
			name:     "Happy Path",
			happy:    true,
			result:   mUserList,
			execErr:  nil,
			ParseErr: nil,
			code:     http.StatusOK,
		},
		{
			name:     "Exec Error",
			happy:    false,
			result:   nil,
			execErr:  fmt.Errorf("test error"),
			ParseErr: nil,
			code:     http.StatusInternalServerError,
		},
		{
			name:     "Parse Error",
			happy:    false,
			result:   nil,
			execErr:  nil,
			ParseErr: fmt.Errorf("test error"),
			code:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			mockRes := mockHttpResponse(t, mAwUser, tt.code)
			h.On("ExecuteRequest", mock.AnythingOfType("*http.Request")).
				Return(mockRes, tt.execErr)

			if tt.execErr == nil {
				h.On("ParseResponse", mock.AnythingOfType("*http.Response"), mock.AnythingOfType("*siogeneric.AwlistResponse")).
					Return(tt.ParseErr)
			}

			result, err := ac.ListUsers()
			if tt.happy && result == nil {
				t.Errorf("expected result but got nil")
				return
			} else if err == nil && !tt.happy {
				if tt.happy && err != nil {
					t.Errorf("error during create() error = %v", err)
					return
				} else if err == nil && !tt.happy {
					t.Errorf("expected error but got nil")
					return
				}
			}
		})
	}
}

func TestAwClient_GetUserByID(t *testing.T) {
	tests := []struct {
		name     string
		happy    bool
		execErr  error
		parseErr error
		code     int
	}{
		{name: "Happy Path", happy: true, execErr: nil, parseErr: nil, code: http.StatusOK},
		{
			name:     "Exec Error",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusInternalServerError,
		},
		{
			name:     "Parse Error",
			happy:    false,
			execErr:  nil,
			parseErr: fmt.Errorf("test error"),
			code:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			mockRes := mockHttpResponse(t, mAwUser, tt.code)
			h.On("ExecuteRequest", mock.AnythingOfType("*http.Request")).
				Return(mockRes, tt.execErr)

			if tt.execErr == nil {
				h.On("ParseResponse", mock.AnythingOfType("*http.Response"), mock.AnythingOfType("*siogeneric.AwUser")).
					Return(tt.parseErr)
			}

			result, err := ac.GetUserByID("test")
			if tt.happy && result == nil {
				t.Errorf("expected result but got nil")
				return
			} else if err == nil && !tt.happy {
				if tt.happy && err != nil {
					t.Errorf("error during create() error = %v", err)
					return
				} else if err == nil && !tt.happy {
					t.Errorf("expected error but got nil")
					return
				}
			}
		})
	}
}

func TestAwClient_CreateUser(t *testing.T) {
	tests := []struct {
		name     string
		happy    bool
		execErr  error
		parseErr error
		code     int
	}{
		{name: "Happy Path", happy: true, execErr: nil, parseErr: nil, code: http.StatusOK},
		{
			name:     "ExecErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusInternalServerError,
		},
		{
			name:     "ParseErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			mockRes := mockHttpResponse(t, mAwUser, tt.code)
			h.On("ExecuteRequest", mock.AnythingOfType("*http.Request")).
				Return(mockRes, tt.execErr)

			if tt.execErr == nil {
				h.On("ParseResponse", mock.AnythingOfType("*http.Response"), mock.AnythingOfType("*siogeneric.AwUser")).
					Return(tt.parseErr)
			}

			result, err := ac.CreateUser(mCr)
			if tt.happy && result == nil {
				t.Errorf("expected result but got nil")
				return
			} else if err == nil && !tt.happy {
				if tt.happy && err != nil {
					t.Errorf("error during create() error = %v", err)
					return
				} else if err == nil && !tt.happy {
					t.Errorf("expected error but got nil")
					return
				}
			}
		})
	}
}

func TestAwClient_UpdatePassword(t *testing.T) {
	tests := []struct {
		name     string
		happy    bool
		execErr  error
		parseErr error
		code     int
	}{
		{name: "Happy Path", happy: true, execErr: nil, parseErr: nil, code: http.StatusOK},
		{
			name:     "ExecErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusInternalServerError,
		},
		{
			name:     "ParseErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			mockRes := mockHttpResponse(t, mAwUser, tt.code)
			h.On("ExecuteRequest", mock.AnythingOfType("*http.Request")).
				Return(mockRes, tt.execErr)

			if tt.execErr == nil {
				h.On("ParseResponse", mock.AnythingOfType("*http.Response"), mock.AnythingOfType("*siogeneric.AwUser")).
					Return(tt.parseErr)
			}

			result, err := ac.CreateUser(mCr)
			if tt.happy && result == nil {
				t.Errorf("expected result but got nil")
				return
			} else if err == nil && !tt.happy {
				if tt.happy && err != nil {
					t.Errorf("error during create() error = %v", err)
					return
				} else if err == nil && !tt.happy {
					t.Errorf("expected error but got nil")
					return
				}
			}
		})
	}
}

func TestAwClient_UpdatePhone(t *testing.T) {
	tests := []struct {
		name     string
		happy    bool
		execErr  error
		parseErr error
		code     int
	}{
		{name: "Happy Path", happy: true, execErr: nil, parseErr: nil, code: http.StatusOK},
		{
			name:     "ExecErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusInternalServerError,
		},
		{
			name:     "ParseErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			mockRes := mockHttpResponse(t, mAwUser, tt.code)
			h.On("ExecuteRequest", mock.AnythingOfType("*http.Request")).
				Return(mockRes, tt.execErr)

			if tt.execErr == nil {
				h.On("ParseResponse", mock.AnythingOfType("*http.Response"), mock.AnythingOfType("*siogeneric.AwUser")).
					Return(tt.parseErr)
			}
			result, err := ac.CreateUser(mCr)
			if tt.happy && result == nil {
				t.Errorf("expected result but got nil")
				return
			} else if err == nil && !tt.happy {
				if tt.happy && err != nil {
					t.Errorf("error during create() error = %v", err)
					return
				} else if err == nil && !tt.happy {
					t.Errorf("expected error but got nil")
					return
				}
			}
		})
	}
}

func TestAwClient_UpdateEmail(t *testing.T) {
	tests := []struct {
		name     string
		happy    bool
		execErr  error
		parseErr error
		code     int
	}{
		{name: "Happy Path", happy: true, execErr: nil, parseErr: nil, code: http.StatusOK},
		{
			name:     "ExecErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusInternalServerError,
		},
		{
			name:     "ParseErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			mockRes := mockHttpResponse(t, mAwUser, tt.code)
			h.On("ExecuteRequest", mock.AnythingOfType("*http.Request")).
				Return(mockRes, tt.execErr)

			if tt.execErr == nil {
				h.On("ParseResponse", mock.AnythingOfType("*http.Response"), mock.AnythingOfType("*siogeneric.AwUser")).
					Return(tt.parseErr)
			}
			result, err := ac.CreateUser(mCr)
			if tt.happy && result == nil {
				t.Errorf("expected result but got nil")
				return
			} else if err == nil && !tt.happy {
				if tt.happy && err != nil {
					t.Errorf("error during create() error = %v", err)
					return
				} else if err == nil && !tt.happy {
					t.Errorf("expected error but got nil")
					return
				}
			}
		})
	}
}

func TestAwClient_DeleteUser(t *testing.T) {
	tests := []struct {
		name    string
		result  *http.Response
		execErr error
		code    int
		happy   bool
	}{
		{
			name: "Happy Path",
			result: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"status":"ok"}`)),
			},
			execErr: nil,
			code:    http.StatusOK,
			happy:   true,
		},
		{
			name:    "ExecErr",
			result:  nil,
			execErr: fmt.Errorf("test error"),
			code:    http.StatusInternalServerError,
			happy:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			mockRes := mockHttpResponse(t, mAwUser, tt.code)
			h.On("ExecuteRequest", mock.AnythingOfType("*http.Request")).
				Return(mockRes, tt.execErr)

			err := ac.DeleteUser("a")
			if tt.result != nil && err != nil {
				t.Errorf("expected request to resolve but got error %v", err)
				return
			} else if tt.result == nil && err == nil {
				t.Error("expected error to resolve but got result")
				return
			}
		})
	}
}

func TestAwClient_CreateEmailSession(t *testing.T) {
	tests := []struct {
		name     string
		happy    bool
		execErr  error
		parseErr error
		code     int
	}{
		{name: "Happy Path", happy: true, execErr: nil, parseErr: nil, code: http.StatusOK},
		{
			name:     "ExecErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusInternalServerError,
		},
		{
			name:     "ParseErr",
			happy:    false,
			execErr:  fmt.Errorf("test error"),
			parseErr: nil,
			code:     http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			mockRes := mockHttpResponse(t, mAwUser, tt.code)
			h.On("ExecuteRequest", mock.AnythingOfType("*http.Request")).
				Return(mockRes, tt.execErr)

			if tt.execErr == nil {
				h.On("ParseResponse", mock.AnythingOfType("*http.Response"), mock.AnythingOfType("*siogeneric.AwSession")).
					Return(tt.parseErr)
			}
			result, err := ac.CreateEmailSession(sessionReq)
			if tt.happy && result == nil {
				t.Errorf("expected result but got nil")
				return
			} else if err == nil && !tt.happy {
				if tt.happy && err != nil {
					t.Errorf("error during create() error = %v", err)
					return
				} else if err == nil && !tt.happy {
					t.Errorf("expected error but got nil")
					return
				}
			}
		})
	}
}

func TestAwClient_DeleteUserSession(t *testing.T) {
	tests := []struct {
		name    string
		result  *http.Response
		execErr error
		code    int
	}{
		{
			name: "Happy Path",
			result: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"status":"ok"}`)),
			},
			execErr: nil,
			code:    http.StatusOK,
		},
		{
			name:    "ExecErr",
			result:  nil,
			execErr: fmt.Errorf("test error"),
			code:    http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			mockRes := mockHttpResponse(t, mAwUser, tt.code)
			h.On("ExecuteRequest", mock.AnythingOfType("*http.Request")).
				Return(mockRes, tt.execErr)

			err := ac.DeleteSession("1", "a")
			if tt.result != nil && err != nil {
				t.Errorf("expected request to resolve but got error %v", err)
				return
			} else if tt.result == nil && err == nil {
				t.Error("expected error to resolve but got result")
				return
			}
		})
	}
}

func mockHttpResponse(t *testing.T, v any, code int) *http.Response {
	jsonData, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Error marshaling struct: %v\n", err)
	}

	response := &http.Response{
		StatusCode: code,
		Body:       io.NopCloser(bytes.NewBuffer(jsonData)),
	}
	return response
}
