package models

// Section represents a section or module within a course.
type Section struct {
	BaseModel
	CourseID    uint64 `gorm:"not null;index;uniqueIndex:idx_section_order" json:"course_id"` // FK to Course, must exist. Added unique index with OrderIndex.
	Title       string `gorm:"size:255;not null" json:"title"`
	Description string `gorm:"type:text" json:"description,omitempty"`   // Optional description
	OrderIndex  int    `gorm:"default:0;not null;uniqueIndex:idx_section_order" json:"order_index"` // Ordering within the course. Added unique index with CourseID.
	IsPublished bool   `gorm:"default:true;not null;index" json:"is_published"`                   // Visibility, indexed

	// Relationships
	// Belongs To Course
	Course Course `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Use CASCADE

	// Has Many Lessons
	Lessons []Lesson `gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"lessons,omitempty"` // Cascade delete lessons if section is deleted
}

// Note on uniqueIndex:idx_section_order: This enforces that the combination of course_id and order_index is unique.
// Adjust if sections can have the same order index within a course.