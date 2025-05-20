package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationErrorsResponse struct {
	Response
	ValidationErrors map[string]string `json:"errors"`
}

func ValidationErrorsResp(
	errs validator.ValidationErrors,
) *ValidationErrorsResponse {
	fields := make(map[string]string, len(errs))
	for _, err := range errs {
		fields[err.Field()] =
			fmt.Sprintf("validation by tag %s failed", err.Tag())
	}

	return &ValidationErrorsResponse{
		Response:         Error("validation failed"),
		ValidationErrors: fields,
	}
}
