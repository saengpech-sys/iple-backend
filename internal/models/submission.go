package models

import "time"

// Submission represents a student's submission for an assignment or quiz (linked to a Lesson).
// Requirements reference: Core submission tracking[cite: 1], Linking to Enrollment [cite: 3]
type Submission struct {
	BaseModel // Includes ID, CreatedAt, UpdatedAt, DeletedAt

	LessonID     uint64    `gorm:"not null;index:idx_submission_lesson_user,unique,where:deleted_at IS NULL" json:"lesson_id"` // FK to Lesson (type assignment/quiz). Unique with UserID for non-deleted.
	UserID       uint64    `gorm:"not null;index:idx_submission_lesson_user,unique,where:deleted_at IS NULL" json:"user_id"`   // FK to User (student). Unique with LessonID for non-deleted.
	EnrollmentID uint64    `gorm:"not null;index" json:"enrollment_id"`                                                         // FK to Enrollment
	SubmittedAt  time.Time `gorm:"autoCreateTime;not null" json:"submitted_at"`                                                 // Use autoCreateTime
	// Stores student's answers, or references to uploaded files (e.g., GCS path)
	SubmissionData ContentData `gorm:"type:jsonb;not null;default:'{}'" json:"submission_data"`
	IsLate         bool        `gorm:"default:false;not null" json:"is_late"`

	// Relationships
	// Belongs To Lesson
	Lesson Lesson `gorm:"foreignKey:LessonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"lesson,omitempty"`

	// Belongs To User (Student)
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	// Belongs To Enrollment
	Enrollment Enrollment `gorm:"foreignKey:EnrollmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Exclude from default JSON

	// Has One Grade (A submission typically has only one primary grade)
	Grade *Grade `gorm:"foreignKey:SubmissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"grade,omitempty"` // Cascade delete grade if submission is deleted
}

// Note on uniqueIndex: idx_submission_lesson_user uses a partial unique index on (LessonID, UserID) for non-deleted records.
// This assumes a student submits only ONCE per lesson. If multiple submissions are allowed, remove the unique constraint. Check requirements.