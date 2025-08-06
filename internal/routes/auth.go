package routes

import (
	"errors"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/saengepch-sys/iple-backend/internal/dto"
	"github.com/saengepch-sys/iple-backend/internal/services"
	"github.com/saengepch-sys/iple-backend/internal/utils"
)

var authService *services.AuthService

// SetupAuthRoutes sets up authentication related routes
func SetupAuthRoutes(router fiber.Router) {
	router.Post("/register", registerHandler)
	router.Post("/login", loginHandler)
	router.Post("/refresh", refreshTokenHandler)
	router.Post("/logout", logoutHandler)
}

// SetAuthService sets the auth service for route handlers
func SetAuthService(service *services.AuthService) {
	authService = service
}

// registerHandler handles user registration
func registerHandler(c *fiber.Ctx) error {
	if authService == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "AUTH_SERVICE_ERROR", "Authentication service not initialized")
	}

	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			details := utils.FormatValidationErrors(validationErrors)
			return utils.SendValidationErrorResponse(c, details)
		}
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	}

	// Call service
	response, err := authService.Register(c.Context(), req)
	if err != nil {
		// Handle specific errors
		if err.Error() == "user with this email already exists" {
			return utils.SendErrorResponse(c, fiber.StatusConflict, "EMAIL_EXISTS", "User with this email already exists")
		}
		if err.Error() == "invalid role" {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "INVALID_ROLE", "Invalid role specified")
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "REGISTRATION_FAILED", "Registration failed")
	}

	return utils.SendSuccessResponse(c, fiber.StatusCreated, "Registration successful", response)
}

// loginHandler handles user login
func loginHandler(c *fiber.Ctx) error {
	if authService == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "AUTH_SERVICE_ERROR", "Authentication service not initialized")
	}

	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			details := utils.FormatValidationErrors(validationErrors)
			return utils.SendValidationErrorResponse(c, details)
		}
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	}

	// Call service
	response, err := authService.Login(c.Context(), req)
	if err != nil {
		// Handle specific errors
		if err.Error() == "invalid email or password" {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		}
		if err.Error() == "account is disabled" {
			return utils.SendErrorResponse(c, fiber.StatusForbidden, "ACCOUNT_DISABLED", "Account is disabled")
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "LOGIN_FAILED", "Login failed")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Login successful", response)
}

// refreshTokenHandler handles token refresh
func refreshTokenHandler(c *fiber.Ctx) error {
	if authService == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "AUTH_SERVICE_ERROR", "Authentication service not initialized")
	}

	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "INVALID_REQUEST_BODY", "Invalid request body")
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			details := utils.FormatValidationErrors(validationErrors)
			return utils.SendValidationErrorResponse(c, details)
		}
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	}

	// Call service
	response, err := authService.RefreshToken(c.Context(), req)
	if err != nil {
		// Handle specific errors
		if err.Error() == "invalid token" {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "INVALID_TOKEN", "Invalid token")
		}
		if err.Error() == "user not found" {
			return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "USER_NOT_FOUND", "User not found")
		}
		if err.Error() == "account is disabled" {
			return utils.SendErrorResponse(c, fiber.StatusForbidden, "ACCOUNT_DISABLED", "Account is disabled")
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "TOKEN_REFRESH_FAILED", "Token refresh failed")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Token refresh successful", response)
}

// logoutHandler handles user logout
func logoutHandler(c *fiber.Ctx) error {
	// For JWT tokens, logout is typically handled client-side by removing the token
	// Optionally, implement token blacklisting for enhanced security
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Logout successful", nil)
}