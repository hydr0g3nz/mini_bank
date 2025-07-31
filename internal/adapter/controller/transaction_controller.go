package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	usecase "github.com/hydr0g3nz/mini_bank/internal/application"
	"github.com/hydr0g3nz/mini_bank/internal/application/dto"
	"github.com/hydr0g3nz/mini_bank/internal/domain/infra"
)

type TransactionController struct {
	transactionUseCase usecase.TransactionUseCase
	logger             infra.Logger
}

func NewTransactionController(transactionUseCase usecase.TransactionUseCase, logger infra.Logger) *TransactionController {
	return &TransactionController{
		transactionUseCase: transactionUseCase,
		logger:             logger,
	}
}

// CreateTransaction creates a new transaction
func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var req dto.CreateTransactionRequest
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

	response, err := c.transactionUseCase.CreateTransaction(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error("Failed to create transaction", "error", err)
		HandleError(ctx, err)
		return
	}

	c.logger.Info("Transaction created successfully", "transactionID", response.ID)
	ctx.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Transaction created successfully",
		Data:    response,
	})
}

// ConfirmTransaction confirms and processes a transaction
func (c *TransactionController) ConfirmTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		c.logger.Error("Transaction ID is required")
		HandleError(ctx, &ValidationError{Field: "id", Message: "transaction ID is required"})
		return
	}

	req := dto.ConfirmTransactionRequest{ID: id}

	response, err := c.transactionUseCase.ConfirmTransaction(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error("Failed to confirm transaction", "error", err, "transactionID", id)
		HandleError(ctx, err)
		return
	}

	c.logger.Info("Transaction confirmed successfully", "transactionID", id)
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Transaction confirmed successfully",
		Data:    response,
	})
}

// GetTransaction retrieves a transaction by ID
func (c *TransactionController) GetTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		c.logger.Error("Transaction ID is required")
		HandleError(ctx, &ValidationError{Field: "id", Message: "transaction ID is required"})
		return
	}

	response, err := c.transactionUseCase.GetTransaction(ctx.Request.Context(), id)
	if err != nil {
		c.logger.Error("Failed to get transaction", "error", err, "transactionID", id)
		HandleError(ctx, err)
		return
	}

	c.logger.Debug("Transaction retrieved successfully", "transactionID", id)
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Transaction retrieved successfully",
		Data:    response,
	})
}

// ListTransactions retrieves transactions with pagination
func (c *TransactionController) ListTransactions(ctx *gin.Context) {
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

	response, err := c.transactionUseCase.ListTransactions(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error("Failed to list transactions", "error", err)
		HandleError(ctx, err)
		return
	}

	c.logger.Debug("Transactions listed successfully", "count", len(response.Transactions))
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Transactions retrieved successfully",
		Data:    response,
	})
}

// GetTransactionsByAccount retrieves transactions for a specific account
func (c *TransactionController) GetTransactionsByAccount(ctx *gin.Context) {
	accountID := ctx.Param("id")
	if accountID == "" {
		c.logger.Error("Account ID is required")
		HandleError(ctx, &ValidationError{Field: "account_id", Message: "account ID is required"})
		return
	}

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

	response, err := c.transactionUseCase.GetTransactionsByAccount(ctx.Request.Context(), accountID, req)
	if err != nil {
		c.logger.Error("Failed to get transactions by account", "error", err, "accountID", accountID)
		HandleError(ctx, err)
		return
	}

	c.logger.Debug("Account transactions retrieved successfully", "accountID", accountID, "count", len(response.Transactions))
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Account transactions retrieved successfully",
		Data:    response,
	})
}

// CancelTransaction cancels a transaction
func (c *TransactionController) CancelTransaction(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		c.logger.Error("Transaction ID is required")
		HandleError(ctx, &ValidationError{Field: "id", Message: "transaction ID is required"})
		return
	}

	req := dto.CancelTransactionRequest{ID: id}

	err := c.transactionUseCase.CancelTransaction(ctx.Request.Context(), req)
	if err != nil {
		c.logger.Error("Failed to cancel transaction", "error", err, "transactionID", id)
		HandleError(ctx, err)
		return
	}

	c.logger.Info("Transaction cancelled successfully", "transactionID", id)
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Transaction cancelled successfully",
	})
}

// GetTransactionsByStatus retrieves transactions by status
func (c *TransactionController) GetTransactionsByStatus(ctx *gin.Context) {
	status := ctx.Param("status")
	if status == "" {
		c.logger.Error("Transaction status is required")
		HandleError(ctx, &ValidationError{Field: "status", Message: "transaction status is required"})
		return
	}

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

	response, err := c.transactionUseCase.GetTransactionsByStatus(ctx.Request.Context(), status, req)
	if err != nil {
		c.logger.Error("Failed to get transactions by status", "error", err, "status", status)
		HandleError(ctx, err)
		return
	}

	c.logger.Debug("Transactions by status retrieved successfully", "status", status, "count", len(response.Transactions))
	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Transactions by status retrieved successfully",
		Data:    response,
	})
}
