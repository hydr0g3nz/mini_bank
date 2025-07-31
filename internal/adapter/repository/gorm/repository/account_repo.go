package repository

import (
	"context"
	"errors"

	"github.com/hydr0g3nz/mini_bank/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/hydr0g3nz/mini_bank/internal/domain/repository"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
	"gorm.io/gorm"
)

type AccountRepositoryImpl struct {
	db *gorm.DB
}

// NewAccountRepository creates a new instance of AccountRepositoryImpl
func NewAccountRepository(db *gorm.DB) repository.AccountRepository {
	return &AccountRepositoryImpl{db: db}
}

// Create creates a new account
func (r *AccountRepositoryImpl) Create(ctx context.Context, account *entity.Account) error {
	accountModel := model.FromDomainAccount(account)

	if err := r.db.WithContext(ctx).Create(accountModel).Error; err != nil {
		// Handle duplicate key constraint
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errs.ErrAccountAlreadyExists
		}
		return err
	}

	return nil
}

// GetByID retrieves an account by ID
func (r *AccountRepositoryImpl) GetByID(ctx context.Context, id vo.AccountID) (*entity.Account, error) {
	var accountModel model.Account

	err := r.db.WithContext(ctx).
		Where("account_id = ?", id.String()).
		First(&accountModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrAccountNotFound
		}
		return nil, err
	}

	return accountModel.ToDomainAccount()
}

// Update updates an existing account
func (r *AccountRepositoryImpl) Update(ctx context.Context, account *entity.Account) error {
	var existingModel model.Account

	// First, find the existing record by account_id
	err := r.db.WithContext(ctx).
		Where("account_id = ?", account.ID.String()).
		First(&existingModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrAccountNotFound
		}
		return err
	}

	// Update the existing model with domain data
	existingModel.UpdateFromDomain(account)

	// Save the updates
	if err := r.db.WithContext(ctx).Save(&existingModel).Error; err != nil {
		return err
	}

	return nil
}

// Delete deletes an account by ID (soft delete)
func (r *AccountRepositoryImpl) Delete(ctx context.Context, id vo.AccountID) error {
	result := r.db.WithContext(ctx).
		Where("account_id = ?", id.String()).
		Delete(&model.Account{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrAccountNotFound
	}

	return nil
}

// List retrieves accounts with pagination
func (r *AccountRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Account, error) {
	var accountModels []model.Account

	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&accountModels).Error

	if err != nil {
		return nil, err
	}

	// Convert models to domain entities
	accounts := make([]*entity.Account, len(accountModels))
	for i, accountModel := range accountModels {
		domainAccount, err := accountModel.ToDomainAccount()
		if err != nil {
			return nil, err
		}
		accounts[i] = domainAccount
	}

	return accounts, nil
}

// GetByAccountName retrieves an account by account name
func (r *AccountRepositoryImpl) GetByAccountName(ctx context.Context, accountName string) (*entity.Account, error) {
	var accountModel model.Account

	err := r.db.WithContext(ctx).
		Where("account_name = ?", accountName).
		First(&accountModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrAccountNotFound
		}
		return nil, err
	}

	return accountModel.ToDomainAccount()
}
