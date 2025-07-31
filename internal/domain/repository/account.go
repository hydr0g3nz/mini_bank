// internal/domain/repository/repository.go (updated)
package repository

import (
	"context"

	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
)

type AccountRepository interface {
	// Create creates a new account
	Create(ctx context.Context, account *entity.Account) error

	// GetByID retrieves an account by ID
	GetByID(ctx context.Context, id vo.AccountID) (*entity.Account, error)

	// Update updates an existing account
	Update(ctx context.Context, account *entity.Account) error

	// Delete deletes an account by ID
	Delete(ctx context.Context, id vo.AccountID) error

	// List retrieves accounts with pagination
	List(ctx context.Context, limit, offset int) ([]*entity.Account, error)

	// GetByAccountName retrieves an account by account name
	GetByAccountName(ctx context.Context, accountName string) (*entity.Account, error)
}
