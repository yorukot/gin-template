package queries

import (
	"github.com/yorukot/go-template/app/models"
	db "github.com/yorukot/go-template/pkg/database"
	"gorm.io/gorm"
)

// Create new session
func CreateSessionQueue(session models.Session) *gorm.DB {
	// Create a new session record in the database
	result := db.GetDB().Create(&session)
	return result
}

// Get session by secretKey
func GetSessionQueueBySecretKey(secretKey string) (models.Session, *gorm.DB) {
	var session models.Session
	// Query the session by secretKey
	result := db.GetDB().Where("secret_key = ?", secretKey).First(&session)
	return session, result
}

// Delete session by secretKey
func DeleteSessionQueue(secretKey string) *gorm.DB {
	// Delete session by secretKey
	result := db.GetDB().Where("secret_key = ?", secretKey).Delete(&models.Session{})
	return result
}
