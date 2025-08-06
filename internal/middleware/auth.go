package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthRequired middleware validates JWT tokens
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format. Use 'Bearer <token>'",
			})
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token required",
			})
		}

		// TODO: Validate JWT token using utils.ValidateJWT
		// For now, just store the token in context
		c.Locals("token", token)
		c.Locals("user_id", "placeholder_user_id") // Will be replaced with actual user ID from JWT

		return c.Next()
	}
}

// OptionalAuth middleware validates JWT tokens if present, but doesn't require them
func OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token != "" {
				// TODO: Validate JWT token if present
				c.Locals("token", token)
				c.Locals("user_id", "placeholder_user_id")
			}
		}
		
		return c.Next()
	}
}

// RequireRole middleware checks if the authenticated user has the required role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Get user role from JWT token or database
		userRole := "student" // Placeholder - will be replaced with actual role checking
		
		if userRole != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}
		
		return c.Next()
	}
}