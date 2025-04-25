package models

import "time"

// Enrollment represents a student's enrollment in a specific course.
// Requirements reference: Core enrollment tracking, progress[cite: 1], Unique enrollment [cite: 5]
type Enrollment struct {
	BaseModel // Includes ID, CreatedAt, UpdatedAt, DeletedAt

	UserID   uint64 `gorm:"not null;uniqueIndex:idx_user_course_unique_nondeleted,where:deleted_at IS NULL" json:"user_id"`   // Composite unique key with CourseID (only for non-deleted records)
	CourseID uint64 `gorm:"not null;uniqueIndex:idx_user_course_unique_nondeleted,where:deleted_at IS NULL" json:"course_id"` // Composite unique key with UserID

	EnrollmentDate time.Time `gorm:"autoCreateTime;not null" json:"enrollment_date"`        // Use autoCreateTime
	Status         string    `gorm:"size:50;default:'active';not null;index" json:"status"` // 'active', 'completed', 'withdrawn' (Consider DB ENUM type)
	CompletedAt    *time.Time `json:"completed_at,omitempty"`                               // Timestamp when completed
	Progress       float32   `gorm:"type:decimal(5,4);default:0.0;not null" json:"progress"` // Use decimal for precision (0.0000 to 1.0000), CHECK constraint in migration recommended

	// Relationships
	// Belongs To User (Student)
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"` // Enrollments deleted if student deleted

	// Belongs To Course
	Course Course `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"course,omitempty"` // Enrollments deleted if course deleted

	// Has Many Grades (Grades specifically for this enrollment)
	Grades []Grade `gorm:"foreignKey:EnrollmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"grades,omitempty"` // Grades deleted if enrollment deleted
}

// Note on uniqueIndex: idx_user_course_unique_nondeleted uses a partial index (where deleted_at IS NULL)
// This allows a user to re-enroll in a course after their previous enrollment was soft-deleted.
// If re-enrollment should be prevented entirely, remove the `where:deleted_at IS NULL` part. Check requirements.
// A CHECK constraint `CHECK (progress >= 0.0 AND progress <= 1.0)` should be added in the SQL migration.
// A CHECK constraint `CHECK (status IN ('active', 'completed', 'withdrawn'))` should be added in the SQL migration.