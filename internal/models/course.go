package models

// Course represents a course offered in the LMS.
// Requirements reference: Core course structure [cite: 1, 5]
type Course struct {
	BaseModel
	Title       string `gorm:"size:255;not null" json:"title"`
	Description string `gorm:"type:text" json:"description,omitempty"` // Allow empty description
	TeacherID   uint64 `gorm:"not null;index" json:"teacher_id"`       // Course must have a teacher
	Subject     string `gorm:"size:100;index" json:"subject,omitempty"`
	GradeLevel  string `gorm:"size:50;index" json:"grade_level,omitempty"` // e.g., 'M1', 'M6'
	IsPublished bool   `gorm:"default:false;not null;index" json:"is_published"` // Index for filtering published courses

	// Relationships
	// Belongs To a Teacher (User)
	Teacher User `gorm:"foreignKey:TeacherID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"teacher,omitempty"` // Ensure Teacher exists

	// Has Many Sections
	Sections []Section `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sections,omitempty"` // Cascade delete sections if course is deleted

	// Has Many Enrollments
	Enrollments []Enrollment `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"enrollments,omitempty"` // Cascade delete enrollments if course is deleted
}