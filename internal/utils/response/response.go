package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type GeneralError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-missing-field", "foreign_key") //  randomly adding header just because I can
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func BaseError(err error) GeneralError {
	return GeneralError{
		Status: "Error",
		Error:  err.Error(),
	}
}

func ValidationErrors(errs validator.ValidationErrors) GeneralError {
	var validationErrors []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			validationErrors = append(validationErrors, fmt.Sprintf("%s field is required", err.Field()))
		default:
			validationErrors = append(validationErrors, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}

	return GeneralError{
		Status: "Validation Errors",
		Error:  strings.Join(validationErrors, ", "),
	}
}
