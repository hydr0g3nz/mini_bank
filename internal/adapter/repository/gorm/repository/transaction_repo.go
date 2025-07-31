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

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new instance of TransactionRepositoryImpl
func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

// Create creates a new transaction
func (r *TransactionRepositoryImpl) Create(ctx context.Context, transaction *entity.Transaction) error {
	transactionModel := model.FromDomainTransaction(transaction)

	if err := r.db.WithContext(ctx).Create(transactionModel).Error; err != nil {
		// Handle duplicate key constraint
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New("transaction with same ID already exists")
		}
		return err
	}

	return nil
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepositoryImpl) GetByID(ctx context.Context, id vo.TransactionID) (*entity.Transaction, error) {
	var transactionModel model.Transaction

	err := r.db.WithContext(ctx).
		Where("transaction_id = ?", id.String()).
		First(&transactionModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrTransactionNotFound
		}
		return nil, err
	}

	return transactionModel.ToDomainTransaction()
}

// Update updates an existing transaction
func (r *TransactionRepositoryImpl) Update(ctx context.Context, transaction *entity.Transaction) error {
	var existingModel model.Transaction

	// First, find the existing record by transaction_id
	err := r.db.WithContext(ctx).
		Where("transaction_id = ?", transaction.ID.String()).
		First(&existingModel).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrTransactionNotFound
		}
		return err
	}

	// Update the existing model with domain data
	existingModel.UpdateFromDomain(transaction)

	// Save the updates
	if err := r.db.WithContext(ctx).Save(&existingModel).Error; err != nil {
		return err
	}

	return nil
}

// List retrieves transactions with pagination
func (r *TransactionRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.Transaction, error) {
	var transactionModels []model.Transaction

	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&transactionModels).Error

	if err != nil {
		return nil, err
	}

	// Convert models to domain entities
	transactions := make([]*entity.Transaction, len(transactionModels))
	for i, transactionModel := range transactionModels {
		domainTransaction, err := transactionModel.ToDomainTransaction()
		if err != nil {
			return nil, err
		}
		transactions[i] = domainTransaction
	}

	return transactions, nil
}

// GetByAccountID retrieves transactions for a specific account
func (r *TransactionRepositoryImpl) GetByAccountID(ctx context.Context, accountID vo.AccountID, limit, offset int) ([]*entity.Transaction, error) {
	var transactionModels []model.Transaction

	accountIDStr := accountID.String()
	err := r.db.WithContext(ctx).
		Where("from_account_id = ? OR to_account_id = ?", accountIDStr, accountIDStr).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&transactionModels).Error

	if err != nil {
		return nil, err
	}

	// Convert models to domain entities
	transactions := make([]*entity.Transaction, len(transactionModels))
	for i, transactionModel := range transactionModels {
		domainTransaction, err := transactionModel.ToDomainTransaction()
		if err != nil {
			return nil, err
		}
		transactions[i] = domainTransaction
	}

	return transactions, nil
}

// GetByStatus retrieves transactions by status
func (r *TransactionRepositoryImpl) GetByStatus(ctx context.Context, status vo.TransactionStatus, limit, offset int) ([]*entity.Transaction, error) {
	var transactionModels []model.Transaction

	err := r.db.WithContext(ctx).
		Where("status = ?", string(status)).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&transactionModels).Error

	if err != nil {
		return nil, err
	}

	// Convert models to domain entities
	transactions := make([]*entity.Transaction, len(transactionModels))
	for i, transactionModel := range transactionModels {
		domainTransaction, err := transactionModel.ToDomainTransaction()
		if err != nil {
			return nil, err
		}
		transactions[i] = domainTransaction
	}

	return transactions, nil
}
