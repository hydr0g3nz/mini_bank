// internal/application/transaction.go
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/hydr0g3nz/mini_bank/internal/application/dto"
	"github.com/hydr0g3nz/mini_bank/internal/domain/entity"
	errs "github.com/hydr0g3nz/mini_bank/internal/domain/error"
	"github.com/hydr0g3nz/mini_bank/internal/domain/infra"
	"github.com/hydr0g3nz/mini_bank/internal/domain/repository"
	"github.com/hydr0g3nz/mini_bank/internal/domain/vo"
)

type transactionUseCase struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
	cache           infra.CacheService
	logger          infra.Logger
	mapper          *dto.TransactionMapper
}

// NewTransactionUseCase creates a new transaction use case
func NewTransactionUseCase(
	transactionRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	cache infra.CacheService,
	logger infra.Logger,
) TransactionUseCase {
	return &transactionUseCase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		cache:           cache,
		logger:          logger,
		mapper:          &dto.TransactionMapper{},
	}
}

// CreateTransaction creates a new transaction
func (uc *transactionUseCase) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (*dto.TransactionResponse, error) {
	uc.logger.Info("Creating new transaction",
		"type", req.TransactionType,
		"amount", req.Amount,
		"fromAccountID", req.FromAccountID,
		"toAccountID", req.ToAccountID)

	// Convert DTO to domain values
	fromAccountID, toAccountID, transactionType, amount, description, reference, err := uc.mapper.FromCreateRequest(req)
	if err != nil {
		uc.logger.Error("Failed to convert create transaction request", "error", err)
		return nil, err
	}

	// Validate accounts exist and can transact
	if err := uc.validateAccountsForTransaction(ctx, fromAccountID, toAccountID, transactionType); err != nil {
		return nil, err
	}

	// Create transaction entity based on type
	var transaction *entity.Transaction
	switch transactionType {
	case vo.TransactionTypeDebit:
		transaction, err = entity.NewDebitTransaction(*fromAccountID, amount, description, reference)
	case vo.TransactionTypeCredit:
		transaction, err = entity.NewCreditTransaction(*toAccountID, amount, description, reference)
	case vo.TransactionTypeTransfer:
		transaction, err = entity.NewTransferTransaction(*fromAccountID, *toAccountID, amount, description, reference)
	default:
		return nil, errs.ErrInvalidInput
	}

	if err != nil {
		uc.logger.Error("Failed to create transaction entity", "error", err)
		return nil, err
	}

	// Save to repository
	if err := uc.transactionRepo.Create(ctx, transaction); err != nil {
		uc.logger.Error("Failed to save transaction to repository", "error", err, "transactionID", transaction.ID.String())
		return nil, err
	}

	// Convert to response DTO
	response := uc.mapper.ToResponse(transaction)

	// Cache the transaction
	cacheKey := fmt.Sprintf("transaction:%s", transaction.ID.String())
	if err := uc.cache.Set(ctx, cacheKey, response, 30*time.Minute); err != nil {
		uc.logger.Warn("Failed to cache transaction", "error", err, "transactionID", transaction.ID.String())
	}

	uc.logger.Info("Transaction created successfully", "transactionID", transaction.ID.String())
	return &response, nil
}

