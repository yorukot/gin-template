package models

import (
	"time"

	db "github.com/yorukot/go-template/pkg/database"
)

func init() {
	db.GetDB().AutoMigrate(&Session{})
}

// Session type / table
type Session struct {
	SessionID uint64    `json:"session_id,string" gorm:"primaryKey"`
	SecretKey string    `json:"secret_key" gorm:"unique;not null;uniqueIndex"`
	UserAgent string    `json:"user_agent" gorm:"size:512"`
	UserID    uint64    `json:"user_id,string" gorm:"not null;index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

// Access token type
type AccessToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
