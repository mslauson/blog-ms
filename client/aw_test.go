package client

import (
	"bytes"
	"fmt"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"io/ioutil"
	"net/http"
	"testing"

	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/stretchr/testify/mock"
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
	uEmailR    = &siogeneric.UpdateEmailRequest{Email: "test"}
	uPhoneR    = &siogeneric.UpdatePhoneRequest{Number: "1235"}
	uPasswordR = &siogeneric.UpdatePasswordRequest{Password: "1235"}
)

func TestAwClient_ListUsers(t *testing.T) {
	tests := []struct {
		name  string
		happy bool
		r     error
	}{
		{name: "Happy Path", happy: true, r: nil},
		{name: "Error Path", happy: false, r: fmt.Errorf("test error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			au := new(siogeneric.AwlistResponse)
			h.On("DoHttpRequestAndParse", mock.AnythingOfType("*http.Request"), au).
				Return(tt.r)

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
		name  string
		happy bool
		r     error
	}{
		{name: "Happy Path", happy: true, r: nil},
		{name: "Error Path", happy: false, r: fmt.Errorf("test error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			au := new(siogeneric.AwUser)
			h.On("DoHttpRequestAndParse", mock.AnythingOfType("*http.Request"), au).
				Return(tt.r)

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
		name  string
		happy bool
		r     error
	}{
		{name: "Happy Path", happy: true, r: nil},
		{name: "Error Path", happy: false, r: fmt.Errorf("test error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			au := new(siogeneric.AwUser)
			h.On("DoHttpRequestAndParse", mock.AnythingOfType("*http.Request"), au).
				Return(tt.r)

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
		name  string
		happy bool
		r     error
	}{
		{name: "Happy Path", happy: true, r: nil},
		{name: "Error Path", happy: false, r: fmt.Errorf("test error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			au := new(siogeneric.AwUser)
			h.On("DoHttpRequestAndParse", mock.AnythingOfType("*http.Request"), au).
				Return(tt.r)

			result, err := ac.UpdatePassword("a", uPasswordR)
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
		name  string
		happy bool
		r     error
	}{
		{name: "Happy Path", happy: true, r: nil},
		{name: "Error Path", happy: false, r: fmt.Errorf("test error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			au := new(siogeneric.AwUser)
			h.On("DoHttpRequestAndParse", mock.AnythingOfType("*http.Request"), au).
				Return(tt.r)

			result, err := ac.UpdatePhone("a", uPhoneR)
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
		name  string
		happy bool
		r     error
	}{
		{name: "Happy Path", happy: true, r: nil},
		{name: "Error Path", happy: false, r: fmt.Errorf("test error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			au := new(siogeneric.AwUser)
			h.On("DoHttpRequestAndParse", mock.AnythingOfType("*http.Request"), au).
				Return(tt.r)

			result, err := ac.UpdateEmail("a", uEmailR)
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
		name string
		req  *http.Response
		err  error
	}{
		{name: "Happy Path", req: &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBufferString(`{"status":"ok"}`))}, err: nil},
		{name: "Error Path", req: nil, err: fmt.Errorf("test error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			h.On("DoHttpRequest", mock.AnythingOfType("*http.Request")).
				Return(tt.req, tt.err)

			err := ac.DeleteUser("a")
			if tt.req != nil && err != nil {
				t.Errorf("expected request to resolve but got error %v", err)
				return
			} else if tt.err != nil && tt.req != nil {
				t.Error("expected error to resolve but got result")
				return
			}
		})
	}
}

func TestAwClient_CreateEmailSession(t *testing.T) {
	tests := []struct {
		name  string
		happy bool
		r     error
	}{
		{name: "Happy Path", happy: true, r: nil},
		{name: "Error Path", happy: false, r: fmt.Errorf("test error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			aw := new(siogeneric.AwSession)
			h.On("DoHttpRequestAndParse", mock.AnythingOfType("*http.Request"), aw).
				Return(tt.r)

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
		name string
		req  *http.Response
		err  error
	}{
		{name: "Happy Path", req: &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBufferString(`{"status":"ok"}`))}, err: nil},
		{name: "Error Path", req: nil, err: fmt.Errorf("test error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac, h := initForTests(t)

			h.On("DoHttpRequest", mock.AnythingOfType("*http.Request")).
				Return(tt.req, tt.err)

			err := ac.DeleteSession("a")
			if tt.req != nil && err != nil {
				t.Errorf("expected request to resolve but got error %v", err)
				return
			} else if tt.err != nil && tt.req != nil {
				t.Error("expected error to resolve but got result")
				return
			}
		})
	}
}
