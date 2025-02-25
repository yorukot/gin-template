package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yorukot/go-template/app/models"
	"github.com/yorukot/go-template/app/queries"
	"github.com/yorukot/go-template/pkg/encryption"
	"github.com/yorukot/go-template/pkg/utils"
	"gorm.io/gorm"
)

// EmailLoginRequest represents the request body for login
type EmailLoginRequest struct {
	Email    string `json:"email" binding:"email,max=320"`
	Password string `json:"password" binding:"required,max=128,min=8"`
}

// Login handles the user login process
func Login(c *gin.Context) {
	request, err := validateLoginRequest(c)
	if err != nil {
		return // Error response already sent in the validation function
	}

	user, err := fetchUserByEmail(c, request.Email)
	if err != nil {
		return // Error response already sent in the fetch function
	}

	if err := validateUserPassword(c, user, request.Password); err != nil {
		return // Error response already sent in the validation function
	}

	if err := generateUserSession(c, user.ID); err != nil {
		return // Error response already sent in the session function
	}

	utils.FullyResponse(c, 200, "Login successful", nil, nil)
}

// validateLoginRequest validates the incoming login request
func validateLoginRequest(c *gin.Context) (*EmailLoginRequest, error) {
	var request EmailLoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.FullyResponse(c, 400, "Invalid request", utils.ErrBadRequest, err.Error())
		return nil, err
	}

	if request.Email == "" {
		utils.FullyResponse(c, 400, "Email is required", utils.ErrBadRequest, nil)
		return nil, errors.New("email is required")
	}

	return &request, nil
}

// fetchUserByEmail retrieves the user by email from the database
func fetchUserByEmail(c *gin.Context, email string) (models.User, error) {
	user, result := queries.GetUserQueueByEmail(email)
	if result.Error == gorm.ErrRecordNotFound {
		utils.FullyResponse(c, 400, "Invalid email", utils.ErrInvalidUsernameOrEmail, nil)
		return models.User{}, result.Error
	} else if result.Error != nil {
		utils.ServerErrorResponse(c, 500, "Error check email", utils.ErrGetData, result.Error)
		return models.User{}, result.Error
	}

	return user, nil
}

// validateUserPassword verifies if the provided password matches the stored hash
func validateUserPassword(c *gin.Context, user models.User, password string) error {
	if user.Password == "" {
		utils.FullyResponse(c, 400, "Invalid password", utils.ErrInvalidPassword, nil)
		return errors.New("invalid password")
	}

	match, err := encryption.ComparePasswordAndHash(password, user.Password)
	if err != nil || !match {
		utils.FullyResponse(c, 400, "Invalid password", utils.ErrInvalidPassword, nil)
		return errors.New("invalid password")
	}

	return nil
}
