package errors

import "errors"

// Domain errors - independent of HTTP or other adapters
var (
	// ErrNotFound indicates the requested resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput indicates invalid input data
	ErrInvalidInput = errors.New("invalid input")

	// ErrDuplicateEmail indicates the email already exists
	ErrDuplicateEmail = errors.New("email already exists")

	// ErrUnauthorized indicates the user is not authenticated
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates the user doesn't have permission
	ErrForbidden = errors.New("forbidden")

	// ErrInternal indicates an internal server error
	ErrInternal = errors.New("internal error")

	// ErrInvalidStatusTransition indicates an invalid order status transition
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return "validation failed"
	}
	return v[0].Field + ": " + v[0].Message
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) ValidationErrors {
	return ValidationErrors{{Field: field, Message: message}}
}

// Is checks if the error matches the target
func (v ValidationErrors) Is(target error) bool {
	return target == ErrInvalidInput
}

