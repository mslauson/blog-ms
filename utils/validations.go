package utils

import (
	"gitea.slauson.io/slausonio/go-types/siogeneric"
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

func (v *IamValidations) ValidateCreateUserRequest(r *siogeneric.AwCreateUserRequest) error {
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

func (v *IamValidations) ValidateUpdatePasswordRequest(r *siogeneric.UpdatePasswordRequest) error {
	err := v.validator.ValidatePassword(r.Password)
	if err != nil {
		return err
	}
	return nil
}

func (v *IamValidations) ValidateUpdateEmailRequest(r *siogeneric.UpdateEmailRequest) error {
	err := v.validator.ValidateEmail(r.Email)
	if err != nil {
		return err
	}
	return nil
}

func (v *IamValidations) ValidateUpdatePhoneRequest(r *siogeneric.UpdatePhoneRequest) error {
	err := v.validator.ValidatePhone(r.Number)
	if err != nil {
		return err
	}
	return nil
}
