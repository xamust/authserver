package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/xamust/xvalidator"
)

func NewValidator() xvalidator.XValidator {
	return xvalidator.NewXValidator(passwordTag())
}

// passwordTag just for fun
func passwordTag() xvalidator.InputTagsData {
	return xvalidator.InputTagsData{
		"custom_password",
		"password can't be empty",
		func(fl validator.FieldLevel) bool {
			if fl.Field().String() == "" {
				return false
			}
			return true
		},
	}
}
