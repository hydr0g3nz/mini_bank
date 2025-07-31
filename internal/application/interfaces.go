// internal/application/interfaces.go
package usecase

import (
	"context"

	"github.com/hydr0g3nz/mini_bank/internal/application/dto"
)

// AccountUseCase defines the interface for account business logic
type AccountUseCase interface {
	// CreateAccount creates a new account
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.AccountResponse, error)

	// GetAccount retrieves an account by ID
	GetAccount(ctx context.Context, id string) (*dto.AccountResponse, error)

	// UpdateAccount updates an existing account
	UpdateAccount(ctx context.Context, req dto.UpdateAccountRequest) (*dto.AccountResponse, error)

	// DeleteAccount deletes an account
	DeleteAccount(ctx context.Context, id string) error

	// ListAccounts retrieves accounts with pagination
	ListAccounts(ctx context.Context, req dto.ListRequest) (*dto.AccountListResponse, error)

	// SuspendAccount suspends an account
	SuspendAccount(ctx context.Context, id string) error

	// ActivateAccount activates an account
	ActivateAccount(ctx context.Context, id string) error
}

// TransactionUseCase defines the interface for transaction business logic
type TransactionUseCase interface {
	// CreateTransaction creates a new transaction
	CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error)
	ConfirmTransaction(ctx context.Context, req dto.ConfirmTransactionRequest) (*dto.TransactionResponse, error)
	// GetTransaction retrieves a transaction by ID
	GetTransaction(ctx context.Context, id string) (*dto.TransactionResponse, error)

	// ListTransactions retrieves transactions with pagination
	ListTransactions(ctx context.Context, req dto.ListRequest) (*dto.TransactionListResponse, error)

	// GetTransactionsByAccount retrieves transactions for a specific account
	GetTransactionsByAccount(ctx context.Context, accountID string, req dto.ListRequest) (*dto.TransactionListResponse, error)

	// CancelTransaction cancels a transaction
	CancelTransaction(ctx context.Context, req dto.CancelTransactionRequest) error

	// GetTransactionsByStatus retrieves transactions by status
	GetTransactionsByStatus(ctx context.Context, status string, req dto.ListRequest) (*dto.TransactionListResponse, error)
}