// ConfirmTransaction confirms and processes a transaction (Idempotent)
func (uc *transactionUseCase) ConfirmTransaction(ctx context.Context, req dto.ConfirmTransactionRequest) (*dto.TransactionResponse, error) {
	uc.logger.Info("Confirming transaction", "transactionID", req.ID)

	// Parse transaction ID
	transactionID, err := vo.NewTransactionIDFromString(req.ID)
	if err != nil {
		uc.logger.Error("Invalid transaction ID format", "error", err, "transactionID", req.ID)
		return nil, err
	}

	// Create idempotency key for confirm operation
	idempotencyKey := fmt.Sprintf("confirm_transaction:%s", req.ID)

	// Check if this confirmation has already been processed (idempotency check)
	var cachedResult dto.TransactionResponse
	if err := uc.cache.Get(ctx, idempotencyKey, &cachedResult); err == nil {
		uc.logger.Info("Transaction confirmation already processed (idempotent)", "transactionID", req.ID)
		return &cachedResult, nil
	}

	// Try to acquire distributed lock for this transaction to prevent concurrent processing
	lockKey := fmt.Sprintf("lock:transaction:%s", req.ID)
	lockAcquired, err := uc.acquireDistributedLock(ctx, lockKey, 30*time.Second)
	if err != nil {
		uc.logger.Error("Failed to acquire distributed lock", "error", err, "transactionID", req.ID)
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !lockAcquired {
		uc.logger.Warn("Another confirmation is in progress", "transactionID", req.ID)
		return nil, errs.ErrTransactionAlreadyInProgress
	}

	// Ensure lock is released
	defer func() {
		if err := uc.releaseLock(ctx, lockKey); err != nil {
			uc.logger.Warn("Failed to release distributed lock", "error", err, "transactionID", req.ID)
		}
	}()

	// Get transaction from repository
	transaction, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		uc.logger.Error("Transaction not found", "error", err, "transactionID", req.ID)
		return nil, errs.ErrTransactionNotFound
	}

	// Check if transaction is already completed (idempotency check)
	if transaction.Status.IsCompleted() {
		uc.logger.Info("Transaction already completed", "transactionID", req.ID)
		response := uc.mapper.ToResponse(transaction)

		// Cache the result for future idempotent calls
		if err := uc.cache.Set(ctx, idempotencyKey, response, 24*time.Hour); err != nil {
			uc.logger.Warn("Failed to cache idempotent result", "error", err, "transactionID", req.ID)
		}

		return &response, nil
	}

	// Check if transaction can be confirmed
	if !transaction.Status.CanTransitionTo(vo.TransactionStatusCompleted) {
		uc.logger.Error("Transaction cannot be confirmed", "status", transaction.Status, "transactionID", req.ID)
		return nil, fmt.Errorf("%w in status : %s", errs.ErrTransactionCannotBeConfirmed, transaction.Status)
	}

	// Process the transaction based on type
	if err := uc.processTransaction(ctx, transaction); err != nil {
		// Mark transaction as failed
		if markErr := transaction.MarkAsFailed(); markErr != nil {
			uc.logger.Error("Failed to mark transaction as failed", "error", markErr, "transactionID", req.ID)
		} else {
			uc.transactionRepo.Update(ctx, transaction)
		}

		uc.logger.Error("Failed to process transaction", "error", err, "transactionID", req.ID)
		return nil, err
	}

	// Mark transaction as completed
	if err := transaction.MarkAsCompleted(); err != nil {
		uc.logger.Error("Failed to mark transaction as completed", "error", err, "transactionID", req.ID)
		return nil, err
	}

	// Update transaction in repository
	if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
		uc.logger.Error("Failed to update transaction in repository", "error", err, "transactionID", req.ID)
		return nil, err
	}

	// Convert to response
	response := uc.mapper.ToResponse(transaction)

	// Cache the result for idempotency (longer TTL since it's completed)
	if err := uc.cache.Set(ctx, idempotencyKey, response, 24*time.Hour); err != nil {
		uc.logger.Warn("Failed to cache confirmed transaction result", "error", err, "transactionID", req.ID)
	}

	// Update transaction cache
	transactionCacheKey := fmt.Sprintf("transaction:%s", req.ID)
	if err := uc.cache.Set(ctx, transactionCacheKey, response, 30*time.Minute); err != nil {
		uc.logger.Warn("Failed to update transaction cache", "error", err, "transactionID", req.ID)
	}

	// Invalidate account caches since balances changed
	uc.invalidateAccountCaches(ctx, transaction)

	uc.logger.Info("Transaction confirmed successfully", "transactionID", req.ID)
	return &response, nil
}

