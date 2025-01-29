package response

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string   `json:"message"`          // user-level status message
	AppCode    int64    `json:"code,omitempty"`   // application-specific error code
	ErrorText  []string `json:"errors,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

var (
	ErrBadRequest          = &ErrResponse{HTTPStatusCode: http.StatusBadRequest, StatusText: "Bad request"}
	ErrNotFound            = &ErrResponse{HTTPStatusCode: http.StatusNotFound, StatusText: "Resource not found"}
	ErrUnprocessableEntity = &ErrResponse{HTTPStatusCode: http.StatusUnprocessableEntity, StatusText: "Validation failed"}
	ErrInternalServerError = &ErrResponse{HTTPStatusCode: http.StatusInternalServerError, StatusText: "Internal server error"}

	// application-level custom errors
	ErrMissingBindedReqObj = &ErrResponse{HTTPStatusCode: http.StatusInternalServerError, StatusText: "Internal server error caused by missing binded payload"}
)

func ErrGeneric(httpStatusCode int, message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: httpStatusCode,
		StatusText:     message,
	}
}

// Error renderer for converting errors of go-playground/validator
func ErrValidationFailed(err error) render.Renderer {
	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value most including myself do not usually have code like this.
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return &ErrResponse{
			HTTPStatusCode: http.StatusUnprocessableEntity,
			StatusText:     err.Error(),
		}
	}

	var errs []string
	for _, err := range err.(validator.ValidationErrors) {
		if err.Param() == "" {
			errs = append(errs, fmt.Sprintf("%s is not a valid '%s' format with type '%s'", err.StructField(), err.Tag(), err.Type()))
		} else {
			errs = append(errs, fmt.Sprintf("%s does not match condition '%s'='%s' with type %s", err.StructField(), err.Tag(), err.Param(), err.Type()))
		}
	}

	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Validation failed",
		ErrorText:      errs,
	}
}

func ErrDuplicateEmail(email string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusConflict,
		StatusText:     fmt.Sprintf("The user with email '%s' has already been registered", email),
	}
}

func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Duplicate Key",
		ErrorText:      []string{err.Error()},
	}
}

func ErrUnauthorized() render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized",
	}
}
