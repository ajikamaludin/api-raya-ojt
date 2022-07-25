package validator

import "github.com/go-playground/validator/v10"

type Validator struct {
}

func New() *Validator {
	return &Validator{}
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func (v *Validator) ValidateRequest(model interface{}) []*ErrorResponse {
	var validate = validator.New()
	var errors []*ErrorResponse
	err := validate.Struct(model)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
