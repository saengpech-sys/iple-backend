package models

// Import time if using EstimatedDurationMinutes

// Lesson represents a single learning unit (text, video, quiz, assignment, document) within a section.
// Requirements reference: Core lesson structure[cite: 1], Content flexibility [cite: 3, 5]
type Lesson struct {
	BaseModel
	SectionID   uint64    `gorm:"not null;index;uniqueIndex:idx_lesson_order" json:"section_id"` // FK to Section, must exist. Added unique index with OrderIndex.
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description,omitempty"`       // Optional description
	ContentType string    `gorm:"size:50;not null;index" json:"content_type"`   // 'text', 'video', 'quiz', 'assignment', 'document' (Consider ENUM in DB)
	ContentData ContentData `gorm:"type:jsonb;not null;default:'{}'" json:"content_data"` // Use jsonb, ensure not null, default empty object
	OrderIndex  int       `gorm:"default:0;not null;uniqueIndex:idx_lesson_order" json:"order_index"` // Ordering within the section. Added unique index with SectionID.
	IsPublished bool      `gorm:"default:true;not null;index" json:"is_published"`                  // Visibility, indexed
	EstimatedDurationMinutes *int `json:"estimated_duration_minutes,omitempty"`                    // Optional: Estimated time

	// Relationships
	// Belongs To Section
	Section Section `gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Use CASCADE

	// Has Many Submissions (Only relevant if ContentType is 'assignment' or 'quiz')
	Submissions []Submission `gorm:"foreignKey:LessonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"submissions,omitempty"` // Cascade delete submissions if lesson is deleted

    // Has Many Grades (If grading lesson directly without submission)
    DirectGrades []Grade `gorm:"foreignKey:LessonID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"direct_grades,omitempty"` // Grades directly linked to this lesson
}

// Note on uniqueIndex:idx_lesson_order: This enforces that the combination of section_id and order_index is unique.