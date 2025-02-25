package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yorukot/go-template/app/models"
	"github.com/yorukot/go-template/app/queries"
	"github.com/yorukot/go-template/pkg/encryption"
	"gorm.io/gorm"
)

var secret bool = false

// Generate new user access_token and refresh_token
func GenerateUserSession(c *gin.Context, userID uint64) error {
	var err error
	secretKey, err := encryption.RandStringRunes(1024, true)
	if err != nil {
		return err
	}
	session := models.Session{
		SessionID: encryption.GenerateID(),
		SecretKey: secretKey,
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(CookieRefreshTokenExpires)),
		CreatedAt: time.Now(),
	}

	// Check if a session with the same secretKey already exists
	for {
		_, result := queries.GetSessionQueueBySecretKey(session.SecretKey)
		if result.Error == gorm.ErrRecordNotFound {
			break
		} else if result.Error != nil {
			// Handle any other errors
			fmt.Println(result.Error.Error())
			return result.Error
		}

		secretKey, err = encryption.RandStringRunes(1024, true)
		if err != nil {
			return err
		}
		session.SecretKey = secretKey
	}

	// Create the new session in the database
	result := queries.CreateSessionQueue(session)
	if result.Error != nil {
		return result.Error
	}

	// Generate the access token
	accessTokenExpiresAt := time.Now().Add(time.Minute * time.Duration(CookieAccessTokenExpires))
	accessToken, err := encryption.GenerateNewJwtToken(userID, []string{}, accessTokenExpiresAt)
	if err != nil {
		return err
	}

	// Set the cookies
	c.SetCookie("refresh_token", session.SecretKey, CookieRefreshTokenExpires*24*60*60, "", "", secret, true)
	c.SetCookie("access_token", accessToken, CookieAccessTokenExpires*60, "/", "", secret, false)

	return nil
}

func init() {
	baseURL := os.Getenv("BASE_URL")

	if strings.HasPrefix(baseURL, "https://") {
		secret = true
	} else {
		secret = false
	}
}