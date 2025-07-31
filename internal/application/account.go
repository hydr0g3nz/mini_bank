// internal/application/account.go
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

type accountUseCase struct {
	accountRepo repository.AccountRepository
	cache       infra.CacheService
	logger      infra.Logger
	mapper      *dto.AccountMapper
}

// NewAccountUseCase creates a new account use case
func NewAccountUseCase(
	accountRepo repository.AccountRepository,
	cache infra.CacheService,
	logger infra.Logger,
) AccountUseCase {
	return &accountUseCase{
		accountRepo: accountRepo,
		cache:       cache,
		logger:      logger,
		mapper:      &dto.AccountMapper{},
	}
}

// CreateAccount creates a new account
func (uc *accountUseCase) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	// Log the operation
	uc.logger.Info("Creating new account", "accountName", req.AccountName, "initialBalance", req.InitialBalance)

	// Convert DTO to domain values
	accountName, money, err := uc.mapper.FromCreateRequest(req)
	if err != nil {
		uc.logger.Error("Failed to convert create request", "error", err)
		return nil, err
	}

	// Check if account with same name already exists
	existingAccount, err := uc.accountRepo.GetByAccountName(ctx, accountName)
	if err == nil && existingAccount != nil {
		uc.logger.Warn("Account with same name already exists", "accountName", accountName)
		return nil, errs.ErrAccountAlreadyExists
	}

	// Create new account entity
	account, err := entity.NewAccount(accountName, money)
	if err != nil {
		uc.logger.Error("Failed to create account entity", "error", err)
		return nil, err
	}

	// Save to repository
	if err := uc.accountRepo.Create(ctx, account); err != nil {
		uc.logger.Error("Failed to save account to repository", "error", err, "accountID", account.ID.String())
		return nil, err
	}

	// Convert to response DTO
	response := uc.mapper.ToResponse(account)

	// Cache the account
	cacheKey := fmt.Sprintf("account:%s", account.ID.String())
	if err := uc.cache.Set(ctx, cacheKey, response, 15*time.Minute); err != nil {
		uc.logger.Warn("Failed to cache account", "error", err, "accountID", account.ID.String())

	}

	uc.logger.Info("Account created successfully", "accountID", account.ID.String(), "accountName", accountName)
	return &response, nil
}

// GetAccount retrieves an account by ID
func (uc *accountUseCase) GetAccount(ctx context.Context, id string) (*dto.AccountResponse, error) {
	uc.logger.Debug("Getting account", "accountID", id)

	// Parse account ID
	accountID, err := vo.NewAccountIDFromString(id)
	if err != nil {
		uc.logger.Error("Invalid account ID format", "error", err, "accountID", id)
		return nil, err
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("account:%s", id)
	var cachedResponse dto.AccountResponse
	if err := uc.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		uc.logger.Debug("Account found in cache", "accountID", id)
		return &cachedResponse, nil
	}

	// Get from repository
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		uc.logger.Error("Failed to get account from repository", "error", err, "accountID", id)
		return nil, errs.ErrAccountNotFound
	}

	// Convert to response DTO
	response := uc.mapper.ToResponse(account)

	// Cache the result
	if err := uc.cache.Set(ctx, cacheKey, response, 15*time.Minute); err != nil {
		uc.logger.Warn("Failed to cache account", "error", err, "accountID", id)
	}

	uc.logger.Debug("Account retrieved successfully", "accountID", id)
	return &response, nil
}

// UpdateAccount updates an existing account
func (uc *accountUseCase) UpdateAccount(ctx context.Context, req dto.UpdateAccountRequest) (*dto.AccountResponse, error) {
	uc.logger.Info("Updating account", "accountID", req.ID, "newName", req.AccountName)

	// Parse account ID
	accountID, err := vo.NewAccountIDFromString(req.ID)
	if err != nil {
		uc.logger.Error("Invalid account ID format", "error", err, "accountID", req.ID)
		return nil, err
	}

	// Get existing account
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		uc.logger.Error("Account not found", "error", err, "accountID", req.ID)
		return nil, errs.ErrAccountNotFound
	}

	// Update account name
	account.AccountName = req.AccountName
	account.UpdatedAt = time.Now()

	// Save to repository
	if err := uc.accountRepo.Update(ctx, account); err != nil {
		uc.logger.Error("Failed to update account in repository", "error", err, "accountID", req.ID)
		return nil, err
	}

	// Convert to response DTO
	response := uc.mapper.ToResponse(account)

	// Update cache
	cacheKey := fmt.Sprintf("account:%s", req.ID)
	if err := uc.cache.Set(ctx, cacheKey, response, 15*time.Minute); err != nil {
		uc.logger.Warn("Failed to update account cache", "error", err, "accountID", req.ID)
	}

	uc.logger.Info("Account updated successfully", "accountID", req.ID)
	return &response, nil
}

