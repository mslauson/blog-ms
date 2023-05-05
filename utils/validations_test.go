package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitea.slauson.io/slausonio/go-types/siogeneric"
	"gitea.slauson.io/slausonio/go-utils/sioerror"
)

func TestValidateCreateUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *siogeneric.AwCreateUserRequest
		error   error
	}{
		{
			name: "valid request",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "10000069",
				Phone:    "5555555555",
				Email:    "fake@fake.com",
				Name:     "Fakey McFakerson",
				Password: "Fake@123",
			},

			error: nil,
		},

		{
			name: "bad phone short",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "10000069",
				Phone:    "5555555",
				Email:    "fake@fake.com",
				Name:     "Fakey McFakerson",
				Password: "Fake@123",
			},
			error: sioerror.NewSioBadRequestError("please enter a ten digit mobile number"),
		},
		{
			name: "bad phone long",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "10000069",
				Phone:    "5555889555555",
				Email:    "fake@fake.com",
				Name:     "Fakey McFakerson",
				Password: "Fake@123",
			},
			error: sioerror.NewSioBadRequestError("please enter a ten digit mobile number"),
		},
		{
			name: "bad email",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "10000069",
				Phone:    "5555555555",
				Email:    "fakefake.com",
				Name:     "Fakey McFakerson",
				Password: "Fake@123",
			},
			error: sioerror.NewSioBadRequestError("invalid email"),
		},
		{
			name: "name too long",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "10000069",
				Phone:    "5555555555",
				Email:    "fake@fake.com",
				Name:     "Fakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfFakey McFakersonasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdf",
				Password: "Fake@123",
			},
			error: sioerror.NewSioBadRequestError("invalid name"),
		},
		{
			name: "Bad Password Missing Number",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "10000069",
				Phone:    "5555555555",
				Email:    "fake@fake.com",
				Name:     "Fakey McFakerson",
				Password: "Faksdfsdfe@",
			},
			error: sioerror.NewSioBadRequestError(
				"invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
			),
		},
		{
			name: "Bad Password Missing Upper",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "10000069",
				Phone:    "5555555555",
				Email:    "fake@fake.com",
				Name:     "Fakey McFakerson",
				Password: "fake@123",
			},
			error: sioerror.NewSioBadRequestError(
				"invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
			),
		},
		{
			name: "Bad Password Missing Special",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "10000069",
				Phone:    "5555555555",
				Email:    "fake@fake.com",
				Name:     "Fakey McFakerson",
				Password: "Fake123",
			},
			error: sioerror.NewSioBadRequestError(
				"invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
			),
		},
		{
			name: "Short Password",
			request: &siogeneric.AwCreateUserRequest{
				UserID:   "10000069",
				Phone:    "5555555555",
				Email:    "fake@fake.com",
				Name:     "Fakey McFakerson",
				Password: "F@123",
			},
			error: sioerror.NewSioBadRequestError(
				"invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
			),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewIamValidations()
			err := v.ValidateCreateUserRequest(test.request)
			if test.error == nil {
				assert.Nilf(t, err, "Expected no error, got %v", err)
			} else {
				assert.Equalf(
					t,
					test.error.Error(),
					err.Error(),
					"Expected error %s, got %s",
					test.error.Error(),
					err.Error(),
				)
			}
		})
	}
}

func TestValidateUpdateEmailRequest(t *testing.T) {
	tests := []struct {
		name    string
		request *siogeneric.UpdateEmailRequest
		error   error
	}{
		{
			name: "Valid",
			request: &siogeneric.UpdateEmailRequest{
				Email: "fake.fake@com",
			},
			error: nil,
		},
		{
			name: "Missing @",
			request: &siogeneric.UpdateEmailRequest{
				Email: "fakefake.com",
			},
			error: sioerror.NewSioBadRequestError("invalid email"),
		},
		{
			name: "Missing Domain",
			request: &siogeneric.UpdateEmailRequest{
				Email: "fake@fake.",
			},
			error: sioerror.NewSioBadRequestError("invalid email"),
		},
		{
			name: "Missing Name",
			request: &siogeneric.UpdateEmailRequest{
				Email: "fakefake.com",
			},
			error: sioerror.NewSioBadRequestError("invalid email"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewIamValidations()
			err := v.ValidateUpdateEmailRequest(test.request)
			if test.error == nil {
				assert.Nilf(t, err, "Expected no error, got %v", err)
			} else {
				assert.Equalf(
					t,
					test.error.Error(),
					err.Error(),
					"Expected error %s, got %s",
					test.error.Error(),
					err.Error(),
				)
			}
		})
	}
}

func TestValidateUpdatePassword(t *testing.T) {
	tests := []struct {
		name    string
		request *siogeneric.UpdatePasswordRequest
		error   error
	}{
		{
			name: "Valid",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "Faksdfsdfe@123!",
			},
			error: nil,
		},
		{
			name: "Bad Password Missing Number",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "Faksdfsdfe@",
			},
			error: sioerror.NewSioBadRequestError(
				"invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
			),
		},
		{
			name: "Bad Password Missing Upper",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "fake@123",
			},
			error: sioerror.NewSioBadRequestError(
				"invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
			),
		},
		{
			name: "Bad Password Missing Special",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "Fake123",
			},
			error: sioerror.NewSioBadRequestError(
				"invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
			),
		},
		{
			name: "Short Password",
			request: &siogeneric.UpdatePasswordRequest{
				Password: "F@123",
			},
			error: sioerror.NewSioBadRequestError(
				"invalid password: Requirements are 8 char min, 1 upper, 1 special, and 1 numerical",
			),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewIamValidations()
			err := v.ValidateUpdatePasswordRequest(test.request)
			if test.error == nil {
				assert.Nilf(t, err, "Expected no error, got %v", err)
			} else {
				assert.Equalf(
					t,
					test.error.Error(),
					err.Error(),
					"Expected error %s, got %s",
					test.error.Error(),
					err.Error(),
				)
			}
		})
	}
}

func TestValidateUpdatePhone(t *testing.T) {
	tests := []struct {
		name    string
		request *siogeneric.UpdatePhoneRequest
		error   error
	}{
		{
			name: "Valid",
			request: &siogeneric.UpdatePhoneRequest{
				Number: "5555555555",
			},
			error: nil,
		},
		{
			name: "bad phone short",
			request: &siogeneric.UpdatePhoneRequest{
				Number: "5555555",
			},
			error: sioerror.NewSioBadRequestError("please enter a ten digit mobile number"),
		},
		{
			name: "bad phone long",
			request: &siogeneric.UpdatePhoneRequest{
				Number: "5555889555555",
			},
			error: sioerror.NewSioBadRequestError("please enter a ten digit mobile number"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := NewIamValidations()
			err := v.ValidateUpdatePhoneRequest(test.request)
			if test.error == nil {
				assert.Nilf(t, err, "Expected no error, got %v", err)
			} else {
				assert.Equalf(
					t,
					test.error.Error(),
					err.Error(),
					"Expected error %s, got %s",
					test.error.Error(),
					err.Error(),
				)
			}
		})
	}
}