// GetTransaction retrieves a transaction by ID
func (uc *transactionUseCase) GetTransaction(ctx context.Context, id string) (*dto.TransactionResponse, error) {
	uc.logger.Debug("Getting transaction", "transactionID", id)

	// Parse transaction ID
	transactionID, err := vo.NewTransactionIDFromString(id)
	if err != nil {
		uc.logger.Error("Invalid transaction ID format", "error", err, "transactionID", id)
		return nil, err
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("transaction:%s", id)
	var cachedResponse dto.TransactionResponse
	if err := uc.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		uc.logger.Debug("Transaction found in cache", "transactionID", id)
		return &cachedResponse, nil
	}

	// Get from repository
	transaction, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		uc.logger.Error("Failed to get transaction from repository", "error", err, "transactionID", id)
		return nil, errs.ErrTransactionNotFound
	}

	// Convert to response DTO
	response := uc.mapper.ToResponse(transaction)

	// Cache the result
	if err := uc.cache.Set(ctx, cacheKey, response, 30*time.Minute); err != nil {
		uc.logger.Warn("Failed to cache transaction", "error", err, "transactionID", id)
	}

	uc.logger.Debug("Transaction retrieved successfully", "transactionID", id)
	return &response, nil
}

// ListTransactions retrieves transactions with pagination
func (uc *transactionUseCase) ListTransactions(ctx context.Context, req dto.ListRequest) (*dto.TransactionListResponse, error) {
	uc.logger.Debug("Listing transactions", "page", req.Page, "pageSize", req.PageSize)

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Try to get from cache first
	cacheKey := fmt.Sprintf("transactions:list:page:%d:size:%d", req.Page, req.PageSize)
	var cachedResponse dto.TransactionListResponse
	if err := uc.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		uc.logger.Debug("Transaction list found in cache")
		return &cachedResponse, nil
	}

	// Get from repository
	transactions, err := uc.transactionRepo.List(ctx, req.PageSize, offset)
	if err != nil {
		uc.logger.Error("Failed to get transactions from repository", "error", err)
		return nil, err
	}

	// Create pagination info
	pagination := dto.PaginationInfo{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: int64(len(transactions)),
		TotalPages: (len(transactions) + req.PageSize - 1) / req.PageSize,
		HasNext:    len(transactions) == req.PageSize,
		HasPrev:    req.Page > 1,
	}

	// Convert to response DTO
	response := uc.mapper.ToResponseList(transactions, pagination)

	// Cache the result for shorter time
	if err := uc.cache.Set(ctx, cacheKey, response, 2*time.Minute); err != nil {
		uc.logger.Warn("Failed to cache transaction list", "error", err)
	}

	uc.logger.Debug("Transaction list retrieved successfully", "count", len(transactions))
	return &response, nil
}

// GetTransactionsByAccount retrieves transactions for a specific account
func (uc *transactionUseCase) GetTransactionsByAccount(ctx context.Context, accountID string, req dto.ListRequest) (*dto.TransactionListResponse, error) {
	uc.logger.Debug("Getting transactions by account", "accountID", accountID, "page", req.Page)

	// Parse account ID
	parsedAccountID, err := vo.NewAccountIDFromString(accountID)
	if err != nil {
		uc.logger.Error("Invalid account ID format", "error", err, "accountID", accountID)
		return nil, err
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Try to get from cache first
	cacheKey := fmt.Sprintf("transactions:account:%s:page:%d:size:%d", accountID, req.Page, req.PageSize)
	var cachedResponse dto.TransactionListResponse
	if err := uc.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		uc.logger.Debug("Account transactions found in cache", "accountID", accountID)
		return &cachedResponse, nil
	}

	// Get from repository
	transactions, err := uc.transactionRepo.GetByAccountID(ctx, parsedAccountID, req.PageSize, offset)
	if err != nil {
		uc.logger.Error("Failed to get transactions by account from repository", "error", err, "accountID", accountID)
		return nil, err
	}

	// Create pagination info
	pagination := dto.PaginationInfo{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: int64(len(transactions)),
		TotalPages: (len(transactions) + req.PageSize - 1) / req.PageSize,
		HasNext:    len(transactions) == req.PageSize,
		HasPrev:    req.Page > 1,
	}

	// Convert to response DTO
	response := uc.mapper.ToResponseList(transactions, pagination)

	// Cache the result
	if err := uc.cache.Set(ctx, cacheKey, response, 5*time.Minute); err != nil {
		uc.logger.Warn("Failed to cache account transactions", "error", err, "accountID", accountID)
	}

	uc.logger.Debug("Account transactions retrieved successfully", "accountID", accountID, "count", len(transactions))
	return &response, nil
}