// DeleteAccount deletes an account
func (uc *accountUseCase) DeleteAccount(ctx context.Context, id string) error {
	uc.logger.Info("Deleting account", "accountID", id)

	// Parse account ID
	accountID, err := vo.NewAccountIDFromString(id)
	if err != nil {
		uc.logger.Error("Invalid account ID format", "error", err, "accountID", id)
		return err
	}

	// Check if account exists
	_, err = uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		uc.logger.Error("Account not found", "error", err, "accountID", id)
		return errs.ErrAccountNotFound
	}

	// Delete from repository
	if err := uc.accountRepo.Delete(ctx, accountID); err != nil { // todo:soft delete
		uc.logger.Error("Failed to delete account from repository", "error", err, "accountID", id)
		return err
	}

	// Remove from cache
	cacheKey := fmt.Sprintf("account:%s", id)
	if err := uc.cache.Delete(ctx, cacheKey); err != nil {
		uc.logger.Warn("Failed to delete account from cache", "error", err, "accountID", id)
	}

	uc.logger.Info("Account deleted successfully", "accountID", id)
	return nil
}

// ListAccounts retrieves accounts with pagination
func (uc *accountUseCase) ListAccounts(ctx context.Context, req dto.ListRequest) (*dto.AccountListResponse, error) {
	uc.logger.Debug("Listing accounts", "page", req.Page, "pageSize", req.PageSize)

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// Try to get from cache first
	cacheKey := fmt.Sprintf("accounts:list:page:%d:size:%d:search:%s", req.Page, req.PageSize, req.Search)
	var cachedResponse dto.AccountListResponse
	if err := uc.cache.Get(ctx, cacheKey, &cachedResponse); err == nil {
		uc.logger.Debug("Account list found in cache")
		return &cachedResponse, nil
	}

	// Get from repository
	accounts, err := uc.accountRepo.List(ctx, req.PageSize, offset)
	if err != nil {
		uc.logger.Error("Failed to get accounts from repository", "error", err)
		return nil, err
	}

	// Create pagination info (simplified - in real implementation, you'd get total count)
	pagination := dto.PaginationInfo{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalItems: int64(len(accounts)), // This should be actual total count from DB
		TotalPages: (len(accounts) + req.PageSize - 1) / req.PageSize,
		HasNext:    len(accounts) == req.PageSize,
		HasPrev:    req.Page > 1,
	}

	// Convert to response DTO
	response := uc.mapper.ToResponseList(accounts, pagination)

	// Cache the result for shorter time since it's a list
	if err := uc.cache.Set(ctx, cacheKey, response, 5*time.Minute); err != nil {
		uc.logger.Warn("Failed to cache account list", "error", err)
	}

	uc.logger.Debug("Account list retrieved successfully", "count", len(accounts))
	return &response, nil
}

// SuspendAccount suspends an account
func (uc *accountUseCase) SuspendAccount(ctx context.Context, id string) error {
	uc.logger.Info("Suspending account", "accountID", id)

	// Parse account ID
	accountID, err := vo.NewAccountIDFromString(id)
	if err != nil {
		uc.logger.Error("Invalid account ID format", "error", err, "accountID", id)
		return err
	}

	// Get account
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		uc.logger.Error("Account not found", "error", err, "accountID", id)
		return errs.ErrAccountNotFound
	}

	// Suspend account
	if err := account.Suspend(); err != nil {
		uc.logger.Error("Failed to suspend account", "error", err, "accountID", id)
		return err
	}

	// Save to repository
	if err := uc.accountRepo.Update(ctx, account); err != nil {
		uc.logger.Error("Failed to update account in repository", "error", err, "accountID", id)
		return err
	}

	// Update cache
	response := uc.mapper.ToResponse(account)
	cacheKey := fmt.Sprintf("account:%s", id)
	if err := uc.cache.Set(ctx, cacheKey, response, 15*time.Minute); err != nil {
		uc.logger.Warn("Failed to update account cache", "error", err, "accountID", id)
	}

	uc.logger.Info("Account suspended successfully", "accountID", id)
	return nil
}

// ActivateAccount activates an account
func (uc *accountUseCase) ActivateAccount(ctx context.Context, id string) error {
	uc.logger.Info("Activating account", "accountID", id)

	// Parse account ID
	accountID, err := vo.NewAccountIDFromString(id)
	if err != nil {
		uc.logger.Error("Invalid account ID format", "error", err, "accountID", id)
		return err
	}

	// Get account
	account, err := uc.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		uc.logger.Error("Account not found", "error", err, "accountID", id)
		return errs.ErrAccountNotFound
	}

	// Activate account
	if err := account.Activate(); err != nil {
		uc.logger.Error("Failed to activate account", "error", err, "accountID", id)
		return err
	}

	// Save to repository
	if err := uc.accountRepo.Update(ctx, account); err != nil {
		uc.logger.Error("Failed to update account in repository", "error", err, "accountID", id)
		return err
	}

	// Update cache
	response := uc.mapper.ToResponse(account)
	cacheKey := fmt.Sprintf("account:%s", id)
	if err := uc.cache.Set(ctx, cacheKey, response, 15*time.Minute); err != nil {
		uc.logger.Warn("Failed to update account cache", "error", err, "accountID", id)
	}

	uc.logger.Info("Account activated successfully", "accountID", id)
	return nil
}
