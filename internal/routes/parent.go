package routes

import (
	"github.com/saengepch-sys/iple-backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// SetupParentRoutes sets up parent-specific routes
func SetupParentRoutes(router fiber.Router) {
	// Children management
	router.Get("/children", getParentChildrenHandler)
	router.Post("/children", linkChildToParentHandler)
	router.Delete("/children/:id", unlinkChildFromParentHandler)
	
	// Children's progress monitoring
	router.Get("/children/:id/progress", getChildProgressHandler)
	router.Get("/children/:id/enrollments", getChildEnrollmentsHandler)
	router.Get("/children/:id/submissions", getChildSubmissionsHandler)
	
	// PDPA consent management
	router.Get("/children/:id/pdpa-consent", getChildPDPAConsentHandler)
	router.Post("/children/:id/pdpa-consent", updateChildPDPAConsentHandler)
}

// getParentChildrenHandler gets all children linked to the parent
func getParentChildrenHandler(c *fiber.Ctx) error {
	// TODO: Implement get parent children logic
	// 1. Get parent user ID from JWT token
	// 2. Fetch all children linked to this parent
	// 3. Return children list

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Get parent children endpoint - coming soon", nil)
}

// linkChildToParentHandler links a child to the current parent
func linkChildToParentHandler(c *fiber.Ctx) error {
	// TODO: Implement link child to parent logic
	// 1. Get parent user ID from JWT token
	// 2. Validate request body (child_user_id, verification_code)
	// 3. Verify the child exists and verification code is correct
	// 4. Create parent-child link
	// 5. Return success response

	return utils.SendSuccessResponse(c, fiber.StatusCreated, "Link child to parent endpoint - coming soon", nil)
}

// unlinkChildFromParentHandler unlinks a child from the current parent
func unlinkChildFromParentHandler(c *fiber.Ctx) error {
	childID := c.Params("id")
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Unlink child from parent endpoint - coming soon", fiber.Map{
		"child_id": childID,
	})
}

// getChildProgressHandler gets progress for a specific child
func getChildProgressHandler(c *fiber.Ctx) error {
	childID := c.Params("id")
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Get child progress endpoint - coming soon", fiber.Map{
		"child_id": childID,
	})
}

// getChildEnrollmentsHandler gets enrollments for a specific child
func getChildEnrollmentsHandler(c *fiber.Ctx) error {
	childID := c.Params("id")
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Get child enrollments endpoint - coming soon", fiber.Map{
		"child_id": childID,
	})
}

// getChildSubmissionsHandler gets submissions for a specific child
func getChildSubmissionsHandler(c *fiber.Ctx) error {
	childID := c.Params("id")
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Get child submissions endpoint - coming soon", fiber.Map{
		"child_id": childID,
	})
}

// getChildPDPAConsentHandler gets PDPA consent status for a child
func getChildPDPAConsentHandler(c *fiber.Ctx) error {
	childID := c.Params("id")
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Get child PDPA consent endpoint - coming soon", fiber.Map{
		"child_id": childID,
	})
}

// updateChildPDPAConsentHandler updates PDPA consent for a child
func updateChildPDPAConsentHandler(c *fiber.Ctx) error {
	childID := c.Params("id")
	return utils.SendSuccessResponse(c, fiber.StatusOK, "Update child PDPA consent endpoint - coming soon", fiber.Map{
		"child_id": childID,
	})
}