package user

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/yorukot/go-template/app/models"
	"github.com/yorukot/go-template/app/queries"
	"github.com/yorukot/go-template/pkg/utils"
	"gorm.io/gorm"
)

// UserProfile represents a user's public profile without sensitive information
type UserProfile struct {
	ID          uint64    `json:"id,string"`
	Avatar      *string   `json:"avatar,omitempty"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetProfile retrieves and returns the user's profile information
func GetProfile(c *gin.Context) {
	userID, err := extractUserIDFromContext(c)
	if err != nil {
		return // Error response sent in the extraction function
	}

	user, err := fetchUserByID(c, userID)
	if err != nil {
		return // Error response sent in the fetch function
	}

	profile, err := createUserProfile(c, user)
	if err != nil {
		return // Error response sent in the create function
	}

	utils.FullyResponse(c, 200, "User profile acquired", nil, profile)
}

// extractUserIDFromContext gets the user ID from the Gin context
func extractUserIDFromContext(c *gin.Context) (uint64, error) {
	jwtContextID, exists := c.Get("userID")
	if !exists {
		utils.FullyResponse(c, 403, "UserID not found in context", utils.ErrUserIDNotFound, nil)
		return 0, errors.New("userID not found in context")
	}

	return jwtContextID.(uint64), nil
}

// fetchUserByID retrieves user information from the database using the user ID
func fetchUserByID(c *gin.Context, userID uint64) (models.User, error) {
	user, result := queries.GetUserQueueByID(userID)
	if result.Error == gorm.ErrRecordNotFound {
		utils.FullyResponse(c, 403, "User not found", utils.ErrGetData, nil)
		return models.User{}, result.Error
	} else if result.Error != nil {
		utils.ServerErrorResponse(c, 500, "Error retrieving user data", utils.ErrGetData, result.Error)
		return models.User{}, result.Error
	}

	return user, nil
}

// createUserProfile creates a UserProfile object from a User model, excluding sensitive data
func createUserProfile(c *gin.Context, user models.User) (*UserProfile, error) {
	var profile UserProfile
	err := copier.Copy(&profile, &user)
	if err != nil {
		utils.ServerErrorResponse(c, 500, "Error processing user data", utils.ErrChangeType, err)
		return nil, err
	}

	return &profile, nil
}
