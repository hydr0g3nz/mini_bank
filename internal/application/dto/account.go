// internal/application/dto/account.go
package dto

import (
	"time"
)

// CreateAccountRequest represents the request to create a new account
type CreateAccountRequest struct {
	AccountName    string  `json:"account_name" validate:"required,min=1,max=100"`
	InitialBalance float64 `json:"initial_balance" validate:"min=0"`
}

// UpdateAccountRequest represents the request to update an account
type UpdateAccountRequest struct {
	ID          string `json:"id" validate:"required"`
	AccountName string `json:"account_name" validate:"required,min=1,max=100"`
}

// AccountResponse represents the response structure for account data
type AccountResponse struct {
	ID          string    `json:"id"`
	AccountName string    `json:"account_name"`
	Balance     float64   `json:"balance"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AccountListResponse represents paginated account list response
type AccountListResponse struct {
	Accounts   []AccountResponse `json:"accounts"`
	Pagination PaginationInfo    `json:"pagination"`
}
