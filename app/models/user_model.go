package models

import (
	"time"

	db "github.com/yorukot/go-template/pkg/database"
)

func init() {
	db.GetDB().AutoMigrate(&User{})
}

// Users data type / table
type User struct {
	ID          uint64    `json:"id,string" gorm:"primaryKey" binding:"required"`
	Avatar      *string   `json:"avatar,omitempty"`
	DisplayName string    `json:"display_name" binding:"required"`
	Email       string    `json:"email" gorm:"unique" binding:"required,email"` // Unique
	Password    string    `json:"password,omitempty"`                           // Hashed password
	CreatedAt   time.Time `json:"created_at" gorm:"autoUpdateTime" binding:"required"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoCreateTime" binding:"required"`
}
