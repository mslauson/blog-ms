package utils

import (
	sioModel "gitea.slauson.io/slausonio/go-libs/model"
	"gitea.slauson.io/slausonio/go-utils/sioUtils"
)

type IamValidations struct {
	validator *sioUtils.SioValidator
}

func NewIamValidations() *IamValidations {
	return &IamValidations{
		validator: sioUtils.NewValidator(),
	}
}

func (v *IamValidations) ValidateCreateUserRequest(r *sioModel.AwCreateUserRequest) error {
	err := v.validator.ValidateEmail(r.Email)
	if err != nil {
		return err
	}
	err = v.validator.ValidatePassword(r.Password)
	if err != nil {
		return err
	}
	err = v.validator.ValidateName(r.Name)
	if err != nil {
		return err
	}
	err = v.validator.ValidatePhone(r.Phone)

	if err != nil {
		return err
	}
	return nil
}

func (v *IamValidations) ValidateUpdatePasswordRequest(r *sioModel.UpdatePasswordRequest) error {
	err := v.validator.ValidatePassword(r.Password)
	if err != nil {
		return err
	}
	return nil
}

func (v *IamValidations) ValidateUpdateEmailRequest(r *sioModel.UpdateEmailRequest) error {
	err := v.validator.ValidateEmail(r.Email)
	if err != nil {
		return err
	}
	return nil
}

func (v *IamValidations) ValidateUpdatePhoneRequest(r *sioModel.UpdatePhoneRequest) error {
	err := v.validator.ValidatePhone(r.Number)
	if err != nil {
		return err
	}
	return nil
}
