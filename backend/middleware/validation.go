package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents a validation error response
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Code    int               `json:"code"`
	Details map[string]string `json:"details,omitempty"`
}

// ValidationMiddleware provides enhanced validation with detailed error messages
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// FormatValidationError formats validator errors into user-friendly messages
func FormatValidationError(err error) ValidationErrorResponse {
	response := ValidationErrorResponse{
		Error:   "Validation failed",
		Code:    400,
		Details: make(map[string]string),
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := strings.ToLower(fieldError.Field())
			tag := fieldError.Tag()
			
			switch tag {
			case "required":
				response.Details[field] = "This field is required"
			case "min":
				response.Details[field] = "Value is too short or too small"
			case "max":
				response.Details[field] = "Value is too long or too large"
			case "email":
				response.Details[field] = "Invalid email format"
			case "oneof":
				response.Details[field] = "Invalid value. Must be one of the allowed values"
			case "alphanum":
				response.Details[field] = "Must contain only letters and numbers"
			case "len":
				response.Details[field] = "Invalid length"
			default:
				response.Details[field] = "Invalid value"
			}
		}
		response.Message = "One or more fields have validation errors"
	} else {
		response.Message = err.Error()
	}

	return response
}

// ContentTypeMiddleware ensures proper content type for API requests
func ContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, ValidationErrorResponse{
					Error:   "Unsupported Media Type",
					Message: "Content-Type must be application/json",
					Code:    415,
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// RequestSizeLimitMiddleware limits request size to prevent abuse
func RequestSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
} 