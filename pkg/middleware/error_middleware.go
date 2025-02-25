package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yorukot/go-template/pkg/logger"
)

// ErrorLoggerMiddleware logs errors in all requests
func ErrorLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there was any error after request processing
		if len(c.Errors) > 0 {
			// Log each error that occurred during the request
			for _, err := range c.Errors {
				logger.LogError(c, err.Err, err.Error(), nil)
			}
		}
	}
}
