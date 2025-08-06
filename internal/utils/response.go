package utils

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response structure
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Success bool                    `json:"success"`
	Error   string                  `json:"error"`
	Code    string                  `json:"code"`
	Message string                  `json:"message"`
	Details []ValidationErrorDetail `json:"details"`
}

// ValidationErrorDetail represents individual validation error
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(statusCode int, code, message string) ErrorResponse {
	return ErrorResponse{
		Success: false,
		Error:   "error",
		Code:    code,
		Message: message,
	}
}

// NewValidationErrorResponse creates a new validation error response
func NewValidationErrorResponse(details []ValidationErrorDetail) ValidationErrorResponse {
	return ValidationErrorResponse{
		Success: false,
		Error:   "validation_error",
		Code:    "VALIDATION_FAILED",
		Message: "Validation failed for one or more fields",
		Details: details,
	}
}

// FormatValidationErrors formats validator.ValidationErrors into our custom format
func FormatValidationErrors(errs validator.ValidationErrors) []ValidationErrorDetail {
	var details []ValidationErrorDetail
	
	for _, err := range errs {
		detail := ValidationErrorDetail{
			Field: err.Field(),
			Tag:   err.Tag(),
		}
		
		// Customize error messages based on validation tag
		switch err.Tag() {
		case "required":
			detail.Message = err.Field() + " is required"
		case "email":
			detail.Message = err.Field() + " must be a valid email address"
		case "min":
			detail.Message = err.Field() + " must be at least " + err.Param() + " characters long"
		case "max":
			detail.Message = err.Field() + " must be at most " + err.Param() + " characters long"
		case "oneof":
			detail.Message = err.Field() + " must be one of: " + err.Param()
		default:
			detail.Message = err.Field() + " is invalid"
		}
		
		details = append(details, detail)
	}
	
	return details
}

// SendSuccessResponse sends a success response
func SendSuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	response := NewSuccessResponse(message, data)
	return c.Status(statusCode).JSON(response)
}

// SendErrorResponse sends an error response
func SendErrorResponse(c *fiber.Ctx, statusCode int, code, message string) error {
	response := NewErrorResponse(statusCode, code, message)
	return c.Status(statusCode).JSON(response)
}

// SendValidationErrorResponse sends a validation error response
func SendValidationErrorResponse(c *fiber.Ctx, details []ValidationErrorDetail) error {
	response := NewValidationErrorResponse(details)
	return c.Status(fiber.StatusBadRequest).JSON(response)
}