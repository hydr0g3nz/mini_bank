package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hydr0g3nz/mini_bank/config"
	"github.com/hydr0g3nz/mini_bank/internal/adapter/controller"
	"github.com/hydr0g3nz/mini_bank/internal/adapter/repository/gorm/repository"
	usecase "github.com/hydr0g3nz/mini_bank/internal/application"
	infra "github.com/hydr0g3nz/mini_bank/internal/infrastructure"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg := config.LoadFromEnv()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatal("Configuration validation failed:", err)
	}

	// Initialize logger
	logger, err := infra.NewSimpleLogger(cfg.IsProduction())
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting Mini Bank API server",
		"environment", cfg.Server.Environment,
		"port", cfg.Server.Port,
	)

	// Connect to database using GORM
	// Connect to database
	db, err := infra.ConnectDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Run migrations
	if err := infra.MigrateDB(db); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	logger.Info("Database connected successfully")

	// Auto-migrate database tables (optional - you might want to use proper migrations)
	// if err := db.AutoMigrate(&model.Account{}, &model.Transaction{}); err != nil {
	// 	logger.Fatal("Failed to migrate database", "error", err)
	// }

	// Initialize Redis cache
	cache := infra.NewRedisClient(infra.CacheConfig{
		Host:     cfg.Cache.Host,
		Port:     cfg.Cache.Port,
		Password: cfg.Cache.Password,
		Db:       cfg.Cache.DB,
	})
	logger.Info("Redis cache connected successfully")

	// Initialize repositories
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	logger.Info("Repositories initialized")

	// Initialize use cases
	accountUseCase := usecase.NewAccountUseCase(accountRepo, cache, logger)
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, accountRepo, cache, logger)
	logger.Info("Use cases initialized")

	// Set Gin mode based on environment
	gin.SetMode(cfg.Server.Environment)

	// Initialize Gin router
	router := gin.New()

	// Setup routes
	routerConfig := controller.RouterConfig{
		APIKey: cfg.API.Key,
		Logger: logger,
	}

	controller.SetupRoutes(router, accountUseCase, transactionUseCase, routerConfig)
	logger.Info("Routes configured")

	// HTTP Server configuration
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting",
			"address", server.Addr,
			"environment", cfg.Server.Environment,
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	} else {
		logger.Info("Server shutdown completed")
	}

	// Close database connection
	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			logger.Error("Failed to close database connection", "error", err)
		} else {
			logger.Info("Database connection closed")
		}
	}

	// Close Redis connection
	if err := cache.Close(); err != nil {
		logger.Error("Failed to close Redis connection", "error", err)
	} else {
		logger.Info("Redis connection closed")
	}

	logger.Info("Server shutdown completed successfully")
}
