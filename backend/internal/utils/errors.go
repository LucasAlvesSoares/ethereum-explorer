package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// APIError represents a standardized API error response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e APIError) Error() string {
	return e.Message
}

// Common error types
var (
	ErrInvalidRequest   = APIError{Code: http.StatusBadRequest, Message: "Invalid request"}
	ErrNotFound         = APIError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrInternalServer   = APIError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrDatabaseError    = APIError{Code: http.StatusInternalServerError, Message: "Database error"}
	ErrBlockchainError  = APIError{Code: http.StatusServiceUnavailable, Message: "Blockchain service unavailable"}
	ErrValidationFailed = APIError{Code: http.StatusBadRequest, Message: "Validation failed"}
)

// NewAPIError creates a new API error with custom message
func NewAPIError(code int, message string, details ...string) APIError {
	err := APIError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

// NewValidationError creates a validation error with details
func NewValidationError(details string) APIError {
	return APIError{
		Code:    http.StatusBadRequest,
		Message: "Validation failed",
		Details: details,
	}
}

// NewDatabaseError creates a database error with logging
func NewDatabaseError(operation string, err error) APIError {
	logrus.Errorf("Database error during %s: %v", operation, err)
	return APIError{
		Code:    http.StatusInternalServerError,
		Message: "Database error",
		Details: fmt.Sprintf("Failed to %s", operation),
	}
}

// NewBlockchainError creates a blockchain service error with logging
func NewBlockchainError(operation string, err error) APIError {
	logrus.Errorf("Blockchain error during %s: %v", operation, err)
	return APIError{
		Code:    http.StatusServiceUnavailable,
		Message: "Blockchain service unavailable",
		Details: fmt.Sprintf("Failed to %s", operation),
	}
}

// HandleError sends a standardized error response
func HandleError(c *gin.Context, err error) {
	var apiErr APIError

	// Check if it's already an APIError
	if e, ok := err.(APIError); ok {
		apiErr = e
	} else {
		// Default to internal server error
		logrus.Errorf("Unhandled error: %v", err)
		apiErr = ErrInternalServer
	}

	c.JSON(apiErr.Code, gin.H{
		"error":   apiErr.Message,
		"details": apiErr.Details,
	})
}

// HandleValidationError sends a validation error response
func HandleValidationError(c *gin.Context, field, message string) {
	err := NewValidationError(fmt.Sprintf("%s: %s", field, message))
	HandleError(c, err)
}

// HandleDatabaseError sends a database error response
func HandleDatabaseError(c *gin.Context, operation string, err error) {
	apiErr := NewDatabaseError(operation, err)
	HandleError(c, apiErr)
}

// HandleBlockchainError sends a blockchain error response
func HandleBlockchainError(c *gin.Context, operation string, err error) {
	apiErr := NewBlockchainError(operation, err)
	HandleError(c, apiErr)
}

// HandleNotFound sends a not found error response
func HandleNotFound(c *gin.Context, resource string) {
	err := NewAPIError(http.StatusNotFound, fmt.Sprintf("%s not found", resource))
	HandleError(c, err)
}

// WrapError wraps an error with additional context
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}
