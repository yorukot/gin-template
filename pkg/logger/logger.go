package logger

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Log *zap.Logger

// Initialize the logger during package initialization
func init() {
	// Create the logger instance based on environment
	Log = zap.Must(zap.NewProduction())
	if os.Getenv("GIN_MODE") == "debug" {
		Log = zap.Must(zap.NewDevelopment())
	}
}

// LogError handles error logging with context
func LogError(c *gin.Context, err error, message string, extraFields map[string]interface{}) {
	// Log the error with context information
	fields := []zap.Field{
		zap.String("error", fmt.Sprintf("%v", err)),
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
	}

	// Add extra fields if provided
	for key, value := range extraFields {
		fields = append(fields, zap.Any(key, value))
	}

	// Log the error
	Log.Error(message, fields...)
}