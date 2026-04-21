package helpers

import "github.com/go-playground/validator/v10"

type errValidation struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

func GetErrorValidations(err error) []errValidation {
	groupErrs := []errValidation{}
	for _, err := range err.(validator.ValidationErrors) {
		groupErrs = append(groupErrs, errValidation{
			Field: err.Field(),
			Err:   err.Error(),
		})
	}
	return groupErrs
}
