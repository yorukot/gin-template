package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yorukot/go-template/pkg/encryption"
	"github.com/yorukot/go-template/pkg/utils"
)

// IsAuthorized is a middleware to check if the user is authorized
func IsAuthorized() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the JWT token from the cookie
		cookie, err := c.Request.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			utils.FullyResponse(c, 403, "Authorization token is empty.", "authentication_key_not_found", nil)
			c.Abort()
			return
		}

		// Parse and validate the JWT token
		claims, err := encryption.ParseAndValidateJWT(cookie.Value)
		if err != nil {
			utils.FullyResponse(c, 403, err.Error(), utils.ErrUnauthorized, nil)
			c.Abort()
			return
		}

		// Retrieve the user ID (subject) from the claims
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			utils.FullyResponse(c, 403, "UserID error", utils.ErrUnauthorized, nil)
			c.Abort()
			return
		}
		userID := uint64(userIDFloat)

		_, ok = claims["pedding"].(bool)
		if ok {
			utils.FullyResponse(c, 403, "Please verify email first", utils.ErrUnauthorized, nil)
			c.Abort()
			return
		}

		// Add the user ID to the request context for further use
		c.Set("userID", userID)
		c.Next()
	}
}

// IsAuthorized is a middleware to check if the user is authorized
func IsPeddingVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the JWT token from the cookie
		cookie, err := c.Request.Cookie("verify_pedding")
		if err != nil || cookie.Value == "" {
			c.Next()
			return
		}

		// Parse and validate the JWT token
		claims, err := encryption.ParseAndValidateJWT(cookie.Value)
		if err != nil {
			c.Next()
			return
		}

		// Retrieve the user ID (subject) from the claims
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			c.Next()
			return
		}
		userID := uint64(userIDFloat)

		pedding, ok := claims["pedding"].(bool)
		if !ok {
			c.Next()
			return
		}
		if !pedding {
			c.Next()
			return
		}

		// Add the user ID to the request context for further use
		c.Set("userID", userID)
		c.Next()
	}
}
