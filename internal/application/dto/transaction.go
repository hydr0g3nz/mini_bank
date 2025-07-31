// internal/application/dto/transaction.go
package dto

import (
	"time"
)

// CreateTransactionRequest represents the request to create a new transaction
type CreateTransactionRequest struct {
	FromAccountID   *string `json:"from_account_id,omitempty"`
	ToAccountID     *string `json:"to_account_id,omitempty"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=DEBIT CREDIT TRANSFER"`
	Amount          float64 `json:"amount" validate:"required,gt=0"`
	Description     string  `json:"description" validate:"max=500"`
	Reference       string  `json:"reference" validate:"max=100"`
}

// TransactionResponse represents the response structure for transaction data
type TransactionResponse struct {
	ID              string     `json:"id"`
	FromAccountID   *string    `json:"from_account_id,omitempty"`
	ToAccountID     *string    `json:"to_account_id,omitempty"`
	TransactionType string     `json:"transaction_type"`
	Amount          float64    `json:"amount"`
	Description     string     `json:"description"`
	Reference       string     `json:"reference"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

// TransactionListResponse represents paginated transaction list response
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Pagination   PaginationInfo        `json:"pagination"`
}

// ProcessTransactionRequest represents the request to process a transaction
type ConfirmTransactionRequest struct {
	ID string `json:"id" validate:"required"`
}

// CancelTransactionRequest represents the request to cancel a transaction
type CancelTransactionRequest struct {
	ID string `json:"id" validate:"required"`
}