// CancelTransaction cancels a transaction
func (uc *transactionUseCase) CancelTransaction(ctx context.Context, req dto.CancelTransactionRequest) error {
	uc.logger.Info("Cancelling transaction", "transactionID", req.ID)

	// Parse transaction ID
	transactionID, err := vo.NewTransactionIDFromString(req.ID)
	if err != nil {
		uc.logger.Error("Invalid transaction ID format", "error", err, "transactionID", req.ID)
		return err
	}

	// Get transaction
	transaction, err := uc.transactionRepo.GetByID(ctx, transactionID)
	if err != nil {
		uc.logger.Error("Transaction not found", "error", err, "transactionID", req.ID)
		return errs.ErrTransactionNotFound
	}

	// Check if transaction can be cancelled
	if !transaction.Status.IsPending() {
		uc.logger.Error("Transaction cannot be cancelled", "status", transaction.Status, "transactionID", req.ID)
		return fmt.Errorf("%w in status: %s", errs.ErrTransactionCannotBeCancelled, transaction.Status)
	}

	// Cancel transaction
	if err := transaction.MarkAsCancelled(); err != nil {
		uc.logger.Error("Failed to cancel transaction", "error", err, "transactionID", req.ID)
		return err
	}

	// Update in repository
	if err := uc.transactionRepo.Update(ctx, transaction); err != nil {
		uc.logger.Error("Failed to update cancelled transaction in repository", "error", err, "transactionID", req.ID)
		return err
	}

	// Update cache
	response := uc.mapper.ToResponse(transaction)
	cacheKey := fmt.Sprintf("transaction:%s", req.ID)
	if err := uc.cache.Set(ctx, cacheKey, response, 30*time.Minute); err != nil {
		uc.logger.Warn("Failed to update transaction cache", "error", err, "transactionID", req.ID)
	}

	uc.logger.Info("Transaction cancelled successfully", "transactionID", req.ID)
	return nil
}

// GetTransactionsByStatus retrieves transactions by status
func (uc *transactionUseCase) GetTransactionsByStatus(ctx context.Context, status string, req dto.ListRequest) (*dto.TransactionListResponse, error) {
	uc.logger.Debug("Getting transactions by status", "status", status, "page", req.Page)

	// Parse status
	transactionStatus := vo.TransactionStatus(status)
	if !transactionStatus.IsValid() {
		uc.logger.Error("Invalid transaction status", "status", status)
		return nil, fmt.Errorf("invalid transaction status: %s", status)
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Try to get from cache first
	cacheKey := fmt.Sprintf("transactions:status:%s:page:%d:size:%d", status, req.Page, req.PageSize)
	var cachedResponse dto.TransactionListResponse
	if err := uc.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		uc.logger.Debug("Transactions by status found in cache", "status", status)
		return &cachedResponse, nil
	}

	// Get from repository
	transactions, err := uc.transactionRepo.GetByStatus(ctx, transactionStatus, req.PageSize, offset)
	if err != nil {
		uc.logger.Error("Failed to get transactions by status from repository", "error", err, "status", status)
		return nil, err
	}

	// Create pagination info
	pagination := dto.PaginationInfo{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: int64(len(transactions)),
		TotalPages: (len(transactions) + req.PageSize - 1) / req.PageSize,
		HasNext:    len(transactions) == req.PageSize,
		HasPrev:    req.Page > 1,
	}

	// Convert to response DTO
	response := uc.mapper.ToResponseList(transactions, pagination)

	// Cache the result
	if err := uc.cache.Set(ctx, cacheKey, response, 5*time.Minute); err != nil {
		uc.logger.Warn("Failed to cache transactions by status", "error", err, "status", status)
	}

	uc.logger.Debug("Transactions by status retrieved successfully", "status", status, "count", len(transactions))
	return &response, nil
}

