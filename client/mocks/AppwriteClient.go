// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	mock "github.com/stretchr/testify/mock"
)

// AppwriteClient is an autogenerated mock type for the AppwriteClient type
type AppwriteClient struct {
	mock.Mock
}

// CreateEmailSession provides a mock function with given fields: r
func (_m *AppwriteClient) CreateEmailSession(r *sioModel.AwEmailSessionRequest) (*sioModel.AwSession, error) {
	ret := _m.Called(r)

	var r0 *sioModel.AwSession
	var r1 error
	if rf, ok := ret.Get(0).(func(*sioModel.AwEmailSessionRequest) (*sioModel.AwSession, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*sioModel.AwEmailSessionRequest) *sioModel.AwSession); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sioModel.AwSession)
		}
	}

	if rf, ok := ret.Get(1).(func(*sioModel.AwEmailSessionRequest) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: r
func (_m *AppwriteClient) CreateUser(r *sioModel.AwCreateUserRequest) (*sioModel.AwUser, error) {
	ret := _m.Called(r)

	var r0 *sioModel.AwUser
	var r1 error
	if rf, ok := ret.Get(0).(func(*sioModel.AwCreateUserRequest) (*sioModel.AwUser, error)); ok {
		return rf(r)
	}
	if rf, ok := ret.Get(0).(func(*sioModel.AwCreateUserRequest) *sioModel.AwUser); ok {
		r0 = rf(r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sioModel.AwUser)
		}
	}

	if rf, ok := ret.Get(1).(func(*sioModel.AwCreateUserRequest) error); ok {
		r1 = rf(r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteSession provides a mock function with given fields: sID
func (_m *AppwriteClient) DeleteSession(sID string) error {
	ret := _m.Called(sID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(sID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserByID provides a mock function with given fields: id
func (_m *AppwriteClient) GetUserByID(id string) (*sioModel.AwUser, error) {
	ret := _m.Called(id)

	var r0 *sioModel.AwUser
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*sioModel.AwUser, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *sioModel.AwUser); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sioModel.AwUser)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUsers provides a mock function with given fields:
func (_m *AppwriteClient) ListUsers() (*sioModel.AwlistResponse, error) {
	ret := _m.Called()

	var r0 *sioModel.AwlistResponse
	var r1 error
	if rf, ok := ret.Get(0).(func() (*sioModel.AwlistResponse, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *sioModel.AwlistResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sioModel.AwlistResponse)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateEmail provides a mock function with given fields: id, r
func (_m *AppwriteClient) UpdateEmail(id string, r *sioModel.UpdateEmailRequest) (*sioModel.AwUser, error) {
	ret := _m.Called(id, r)

	var r0 *sioModel.AwUser
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *sioModel.UpdateEmailRequest) (*sioModel.AwUser, error)); ok {
		return rf(id, r)
	}
	if rf, ok := ret.Get(0).(func(string, *sioModel.UpdateEmailRequest) *sioModel.AwUser); ok {
		r0 = rf(id, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sioModel.AwUser)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *sioModel.UpdateEmailRequest) error); ok {
		r1 = rf(id, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePassword provides a mock function with given fields: id, r
func (_m *AppwriteClient) UpdatePassword(id string, r *sioModel.UpdatePasswordRequest) (*sioModel.AwUser, error) {
	ret := _m.Called(id, r)

	var r0 *sioModel.AwUser
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *sioModel.UpdatePasswordRequest) (*sioModel.AwUser, error)); ok {
		return rf(id, r)
	}
	if rf, ok := ret.Get(0).(func(string, *sioModel.UpdatePasswordRequest) *sioModel.AwUser); ok {
		r0 = rf(id, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sioModel.AwUser)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *sioModel.UpdatePasswordRequest) error); ok {
		r1 = rf(id, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePhone provides a mock function with given fields: id, r
func (_m *AppwriteClient) UpdatePhone(id string, r *sioModel.UpdatePhoneRequest) (*sioModel.AwUser, error) {
	ret := _m.Called(id, r)

	var r0 *sioModel.AwUser
	var r1 error
	if rf, ok := ret.Get(0).(func(string, *sioModel.UpdatePhoneRequest) (*sioModel.AwUser, error)); ok {
		return rf(id, r)
	}
	if rf, ok := ret.Get(0).(func(string, *sioModel.UpdatePhoneRequest) *sioModel.AwUser); ok {
		r0 = rf(id, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sioModel.AwUser)
		}
	}

	if rf, ok := ret.Get(1).(func(string, *sioModel.UpdatePhoneRequest) error); ok {
		r1 = rf(id, r)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewAppwriteClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewAppwriteClient creates a new instance of AppwriteClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAppwriteClient(t mockConstructorTestingTNewAppwriteClient) *AppwriteClient {
	mock := &AppwriteClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
