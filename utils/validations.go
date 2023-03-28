package utils

import (
	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
	"github.com/gin-gonic/gin"
)

type IamValidations struct {
	enc       *sioUtils.EncryptionUtil
	validator *sioUtils.SioValidator
}

func NewIamValidations(c *gin.Context) *IamValidations {
	return &IamValidations{
		enc:       sioUtils.NewEncryptionUtil(),
		validator: sioUtils.NewValidator(c),
	}
}

func (v *IamValidations) ValidateCreateUserRequest(r *sioModel.AwCreateUserRequest) bool {
	var (
		emailValid    bool = false
		passwordValid bool = false
		nameValid     bool = false
		phoneValid    bool = false
	)

	emailValid = v.validator.ValidateEmail(r.Email)
	passwordValid = v.validator.ValidatePassword(r.Password)
	nameValid = v.validator.ValidateName(r.Name)
	phoneValid = v.validator.ValidatePhone(r.Phone)

	return emailValid && passwordValid && nameValid && phoneValid
}

func (v *IamValidations) ValidateUpdatePasswordRequest(r *sioModel.UpdatePasswordRequest) bool {
	return v.validator.ValidatePassword(r.Password)
}

func (v *IamValidations) ValidateUpdateEmailRequest(r *sioModel.UpdateEmailRequest) bool {
	return v.validator.ValidateEmail(r.Email)
}

func (v *IamValidations) ValidateUpdatePhoneRequest(r *sioModel.UpdatePhoneRequest) bool {
	return v.validator.ValidatePhone(r.Number)
}