// Helper methods

// validateAccountsForTransaction validates that accounts exist and can perform the transaction
func (uc *transactionUseCase) validateAccountsForTransaction(
	ctx context.Context,
	fromAccountID *vo.AccountID,
	toAccountID *vo.AccountID,
	transactionType vo.TransactionType,
) error {
	switch transactionType {
	case vo.TransactionTypeDebit:
		if fromAccountID == nil {
			return errs.ErrMissingAccountID
		}
		return uc.validateAccountCanTransact(ctx, *fromAccountID)

	case vo.TransactionTypeCredit:
		if toAccountID == nil {
			return errs.ErrMissingAccountID
		}
		return uc.validateAccountCanTransact(ctx, *toAccountID)

	case vo.TransactionTypeTransfer:
		if fromAccountID == nil || toAccountID == nil {
			return errs.ErrMissingAccountID
		}
		if err := uc.validateAccountCanTransact(ctx, *fromAccountID); err != nil {
			return err
		}
		return uc.validateAccountCanTransact(ctx, *toAccountID)
	}

	return nil
}

// validateAccountCanTransact checks if an account exists and can perform transactions
func (uc *transactionUseCase) validateAccountCanTransact(ctx context.Context, accountID vo.AccountID) error {
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		uc.logger.Error("Account not found for transaction validation", "error", err, "accountID", accountID.String())
		return errs.ErrAccountNotFound
	}

	if !account.CanTransact() {
		uc.logger.Error("Account cannot perform transactions", "accountID", accountID.String(), "status", account.Status)
		return fmt.Errorf("%w : %s", errs.ErrAccountCannotTransact, account.Status)
	}

	return nil
}

// processTransaction executes the actual transaction logic
func (uc *transactionUseCase) processTransaction(ctx context.Context, transaction *entity.Transaction) error {
	switch transaction.TransactionType {
	case vo.TransactionTypeDebit:
		return uc.processDebitTransaction(ctx, transaction)
	case vo.TransactionTypeCredit:
		return uc.processCreditTransaction(ctx, transaction)
	case vo.TransactionTypeTransfer:
		return uc.processTransferTransaction(ctx, transaction)
	default:
		return fmt.Errorf("%w : %s", errs.ErrUnsupportedType, transaction.TransactionType)
	}
}

// processDebitTransaction processes a debit transaction
func (uc *transactionUseCase) processDebitTransaction(ctx context.Context, transaction *entity.Transaction) error {
	if transaction.FromAccountID == nil {
		return errs.ErrMissingAccountID
	}

	// Get account
	account, err := uc.accountRepo.GetByID(ctx, *transaction.FromAccountID)
	if err != nil {
		return errs.ErrAccountNotFound
	}

	// Check if account can transact
	if !account.CanTransact() {
		return errs.ErrAccountCannotTransact
	}

	// Perform debit
	if err := account.Debit(transaction.Amount); err != nil {
		return err
	}

	// Update account
	return uc.accountRepo.Update(ctx, account)
}

