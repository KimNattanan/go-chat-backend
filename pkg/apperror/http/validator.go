package http

import "github.com/go-playground/validator/v10"

func ParseValidationErrors(errs validator.ValidationErrors) map[string]string {
	out := make(map[string]string)

	for _, e := range errs {
		out[e.Field()] = validationMessage(e)
	}

	return out
}

func validationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "invalid email format"
	case "min":
		return "value is below minimum"
	case "max":
		return "value exceeds maximum"
	case "len":
		return "invalid length"
	case "oneof":
		return "invalid option"
	case "gt":
		return "must be greater than allowed value"
	case "gte":
		return "must be greater than or equal"
	case "lt":
		return "must be less than allowed value"
	case "lte":
		return "must be less than or equal"
	default:
		return "invalid value"
	}
}
