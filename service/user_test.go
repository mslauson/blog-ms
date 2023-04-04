package service

import (
	"fmt"
	siotest "gitea.slauson.io/slausonio/go-testing/sio_test"
	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
	"gitea.slauson.io/slausonio/iam-ms/client/mocks"
	"gitea.slauson.io/slausonio/iam-ms/constants"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
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
	mCreateReq = &siogeneric.AwCreateUserRequest{
		Email:    "t@t.com",
		Password: "test_password",
		Name:     "test_name",
		Phone:    "test_phone",
	}
	tError       = fmt.Errorf("test error")
	uEmailReq    = &siogeneric.UpdateEmailRequest{Email: "test"}
	uPhoneReq    = &siogeneric.UpdatePhoneRequest{Number: "1235"}
	uPasswordReq = &siogeneric.UpdatePasswordRequest{Password: "1235"}
)

func initUserServiceTest(t *testing.T) (*UserService, *mocks.AppwriteClient) {
	awClient := mocks.NewAppwriteClient(t)
	us := &UserService{
		awClient: awClient,
	}
	return us, awClient
}

func TestUserService_ListUsers(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("ListUsers").Return(mUserList, nil)
	actual, err := us.ListUsers()
	assert.Equalf(t, mUserList, actual, "actual: %v", actual)
	assert.Emptyf(t, err, "error should have been nil. err: %v", err)
}

func TestUserService_ListUsers_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("ListUsers").Return(nil, tError)
	actual, err := us.ListUsers()
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, err, sioerror.NewSioNotFoundError(constants.NoCustomersFound), "actual error: %v", err)
}

func TestUserService_GetUserByID(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("GetUserByID", "a").Return(mAwUserPtr, nil)
	actual, err := us.GetUserByID("a")
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, err, "error should have been nil. err: %v", err)
}

func TestUserService_GetUserByID_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("GetUserByID", "a").Return(nil, tError)
	actual, err := us.GetUserByID("a")
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equal(t, err, sioerror.NewSioNotFoundError(constants.NoCustomerFound))
}

func TestUserService_CreateUser(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("CreateUser", mock.AnythingOfType("*siogeneric.AwCreateUserRequest")).Return(mAwUserPtr, nil)
	actual, err := us.CreateUser(mCreateReq)
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, err, "error should have been nil. err: %v", err)
}

func TestUserService_CreateUser_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("CreateUser", mock.AnythingOfType("*siogeneric.AwCreateUserRequest")).Return(nil, tError)
	actual, err := us.CreateUser(mCreateReq)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, err, siotest.TBadRequestError, "error: %v wanted: %v", err, siotest.TBadRequestError)
}

func TestUserService_UpdateEmail(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("UpdateEmail", "a", mock.AnythingOfType("*siogeneric.UpdateEmailRequest")).Return(mAwUserPtr, nil)
	actual, err := us.UpdateEmail("a", uEmailReq)
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, err, "error should have been nil. err: %v", err)
}

func TestUserService_UpdateEmail_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("UpdateEmail", "a", mock.AnythingOfType("*siogeneric.UpdateEmailRequest")).Return(nil, tError)
	actual, err := us.UpdateEmail("a", uEmailReq)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, err, siotest.TBadRequestError, "error: %v wanted: %v", err, siotest.TBadRequestError)
}

func TestUserService_UpdatePhone(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("UpdatePhone", "a", mock.AnythingOfType("*siogeneric.UpdatePhoneRequest")).Return(mAwUserPtr, nil)
	actual, err := us.UpdatePhone("a", uPhoneReq)
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, err, "error should have been nil. err: %v", err)
}

func TestUserService_UpdatePhone_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("UpdatePhone", "a", mock.AnythingOfType("*siogeneric.UpdatePhoneRequest")).Return(nil, tError)
	actual, err := us.UpdatePhone("a", uPhoneReq)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, err, siotest.TBadRequestError, "error: %v wanted: %v", err, siotest.TBadRequestError)
}

func TestUserService_UpdatePassword(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("UpdatePassword", "a", mock.AnythingOfType("*siogeneric.UpdatePasswordRequest")).Return(mAwUserPtr, nil)
	actual, err := us.UpdatePassword("a", uPasswordReq)
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, err, "error should have been nil. err: %v", err)
}

func TestUserService_UpdatePassword_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("UpdatePassword", "a", mock.AnythingOfType("*siogeneric.UpdatePasswordRequest")).Return(nil, tError)
	actual, err := us.UpdatePassword("a", uPasswordReq)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, err, siotest.TBadRequestError, "error: %v wanted: %v", err, siotest.TBadRequestError)
}

func TestUserService_DeleteUser(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("DeleteUser", "a").Return(nil)
	actual, err := us.DeleteUser("a")
	assert.Truef(t, actual.Success, "actual.Success: %v", actual.Success)
	assert.Emptyf(t, err, "error should have been nil. err: %v", err)
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)

	awClient.On("DeleteUser", "a").Return(tError)
	actual, err := us.DeleteUser("a")
	assert.False(t, actual.Success)
	assert.Equalf(t, err, siotest.TNotFoundError, "error: %v wanted: %v", err, siotest.TBadRequestError)
}