// processCreditTransaction processes a credit transaction
func (uc *transactionUseCase) processCreditTransaction(ctx context.Context, transaction *entity.Transaction) error {
	if transaction.ToAccountID == nil {
		return errs.ErrMissingAccountID
	}

	// Get account
	account, err := uc.accountRepo.GetByID(ctx, *transaction.ToAccountID)
	if err != nil {
		return errs.ErrAccountNotFound
	}

	// Check if account can transact
	if !account.CanTransact() {
		return errs.ErrAccountCannotTransact
	}

	// Perform credit
	if err := account.Credit(transaction.Amount); err != nil {
		return err
	}

	// Update account
	return uc.accountRepo.Update(ctx, account)
}

// processTransferTransaction processes a transfer transaction
func (uc *transactionUseCase) processTransferTransaction(ctx context.Context, transaction *entity.Transaction) error {
	if transaction.FromAccountID == nil || transaction.ToAccountID == nil {
		return errs.ErrMissingAccountID
	}

	// Get both accounts
	fromAccount, err := uc.accountRepo.GetByID(ctx, *transaction.FromAccountID)
	if err != nil {
		return errs.ErrAccountNotFound
	}

	toAccount, err := uc.accountRepo.GetByID(ctx, *transaction.ToAccountID)
	if err != nil {
		return errs.ErrAccountNotFound
	}

	// Check if both accounts can transact
	if !fromAccount.CanTransact() {
		return errs.ErrAccountCannotTransact
	}
	if !toAccount.CanTransact() {
		return errs.ErrAccountCannotTransact
	}

	// Perform debit from source account
	if err := fromAccount.Debit(transaction.Amount); err != nil {
		return fmt.Errorf("failed to debit from account: %w", err)
	}

	// Perform credit to destination account
	if err := toAccount.Credit(transaction.Amount); err != nil {
		// Rollback the debit if credit fails
		fromAccount.Credit(transaction.Amount) // Ignore error on rollback
		return fmt.Errorf("failed to credit to account: %w", err)
	}

	// Update both accounts
	if err := uc.accountRepo.Update(ctx, fromAccount); err != nil {
		return fmt.Errorf("failed to update from account: %w", err)
	}

	if err := uc.accountRepo.Update(ctx, toAccount); err != nil {
		return fmt.Errorf("failed to update to account: %w", err)
	}

	return nil
}

// acquireDistributedLock acquires a distributed lock using Redis
func (uc *transactionUseCase) acquireDistributedLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	// This is a simplified implementation. In production, consider using a more robust
	// distributed lock implementation like Redlock
	lockValue := fmt.Sprintf("lock_%d", time.Now().UnixNano())

	// Try to set the lock with expiration
	// This should be implemented using Redis SETNX with expiration
	// For now, we'll use a simple cache set operation
	err := uc.cache.Set(ctx, key, lockValue, expiration)
	if err != nil {
		return false, err
	}

	return true, nil
}

// releaseLock releases a distributed lock
func (uc *transactionUseCase) releaseLock(ctx context.Context, key string) error {
	return uc.cache.Delete(ctx, key)
}

// invalidateAccountCaches invalidates account caches after balance changes
func (uc *transactionUseCase) invalidateAccountCaches(ctx context.Context, transaction *entity.Transaction) {
	if transaction.FromAccountID != nil {
		cacheKey := fmt.Sprintf("account:%s", transaction.FromAccountID.String())
		if err := uc.cache.Delete(ctx, cacheKey); err != nil {
			uc.logger.Warn("Failed to invalidate from account cache",
				"error", err,
				"accountID", transaction.FromAccountID.String())
		}
	}

	if transaction.ToAccountID != nil {
		cacheKey := fmt.Sprintf("account:%s", transaction.ToAccountID.String())
		if err := uc.cache.Delete(ctx, cacheKey); err != nil {
			uc.logger.Warn("Failed to invalidate to account cache",
				"error", err,
				"accountID", transaction.ToAccountID.String())
		}
	}

	// Also invalidate account list caches since balances changed
	// In a more sophisticated implementation, you might use cache tags or patterns
	// For now, we'll just log that lists should be invalidated
	uc.logger.Debug("Account balances changed, consider invalidating account list caches")
}
