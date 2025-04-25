package models

import "time"

// Grade represents a grade given for a specific submission or potentially a lesson directly.
// Requirements reference: Core grading functionality [cite: 1]
type Grade struct {
	BaseModel // Includes ID, CreatedAt, UpdatedAt, DeletedAt

	EnrollmentID uint64 `gorm:"not null;index" json:"enrollment_id"`   // FK to Enrollment
	LessonID     uint64 `gorm:"not null;index" json:"lesson_id"`       // FK to the specific Lesson being graded
	// Link to a specific submission (Optional, allows grading non-submission activities like participation)
	// UNIQUE constraint ensures one grade record per submission.
	SubmissionID *uint64 `gorm:"index;unique" json:"submission_id,omitempty"` // FK to Submission
	GraderID     uint64 `gorm:"not null;index" json:"grader_id"`       // FK to User (teacher or potentially an AI system user)
	// Using numeric for potentially more precision than float32
	Score    string `gorm:"type:numeric(10,2);not null" json:"score"`         // Actual score (e.g., 85.50) - Use string for precise decimal handling with DB
	MaxScore string `gorm:"type:numeric(10,2);not null" json:"max_score"`     // Max possible score (e.g., 100.00) - Use string
	Feedback string `gorm:"type:text" json:"feedback,omitempty"`              // Optional text feedback
	GradedAt time.Time `gorm:"autoCreateTime;not null" json:"graded_at"`      // Use autoCreateTime for initial grade, might be updated

	// Relationships
	// Belongs To Enrollment
	Enrollment Enrollment `gorm:"foreignKey:EnrollmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Delete grade if enrollment deleted

	// Belongs To Lesson
	Lesson Lesson `gorm:"foreignKey:LessonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Delete grade if lesson deleted

	// Belongs To Submission (Optional)
	Submission *Submission `gorm:"foreignKey:SubmissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Delete grade if submission deleted

	// Belongs To User (Grader)
	Grader User `gorm:"foreignKey:GraderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"grader,omitempty"` // Keep grade record but set grader to NULL if grader user is deleted
}

// Note on Score/MaxScore type: Using string with `type:numeric(10,2)` can sometimes help avoid floating-point inaccuracies
// when interacting with databases, especially across different libraries/languages.
// You'll need to handle conversion to/from float in your Service layer.
// Alternatively, use `float64` or a dedicated decimal type library if preferred, but ensure consistency.
// A CHECK constraint `CHECK (score <= max_score)` should be added in the SQL migration.
// The unique constraint on SubmissionID ensures only one grade row per submission.