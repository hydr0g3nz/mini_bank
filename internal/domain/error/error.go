package errs

import (
	"errors"
	"fmt"
)

// Domain Error Types
var (
	// Transaction Errors
	ErrInvalidTransactionAmount     = errors.New("transaction amount must be greater than zero")
	ErrMissingAccountID             = errors.New("account ID is required")
	ErrSameAccountTransfer          = errors.New("from and to account cannot be the same")
	ErrInvalidTransactionStatus     = errors.New("invalid transaction status transition")
	ErrTransactionAlreadyInProgress = errors.New("transaction confirmation already in progress")
	ErrTransactionNotFound          = errors.New("transaction not found")
	ErrTransactionCannotBeConfirmed = errors.New("transaction cannot be confirmed")
	ErrTransactionCannotBeCancelled = errors.New("transaction cannot be cancelled")

	// Account Errors
	ErrAccountNotFound       = errors.New("account not found")
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrAccountAlreadyExists  = errors.New("account already exists")
	ErrAccountCannotTransact = errors.New("account cannot perform transactions")

	// General Errors
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized access")
	ErrInternalError = errors.New("internal server error")
	// validation errors
	ErrInvalidAccountID     = errors.New("invalid account ID format")
	ErrInvalidTransactionID = errors.New("invalid transaction ID format")
	ErrUnsupportedType      = errors.New("unsupported transaction type")
)

// Custom Error Types
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

type BusinessError struct {
	Code    string
	Message string
}

func (e BusinessError) Error() string {
	return fmt.Sprintf("business error [%s]: %s", e.Code, e.Message)
}
