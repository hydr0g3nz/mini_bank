package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hydr0g3nz/mini_bank/internal/application/dto"
	"github.com/hydr0g3nz/mini_bank/internal/domain/infra"
)

// APIKeyMiddleware creates a middleware that validates API key from x-api-key header
func APIKeyMiddleware(validAPIKey string, logger infra.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get API key from header
		apiKey := ctx.GetHeader("x-api-key")

		// Check if API key is provided
		if apiKey == "" {
			logger.Warn("API key missing in request",
				"path", ctx.Request.URL.Path,
				"method", ctx.Request.Method,
				"ip", ctx.ClientIP(),
			)

			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "MISSING_API_KEY",
				Message: "API key is required. Please provide x-api-key header",
			})
			ctx.Abort()
			return
		}

		// Validate API key
		if strings.TrimSpace(apiKey) != validAPIKey {
			logger.Warn("Invalid API key provided",
				"path", ctx.Request.URL.Path,
				"method", ctx.Request.Method,
				"ip", ctx.ClientIP(),
				"providedKey", apiKey[:min(len(apiKey), 8)]+"...", // Log only first 8 chars for security
			)

			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Code:    "INVALID_API_KEY",
				Message: "Invalid API key provided",
			})
			ctx.Abort()
			return
		}

		// Log successful authentication for monitoring
		logger.Debug("API key validated successfully",
			"path", ctx.Request.URL.Path,
			"method", ctx.Request.Method,
			"ip", ctx.ClientIP(),
		)

		// Continue to next handler
		ctx.Next()
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-api-key")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger infra.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log with structured format
		logger.Info("HTTP Request",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"ip", param.ClientIP,
			"userAgent", param.Request.UserAgent(),
			"bodySize", param.BodySize,
		)
		return ""
	})
}

// RecoveryMiddleware handles panics and recovers gracefully
func RecoveryMiddleware(logger infra.Logger) gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultWriter, func(ctx *gin.Context, recovered interface{}) {
		logger.Error("Panic recovered",
			"error", recovered,
			"path", ctx.Request.URL.Path,
			"method", ctx.Request.Method,
			"ip", ctx.ClientIP(),
		)

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Internal server error occurred",
		})
	})
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a simple request ID (in production, use a proper UUID library)
			requestID = generateRequestID()
		}

		ctx.Set("requestID", requestID)
		ctx.Header("X-Request-ID", requestID)
		ctx.Next()
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	// Simple implementation - in production use UUID or similar
	return "req_" + string(rune(1000000+(int(time.Now().UnixNano())%9000000)))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
