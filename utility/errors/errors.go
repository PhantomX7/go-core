package errors

import (
	"fmt"
	"net/http"
)

var (
	// ErrTetapTenangTetapSemangat custom error on unexpected error
	ErrTetapTenangTetapSemangat = CustomError{
		Message:  "Tetap Tenang Tetap Semangat",
		HTTPCode: http.StatusInternalServerError,
	}

	ErrBadRequest = CustomError{
		Message:  "Bad Request",
		HTTPCode: http.StatusBadRequest,
	}

	ErrUnauthorized = CustomError{
		Message:  "Unauthorized",
		HTTPCode: http.StatusUnauthorized,
	}

	ErrForbidden = CustomError{
		Message:  "Forbidden",
		HTTPCode: http.StatusForbidden,
	}

	ErrNotFound = CustomError{
		Message:  "Record not exist",
		HTTPCode: http.StatusNotFound,
	}

	ErrUnprocessableEntity = CustomError{
		Message:  "Unprocessable Entity",
		HTTPCode: http.StatusUnprocessableEntity,
	}

	ErrFailedAuthentication = CustomError{
		Message:  "Invalid Credentials",
		HTTPCode: http.StatusUnauthorized,
	}
)

// CustomError holds data for customized error
type CustomError struct {
	Message  interface{} `json:"message"`
	HTTPCode int         `json:"code"`
}

// Error is a function to convert error to string.
// It exists to satisfy error interface
func (c CustomError) Error() string {
	return fmt.Sprint(c.Message)
}
