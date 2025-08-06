package routes

import (
	"github.com/saengepch-sys/iple-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all application routes
func SetupRoutes(app *fiber.App) {
	// API version prefix
	api := app.Group("/api/v1")

	// Health check endpoint (no auth required)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "IPLE Backend API is running",
			"version": "1.0.0",
		})
	})

	// Auth routes (no auth required)
	auth := api.Group("/auth")
	SetupAuthRoutes(auth)

	// Protected routes (require authentication)
	protected := api.Group("/", middleware.AuthRequired())

	// User management routes
	users := protected.Group("/users")
	SetupUserRoutes(users)

	// Course management routes
	courses := protected.Group("/courses")
	SetupCourseRoutes(courses)

	// Enrollment routes
	enrollments := protected.Group("/enrollments")
	SetupEnrollmentRoutes(enrollments)

	// Admin routes (require admin role)
	admin := protected.Group("/admin", middleware.RequireRole("admin"))
	SetupAdminRoutes(admin)

	// Parent routes (require parent role)
	parent := protected.Group("/parent", middleware.RequireRole("parent"))
	SetupParentRoutes(parent)
}