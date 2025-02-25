// pkg/utils/context.go
package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext retrieves the user ID from the request context.
func GetUserIDFromContext(c *gin.Context) (uint64, error) {
	ContextUserID, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("userID not found in context")
	}
	return ContextUserID.(uint64), nil
}
