package dto

import (
	"time"

	"github.com/saengepch-sys/iple-backend/internal/models"
)

// --- Request DTOs ---

// No specific request DTO needed for GET enrollments (uses query params)
// POST /courses/{id}/enroll doesn't need a body if student ID comes from auth token

// (Optional) If updating status via API needs specific fields
// type UpdateEnrollmentStatusRequest struct {
// 	Status string `json:"status" validate:"required,oneof=active completed withdrawn"`
// }

// --- Response DTOs ---

// Basic enrollment info for lists
type EnrollmentBasicResponse struct {
	ID             uint64    `json:"id"`
	UserID         uint64    `json:"user_id"` // Student ID
	CourseID       uint64    `json:"course_id"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Status         string    `json:"status"`
	Progress       float32   `json:"progress"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// Detailed enrollment info (includes related data)
type EnrollmentDetailResponse struct {
	EnrollmentBasicResponse             // Embed basic info
	Student                 UserBasicInfoResponse `json:"student,omitempty"` // Basic student info
	Course                  CourseBasicResponse   `json:"course,omitempty"`  // Basic course info
	// Grades                  []GradeBasicResponse `json:"grades,omitempty"`  // Optional: List of grades for this enrollment
	// Add other relevant details
}

// --- Helper Functions ---

func BuildEnrollmentBasicResponse(enrl *models.Enrollment) EnrollmentBasicResponse {
	return EnrollmentBasicResponse{
		ID:             enrl.ID,
		UserID:         enrl.UserID,
		CourseID:       enrl.CourseID,
		EnrollmentDate: enrl.EnrollmentDate,
		Status:         enrl.Status,
		Progress:       enrl.Progress,
		CompletedAt:    enrl.CompletedAt,
	}
}

func BuildEnrollmentDetailResponse(enrl *models.Enrollment) EnrollmentDetailResponse {
	resp := EnrollmentDetailResponse{
		EnrollmentBasicResponse: BuildEnrollmentBasicResponse(enrl),
	}
	if enrl.User.ID != 0 {
		resp.Student = BuildUserBasicInfoResponse(&enrl.User)
	}
	if enrl.Course.ID != 0 {
		resp.Course = BuildCourseBasicResponse(&enrl.Course) // Assuming this helper exists in course_dto.go
	}
    // TODO: Build and add grades if preloaded
	// if enrl.Grades != nil {
	// 	resp.Grades = BuildGradeBasicResponses(enrl.Grades)
	// }
	return resp
}

// Helper to build list of basic responses
func BuildEnrollmentBasicResponses(enrollments []models.Enrollment) []EnrollmentBasicResponse {
	responses := make([]EnrollmentBasicResponse, len(enrollments))
	for i, enrl := range enrollments {
		responses[i] = BuildEnrollmentBasicResponse(&enrl)
	}
	return responses
}