package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yorukot/go-template/pkg/logger"
	"go.uber.org/zap"
)

// CustomLogger is a middleware to log HTTP requests with execution time and status code
func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate execution time
		latency := time.Since(startTime)

		// Get request status code
		statusCode := c.Writer.Status()

		// Log the request details
		logger.Log.Info("HTTP Request",
			zap.String("timestamp", time.Now().Format("2006/01/02 - 15:04:05")),
			zap.Int("status_code", statusCode),
			zap.String("method", c.Request.Method),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("path", c.Request.URL.Path),
		)
	}
}
