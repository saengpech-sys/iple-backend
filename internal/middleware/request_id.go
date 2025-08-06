package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

// RequestID middleware adds a unique request ID to each request
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if request ID is already set (e.g., from load balancer)
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			// Generate a new UUID for the request
			requestID = uuid.New().String()
		}
		
		// Store in context locals for use by other middleware/handlers
		c.Locals(RequestIDKey, requestID)
		
		// Add to response headers for client debugging
		c.Set("X-Request-ID", requestID)
		
		return c.Next()
	}
}