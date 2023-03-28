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

func (v *IamValidations) ValidateCreateUserRequest(r *sioModel.AwCreateUserRequest) {
	var (
		emailValid    bool = false
		passwordValid bool = false
		nameValid     bool = false
		phoneValid    bool = false
	)

	emailValid = v.validator.ValidateEmail(r.Email)
}
