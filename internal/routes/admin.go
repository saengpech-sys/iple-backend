package routes

import (
	"github.com/saengepch-sys/iple-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes sets up admin-only routes
func SetupAdminRoutes(router fiber.Router) {
	// User management
	router.Get("/users", adminGetAllUsersHandler)
	router.Post("/users", adminCreateUserHandler)
	router.Put("/users/:id/status", adminUpdateUserStatusHandler)
	
	// Course management
	router.Get("/courses", adminGetAllCoursesHandler)
	router.Put("/courses/:id/status", adminUpdateCourseStatusHandler)
	
	// Enrollment management
	router.Get("/enrollments", adminGetAllEnrollmentsHandler)
	router.Post("/enrollments", adminCreateEnrollmentHandler)
	
	// System statistics
	router.Get("/stats", adminGetSystemStatsHandler)
}

// adminGetAllUsersHandler gets all users with admin privileges
func adminGetAllUsersHandler(c *fiber.Ctx) error {
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Admin - Get all users endpoint - coming soon", nil)
}

// adminCreateUserHandler creates a new user (admin only)
func adminCreateUserHandler(c *fiber.Ctx) error {
	return utils.SendSuccessResponse(c, fiber.StatusCreated, "Admin - Create user endpoint - coming soon", nil)
}

// adminUpdateUserStatusHandler updates user status (active/inactive)
func adminUpdateUserStatusHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Admin - Update user status endpoint - coming soon", fiber.Map{
		"user_id": userID,
	})
}

// adminGetAllCoursesHandler gets all courses with admin privileges
func adminGetAllCoursesHandler(c *fiber.Ctx) error {
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Admin - Get all courses endpoint - coming soon", nil)
}

// adminUpdateCourseStatusHandler updates course status
func adminUpdateCourseStatusHandler(c *fiber.Ctx) error {
	courseID := c.Params("id")
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Admin - Update course status endpoint - coming soon", fiber.Map{
		"course_id": courseID,
	})
}

// adminGetAllEnrollmentsHandler gets all enrollments with admin privileges
func adminGetAllEnrollmentsHandler(c *fiber.Ctx) error {
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Admin - Get all enrollments endpoint - coming soon", nil)
}

// adminCreateEnrollmentHandler creates an enrollment for any user (admin only)
func adminCreateEnrollmentHandler(c *fiber.Ctx) error {
	return utils.SendSuccessResponse(c, fiber.StatusCreated, "Admin - Create enrollment endpoint - coming soon", nil)
}

// adminGetSystemStatsHandler gets system statistics
func adminGetSystemStatsHandler(c *fiber.Ctx) error {
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Admin - System stats endpoint - coming soon", nil)
}