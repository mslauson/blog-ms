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
	if err := v.validator.ValidateEmail(r.Email); err != nil {
		return err
	}

	if err := v.validator.ValidatePassword(r.Password); err != nil {
		return err
	}

	if err := v.validator.ValidateName(r.Name); err != nil {
		return err
	}

	if err := v.validator.ValidatePhone(r.Phone); err != nil {
		return err
	}

	return nil
}

func (v *IamValidations) ValidateUpdatePasswordRequest(r *siogeneric.UpdatePasswordRequest) error {
	if err := v.validator.ValidatePassword(r.Password); err != nil {
		return err
	}

	return nil
}

func (v *IamValidations) ValidateUpdateEmailRequest(r *siogeneric.UpdateEmailRequest) error {
	if err := v.validator.ValidateEmail(r.Email); err != nil {
		return err
	}

	return nil
}

func (v *IamValidations) ValidateUpdatePhoneRequest(r *siogeneric.UpdatePhoneRequest) error {
	if err := v.validator.ValidatePhone(r.Number); err != nil {
		return err
	}

	return nil
}
