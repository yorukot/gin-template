package encryption

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yorukot/go-template/pkg/logger"
)

var JwtSecretKey = ""

func init() {
	JwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	if JwtSecretKey == "" {
		logger.Log.Fatal("missing JWT secret key")
	}
}



// Generate new jwt token with credentials
func GenerateNewJwtToken(id uint64, credentials []string, expiresAt time.Time) (string, error) {
	// Create a new claims.
	claims := jwt.MapClaims{}

	// Set public claims:
	claims["sub"] = id
	claims["exp"] = expiresAt.Unix()

	// Set private token credentials:
	for _, credential := range credentials {
		claims[credential] = true
	}

	// Create a new JWT access token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token.
	t, err := token.SignedString([]byte(JwtSecretKey))
	if err != nil {
		// Return error, it JWT token generation failed.
		return "", err
	}

	return t, nil
}

// ParseAndValidateJWT parses and validates the JWT, returning the claims and any error
func ParseAndValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Retrieve the JWT secret key from environment variables

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(JwtSecretKey), nil
	})

	// Validate token and check for errors
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Ensure the claims are of type jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify the token expiration time
	expiresAtFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expiration time datatype")
	}
	expiresAt := int64(expiresAtFloat)
	if time.Now().Unix() >= expiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}