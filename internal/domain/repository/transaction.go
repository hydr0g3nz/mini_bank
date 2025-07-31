package repository

import (
	"context"

	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
)

type TransactionRepository interface {
	// Create creates a new transaction
	Create(ctx context.Context, transaction *entity.Transaction) error

	// GetByID retrieves a transaction by ID
	GetByID(ctx context.Context, id vo.TransactionID) (*entity.Transaction, error)

	// Update updates an existing transaction
	Update(ctx context.Context, transaction *entity.Transaction) error

	// List retrieves transactions with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.Transaction, error)

	// GetByAccountID retrieves transactions for a specific account
	GetByAccountID(ctx context.Context, accountID vo.AccountID, limit, offset int) ([]*entity.Transaction, error)

	// GetByStatus retrieves transactions by status
	GetByStatus(ctx context.Context, status vo.TransactionStatus, limit, offset int) ([]*entity.Transaction, error)
}
