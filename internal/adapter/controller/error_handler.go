package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hydr0g3nz/mini_bank/internal/application/dto"
	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
)

// Custom validation error type
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Global validator instance
var validate = validator.New()

// ValidateStruct validates a struct using the validator package
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		return &ValidationError{
			Field:   "validation",
			Message: strings.Join(validationErrors, ", "),
		}
	}
	return nil
}

// formatValidationError formats validator.FieldError to a readable message
func formatValidationError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())
	tag := err.Tag()

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + err.Param() + " characters long"
	case "max":
		return field + " must be at most " + err.Param() + " characters long"
	case "gt":
		return field + " must be greater than " + err.Param()
	case "gte":
		return field + " must be greater than or equal to " + err.Param()
	case "lt":
		return field + " must be less than " + err.Param()
	case "lte":
		return field + " must be less than or equal to " + err.Param()
	case "oneof":
		return field + " must be one of: " + err.Param()
	default:
		return field + " is invalid"
	}
}

// HandleError handles different types of errors and returns appropriate HTTP responses
func HandleError(ctx *gin.Context, err error) {
	var errorResponse dto.ErrorResponse
	var statusCode int

	switch {
	// Domain-specific errors
	case errors.Is(err, errs.ErrAccountNotFound):
		statusCode = http.StatusNotFound
		errorResponse = dto.ErrorResponse{
			Code:    "ACCOUNT_NOT_FOUND",
			Message: "Account not found",
		}

	case errors.Is(err, errs.ErrAccountAlreadyExists):
		statusCode = http.StatusConflict
		errorResponse = dto.ErrorResponse{
			Code:    "ACCOUNT_ALREADY_EXISTS",
			Message: "Account already exists",
		}

	case errors.Is(err, errs.ErrInsufficientBalance):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "INSUFFICIENT_BALANCE",
			Message: "Insufficient balance for this transaction",
		}

	case errors.Is(err, errs.ErrAccountCannotTransact):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "ACCOUNT_CANNOT_TRANSACT",
			Message: "Account cannot perform transactions",
		}

	case errors.Is(err, errs.ErrTransactionNotFound):
		statusCode = http.StatusNotFound
		errorResponse = dto.ErrorResponse{
			Code:    "TRANSACTION_NOT_FOUND",
			Message: "Transaction not found",
		}

	case errors.Is(err, errs.ErrInvalidTransactionAmount):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "INVALID_TRANSACTION_AMOUNT",
			Message: "Transaction amount must be greater than zero",
		}

	case errors.Is(err, errs.ErrSameAccountTransfer):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "SAME_ACCOUNT_TRANSFER",
			Message: "Cannot transfer to the same account",
		}

	case errors.Is(err, errs.ErrTransactionCannotBeConfirmed):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "TRANSACTION_CANNOT_BE_CONFIRMED",
			Message: "Transaction cannot be confirmed in its current state",
		}

	case errors.Is(err, errs.ErrTransactionCannotBeCancelled):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "TRANSACTION_CANNOT_BE_CANCELLED",
			Message: "Transaction cannot be cancelled in its current state",
		}

	case errors.Is(err, errs.ErrTransactionAlreadyInProgress):
		statusCode = http.StatusConflict
		errorResponse = dto.ErrorResponse{
			Code:    "TRANSACTION_IN_PROGRESS",
			Message: "Transaction confirmation is already in progress",
		}

	case errors.Is(err, errs.ErrMissingAccountID):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "MISSING_ACCOUNT_ID",
			Message: "Account ID is required for this transaction type",
		}

	case errors.Is(err, errs.ErrInvalidAccountID):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "INVALID_ACCOUNT_ID",
			Message: "Invalid account ID format",
		}

	case errors.Is(err, errs.ErrInvalidTransactionID):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "INVALID_TRANSACTION_ID",
			Message: "Invalid transaction ID format",
		}

	case errors.Is(err, errs.ErrUnsupportedType):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "UNSUPPORTED_TRANSACTION_TYPE",
			Message: "Unsupported transaction type",
		}

	case errors.Is(err, errs.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		errorResponse = dto.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: "Invalid input provided",
		}

	case errors.Is(err, errs.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		errorResponse = dto.ErrorResponse{
			Code:    "UNAUTHORIZED",
			Message: "Unauthorized access",
		}

	// Custom error types
	default:
		var validationErr *ValidationError
		var businessErr errs.BusinessError
		var domainValidationErr errs.ValidationError

		switch {
		case errors.As(err, &validationErr):
			statusCode = http.StatusBadRequest
			errorResponse = dto.ErrorResponse{
				Code:    "VALIDATION_ERROR",
				Message: validationErr.Message,
				Details: map[string]string{
					"field": validationErr.Field,
				},
			}

		case errors.As(err, &businessErr):
			statusCode = http.StatusBadRequest
			errorResponse = dto.ErrorResponse{
				Code:    businessErr.Code,
				Message: businessErr.Message,
			}

		case errors.As(err, &domainValidationErr):
			statusCode = http.StatusBadRequest
			errorResponse = dto.ErrorResponse{
				Code:    "DOMAIN_VALIDATION_ERROR",
				Message: domainValidationErr.Message,
				Details: map[string]string{
					"field": domainValidationErr.Field,
				},
			}

		// JSON binding errors
		case strings.Contains(err.Error(), "cannot unmarshal"):
			statusCode = http.StatusBadRequest
			errorResponse = dto.ErrorResponse{
				Code:    "INVALID_JSON",
				Message: "Invalid JSON format",
			}

		case strings.Contains(err.Error(), "required"):
			statusCode = http.StatusBadRequest
			errorResponse = dto.ErrorResponse{
				Code:    "MISSING_REQUIRED_FIELD",
				Message: "Required field is missing",
			}

		// Default internal server error
		default:
			statusCode = http.StatusInternalServerError
			errorResponse = dto.ErrorResponse{
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			}
		}
	}

	ctx.JSON(statusCode, errorResponse)
}
