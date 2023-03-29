package service

import (
	"fmt"
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
	mAwUser = sioModel.AwUser{
		Email: "t@t.com",
	}
	mAwUserPtr = &mAwUser
	mUserList  = &sioModel.AwlistResponse{
		Total: 1,
		Users: []sioModel.AwUser{mAwUser},
	}
	mCreateReq = &sioModel.AwCreateUserRequest{
		Email:    "t@t.com",
		Password: "test_password",
		Name:     "test_name",
		Phone:    "test_phone",
	}
	tError       = fmt.Errorf("test error")
	uEmailReq    = &sioModel.UpdateEmailRequest{Email: "test"}
	uPhoneReq    = &sioModel.UpdatePhoneRequest{Number: "1235"}
	uPasswordReq = &sioModel.UpdatePasswordRequest{Password: "1235"}
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
	gc, _ := siotest.InitGinTestContext()

	awClient.On("ListUsers").Return(mUserList, nil)
	actual := us.ListUsers(gc)
	assert.Equalf(t, mUserList, actual, "actual: %v", actual)
	assert.Emptyf(t, gc.Errors, "gc.Errors: %v", gc.Errors)
}

func TestUserService_ListUsers_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, w := siotest.InitGinTestContext()

	awClient.On("ListUsers").Return(nil, tError)
	actual := us.ListUsers(gc)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, gc.Errors[0].Err, iamError.NoCustomersFound, "gc.Errors: %v", gc.Errors)
	assert.Equalf(t, http.StatusNotFound, w.Code, "w.Code: %v", w.Code)
}

func TestUserService_GetUserByID(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, _ := siotest.InitGinTestContext()

	awClient.On("GetUserByID", "a").Return(mAwUserPtr, nil)
	actual := us.GetUserByID("a", gc)
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, gc.Errors, "gc.Errors: %v", gc.Errors)
}

func TestUserService_GetUserByID_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, w := siotest.InitGinTestContext()

	awClient.On("GetUserByID", "a").Return(nil, tError)
	actual := us.GetUserByID("a", gc)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equal(t, gc.Errors[0].Err, iamError.NoCustomerFound)
	assert.Equalf(t, http.StatusNotFound, w.Code, "w.Code: %v", w.Code)
}

func TestUserService_CreateUser(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, _ := siotest.InitGinTestContext()

	awClient.On("CreateUser", mock.AnythingOfType("*sioModel.AwCreateUserRequest")).Return(mAwUserPtr, nil)
	actual := us.CreateUser(mCreateReq, gc)
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, gc.Errors, "gc.Errors: %v", gc.Errors)
}

func TestUserService_CreateUser_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, w := siotest.InitGinTestContext()

	awClient.On("CreateUser", mock.AnythingOfType("*sioModel.AwCreateUserRequest")).Return(nil, tError)
	actual := us.CreateUser(mCreateReq, gc)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, gc.Errors[0].Err, tError, "error: %v", gc.Errors[0].Err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserService_UpdateEmail(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, _ := siotest.InitGinTestContext()

	awClient.On("UpdateEmail", "a", mock.AnythingOfType("*sioModel.UpdateEmailRequest")).Return(mAwUserPtr, nil)
	actual := us.UpdateEmail("a", uEmailReq, gc)
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, gc.Errors, "gc.Errors: %v", gc.Errors)
}

func TestUserService_UpdateEmail_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, w := siotest.InitGinTestContext()

	awClient.On("UpdateEmail", "a", mock.AnythingOfType("*sioModel.UpdateEmailRequest")).Return(nil, tError)
	actual := us.UpdateEmail("a", uEmailReq, gc)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, gc.Errors[0].Err, tError, "error: %v", gc.Errors[0].Err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserService_UpdatePhone(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, _ := siotest.InitGinTestContext()

	awClient.On("UpdatePhone", "a", mock.AnythingOfType("*sioModel.UpdatePhoneRequest")).Return(mAwUserPtr, nil)
	actual := us.UpdatePhone("a", uPhoneReq, gc)
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, gc.Errors, "gc.Errors: %v", gc.Errors)
}

func TestUserService_UpdatePhone_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, w := siotest.InitGinTestContext()

	awClient.On("UpdatePhone", "a", mock.AnythingOfType("*sioModel.UpdatePhoneRequest")).Return(nil, tError)
	actual := us.UpdatePhone("a", uPhoneReq, gc)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, gc.Errors[0].Err, tError, "error: %v", gc.Errors[0].Err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserService_UpdatePassword(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, _ := siotest.InitGinTestContext()

	awClient.On("UpdatePassword", "a", mock.AnythingOfType("*sioModel.UpdatePasswordRequest")).Return(mAwUserPtr, nil)
	actual := us.UpdatePassword("a", uPasswordReq, gc)
	assert.Equalf(t, mAwUserPtr, actual, "actual: %v", actual)
	assert.Emptyf(t, gc.Errors, "gc.Errors: %v", gc.Errors)
}

func TestUserService_UpdatePassword_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, w := siotest.InitGinTestContext()

	awClient.On("UpdatePassword", "a", mock.AnythingOfType("*sioModel.UpdatePasswordRequest")).Return(nil, tError)
	actual := us.UpdatePassword("a", uPasswordReq, gc)
	assert.Nilf(t, actual, "expected nil, actual: %v", actual)
	assert.Equalf(t, gc.Errors[0].Err, tError, "error: %v", gc.Errors[0].Err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUserService_DeleteUser(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, _ := siotest.InitGinTestContext()

	awClient.On("DeleteUser", "a").Return(nil)
	actual := us.DeleteUser("a", gc)
	assert.Truef(t, actual.Success, "actual.Success: %v", actual.Success)
	assert.Emptyf(t, gc.Errors, "gc.Errors: %v", gc.Errors)
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	us, awClient := initUserServiceTest(t)
	gc, w := siotest.InitGinTestContext()

	awClient.On("DeleteUser", "a").Return(tError)
	actual := us.DeleteUser("a", gc)
	assert.False(t, actual.Success)
	assert.Equalf(t, gc.Errors[0].Err, iamError.NoCustomerFound, "error: %v", gc.Errors[0].Err)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
