package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	usecase "github.com/hydr0g3nz/mini_bank/internal/application"
	"github.com/hydr0g3nz/mini_bank/internal/application/dto"
	"github.com/hydr0g3nz/mini_bank/internal/domain/infra"
)

type AccountController struct {
	accountUseCase usecase.AccountUseCase
	logger         infra.Logger
}

func NewAccountController(accountUseCase usecase.AccountUseCase, logger infra.Logger) *AccountController {
	return &AccountController{
		accountUseCase: accountUseCase,
		logger:         logger,
	}
}

// CreateAccount creates a new account
func (c *AccountController) CreateAccount(ctx *gin.Context) {
	var req dto.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error("Failed to bind JSON", "error", err)
		HandleError(ctx, err)
		return
	}

	// Validate request
	if err := ValidateStruct(req); err != nil {
		c.logger.Error("Validation failed", "error", err)
		HandleError(ctx, err)
		return
	}

	response, err := c.accountUseCase.CreateAccount(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error("Failed to create account", "error", err)
		HandleError(ctx, err)
		return
	}

	c.logger.Info("Account created successfully", "accountID", response.ID)
	ctx.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Account created successfully",
		Data:    response,
	})
}

// GetAccount retrieves an account by ID
func (c *AccountController) GetAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		c.logger.Error("Account ID is required")
		HandleError(ctx, &ValidationError{Field: "id", Message: "account ID is required"})
		return
	}

	response, err := c.accountUseCase.GetAccount(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Error("Failed to get account", "error", err, "accountID", id)
		HandleError(ctx, err)
		return
	}

	c.logger.Debug("Account retrieved successfully", "accountID", id)
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Account retrieved successfully",
		Data:    response,
	})
}

// UpdateAccount updates an existing account
func (c *AccountController) UpdateAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		c.logger.Error("Account ID is required")
		HandleError(ctx, &ValidationError{Field: "id", Message: "account ID is required"})
		return
	}

	var req dto.UpdateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error("Failed to bind JSON", "error", err)
		HandleError(ctx, err)
		return
	}

	// Set ID from URL parameter
	req.ID = id

	// Validate request
	if err := ValidateStruct(req); err != nil {
		c.logger.Error("Validation failed", "error", err)
		HandleError(ctx, err)
		return
	}

	response, err := c.accountUseCase.UpdateAccount(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error("Failed to update account", "error", err, "accountID", id)
		HandleError(ctx, err)
		return
	}

	c.logger.Info("Account updated successfully", "accountID", id)
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Account updated successfully",
		Data:    response,
	})
}

// DeleteAccount deletes an account
func (c *AccountController) DeleteAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		c.logger.Error("Account ID is required")
		HandleError(ctx, &ValidationError{Field: "id", Message: "account ID is required"})
		return
	}

	err := c.accountUseCase.DeleteAccount(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Error("Failed to delete account", "error", err, "accountID", id)
		HandleError(ctx, err)
		return
	}

	c.logger.Info("Account deleted successfully", "accountID", id)
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Account deleted successfully",
	})
}

// ListAccounts retrieves accounts with pagination
func (c *AccountController) ListAccounts(ctx *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	search := ctx.Query("search")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	sortDir := ctx.DefaultQuery("sort_dir", "desc")

	req := dto.ListRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}

	// Validate request
	if err := ValidateStruct(req); err != nil {
		c.logger.Error("Validation failed", "error", err)
		HandleError(ctx, err)
		return
	}

	response, err := c.accountUseCase.ListAccounts(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error("Failed to list accounts", "error", err)
		HandleError(ctx, err)
		return
	}

	c.logger.Debug("Accounts listed successfully", "count", len(response.Accounts))
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Accounts retrieved successfully",
		Data:    response,
	})
}

// SuspendAccount suspends an account
func (c *AccountController) SuspendAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		c.logger.Error("Account ID is required")
		HandleError(ctx, &ValidationError{Field: "id", Message: "account ID is required"})
		return
	}

	err := c.accountUseCase.SuspendAccount(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Error("Failed to suspend account", "error", err, "accountID", id)
		HandleError(ctx, err)
		return
	}

	c.logger.Info("Account suspended successfully", "accountID", id)
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Account suspended successfully",
	})
}

// ActivateAccount activates an account
func (c *AccountController) ActivateAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		c.logger.Error("Account ID is required")
		HandleError(ctx, &ValidationError{Field: "id", Message: "account ID is required"})
		return
	}

	err := c.accountUseCase.ActivateAccount(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Error("Failed to activate account", "error", err, "accountID", id)
		HandleError(ctx, err)
		return
	}

	c.logger.Info("Account activated successfully", "accountID", id)
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Account activated successfully",
	})
}
