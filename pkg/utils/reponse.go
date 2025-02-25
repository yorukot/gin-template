package utils

import (
	"github.com/gin-gonic/gin"
)

// ServerErrorResponse is a helper function to handle errors and send responses
func ServerErrorResponse(c *gin.Context, statusCode int, message string, errorCode string, err error) {
	if err != nil {
		c.Error(err)
	}
	FullyResponse(c, statusCode, message, errorCode, nil)
}

// Fully response
func FullyResponse(c *gin.Context, statusCode int, message string, errorCode interface{}, data interface{}) {
	// Prepare the response payload
	response := gin.H{"message": message}
	if data != nil {
		response["result"] = data
	}

	// If it's an error, include an "error" field
	var errorCodePtr *string

	switch v := errorCode.(type) {
	case string:
		errorCodePtr = &v
	case *string:
		errorCodePtr = v
	case nil:
		errorCodePtr = nil
	default:
		panic("invalid errorCode type")
	}

	if errorCodePtr != nil {
		response["error"] = *errorCodePtr
	}

	// Send JSON response
	c.JSON(statusCode, response)
}
