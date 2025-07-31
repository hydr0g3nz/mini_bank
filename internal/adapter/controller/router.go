package controller

import (
	"github.com/gin-gonic/gin"
	usecase "github.com/hydr0g3nz/mini_bank/internal/application"
	"github.com/hydr0g3nz/mini_bank/internal/domain/infra"
)

type RouterConfig struct {
	APIKey string
	Logger infra.Logger
}

// SetupRoutes configures all routes for the application
func SetupRoutes(
	router *gin.Engine,
	accountUseCase usecase.AccountUseCase,
	transactionUseCase usecase.TransactionUseCase,
	config RouterConfig,
) {
	// Initialize controllers
	accountController := NewAccountController(accountUseCase, config.Logger)
	transactionController := NewTransactionController(transactionUseCase, config.Logger)

	// Apply global middlewares
	router.Use(CORSMiddleware())
	router.Use(RequestIDMiddleware())
	router.Use(LoggingMiddleware(config.Logger))
	router.Use(RecoveryMiddleware(config.Logger))

	// Health check endpoint (no API key required)
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"service": "mini-bank-api",
		})
	})

	// API v1 routes with API key middleware
	v1 := router.Group("/api/v1")
	v1.Use(APIKeyMiddleware(config.APIKey, config.Logger))
	{
		// Account routes
		accounts := v1.Group("/accounts")
		{
			// Account-specific transaction routes
			accounts.GET("/:id/transactions", transactionController.GetTransactionsByAccount)

			accounts.POST("", accountController.CreateAccount)
			accounts.GET("", accountController.ListAccounts)
			accounts.GET("/:id", accountController.GetAccount)
			accounts.PUT("/:id", accountController.UpdateAccount)
			accounts.DELETE("/:id", accountController.DeleteAccount)
			accounts.PATCH("/:id/suspend", accountController.SuspendAccount)
			accounts.PATCH("/:id/activate", accountController.ActivateAccount)

		}

		// Transaction routes
		transactions := v1.Group("/transactions")
		{
			transactions.POST("", transactionController.CreateTransaction)
			transactions.GET("", transactionController.ListTransactions)
			transactions.GET("/:id", transactionController.GetTransaction)
			transactions.PATCH("/:id/confirm", transactionController.ConfirmTransaction)
			transactions.PATCH("/:id/cancel", transactionController.CancelTransaction)

			// Transaction status routes
			transactions.GET("/status/:status", transactionController.GetTransactionsByStatus)
		}
	}

	// Add a catch-all route for undefined endpoints
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"error":   "Not Found",
			"message": "The requested endpoint does not exist",
			"path":    ctx.Request.URL.Path,
		})
	})
}
