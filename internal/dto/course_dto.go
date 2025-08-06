package dto

import (
	"time"

	"github.com/saengepch-sys/iple-backend/internal/models"
)

// --- Course Request DTOs ---

type CreateCourseRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=200"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	IsPublished *bool  `json:"is_published,omitempty"` // Optional, defaults to false
}

type UpdateCourseRequest struct {
	Title       string `json:"title" validate:"omitempty,min=3,max=200"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	IsPublished *bool  `json:"is_published,omitempty"`
}

// --- Course Response DTOs ---

type CourseBasicResponse struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TeacherID   uint64    `json:"teacher_id"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CourseDetailResponse struct {
	CourseBasicResponse
	Teacher  UserBasicInfoResponse    `json:"teacher,omitempty"`
	Sections []SectionBasicResponse   `json:"sections,omitempty"`
	Stats    CourseStatsResponse      `json:"stats,omitempty"`
}

type CourseStatsResponse struct {
	TotalEnrollments   int `json:"total_enrollments"`
	ActiveEnrollments  int `json:"active_enrollments"`
	CompletedEnrollments int `json:"completed_enrollments"`
	TotalSections      int `json:"total_sections"`
	TotalLessons       int `json:"total_lessons"`
}

type SectionBasicResponse struct {
	ID          uint64    `json:"id"`
	CourseID    uint64    `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OrderIndex  int       `json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// --- Helper Functions ---

func BuildCourseBasicResponse(course *models.Course) CourseBasicResponse {
	return CourseBasicResponse{
		ID:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		TeacherID:   course.TeacherID,
		IsPublished: course.IsPublished,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}
}

func BuildCourseDetailResponse(course *models.Course) CourseDetailResponse {
	resp := CourseDetailResponse{
		CourseBasicResponse: BuildCourseBasicResponse(course),
	}
	
	// Add teacher info if preloaded
	if course.Teacher.ID != 0 {
		resp.Teacher = BuildUserBasicInfoResponse(&course.Teacher)
	}
	
	// Add sections if preloaded
	if course.Sections != nil {
		resp.Sections = BuildSectionBasicResponses(course.Sections)
	}
	
	return resp
}

func BuildCourseBasicResponses(courses []models.Course) []CourseBasicResponse {
	responses := make([]CourseBasicResponse, len(courses))
	for i, course := range courses {
		responses[i] = BuildCourseBasicResponse(&course)
	}
	return responses
}

func BuildSectionBasicResponse(section *models.Section) SectionBasicResponse {
	return SectionBasicResponse{
		ID:          section.ID,
		CourseID:    section.CourseID,
		Title:       section.Title,
		Description: section.Description,
		OrderIndex:  section.OrderIndex,
		CreatedAt:   section.CreatedAt,
		UpdatedAt:   section.UpdatedAt,
	}
}

func BuildSectionBasicResponses(sections []models.Section) []SectionBasicResponse {
	responses := make([]SectionBasicResponse, len(sections))
	for i, section := range sections {
		responses[i] = BuildSectionBasicResponse(&section)
	}
	return responses
}