package middleware

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

// StructuredLogger creates a structured logger middleware using slog
func StructuredLogger(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Process request
		err := c.Next()
		
		// Calculate latency
		latency := time.Since(start)
		
		// Get request ID from context
		requestID, _ := c.Locals(RequestIDKey).(string)
		
		// Log the request
		logger.Info("HTTP Request",
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency", latency.String(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
			"request_id", requestID,
		)
		
		return err
	}
}