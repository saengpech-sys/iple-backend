package models

import "time"

// User represents a user in the system (student, teacher, admin, parent).
// Requirements reference: Core LMS user roles[cite: 1], PDPA considerations[cite: 3], Base user info [cite: 5]
type User struct {
	BaseModel
	ClerkUserID   *string    `gorm:"size:255;uniqueIndex" json:"clerk_user_id,omitempty"` // Optional Clerk ID, ensure unique index
	FirstName     string     `gorm:"size:100;not null" json:"first_name"`
	LastName      string     `gorm:"size:100;not null" json:"last_name"`
	Email         string     `gorm:"size:255;uniqueIndex;not null" json:"email"`           // Email must be unique and not null
	Password      string     `gorm:"type:varchar(255);not null" json:"-"`                  // Store hashed password, never expose
	Role          string     `gorm:"size:50;not null;index" json:"role"`                   // Role: student, teacher, admin, parent (Consider ENUM type in DB if possible)
	IsActive      bool       `gorm:"default:true;not null" json:"is_active"`               // User status, active by default
	LastLogin     *time.Time `gorm:"index" json:"last_login,omitempty"`                    // Track last login, index for potential queries

	// PDPA Consent Information [cite: 3, 5]
	PdpaConsentTimestamp *time.Time `json:"pdpa_consent_timestamp,omitempty"` // Timestamp of the latest general consent action

	// --- Relationships --- (omitempty helps reduce JSON size if not preloaded)

	// If Teacher: Courses they teach
	TaughtCourses []Course `gorm:"foreignKey:TeacherID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"taught_courses,omitempty"` // Prevent deleting teacher if they still have courses

	// If Student: Enrollments they are in
	Enrollments []Enrollment `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"enrollments,omitempty"` // Cascade delete enrollments if student is deleted

	// If Student: Submissions they made
	Submissions []Submission `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"submissions,omitempty"` // Cascade delete submissions if student is deleted

	// If Parent: Links to their children [cite: 1]
	ChildrenLinks []ParentChildLink `gorm:"foreignKey:ParentUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"children_links,omitempty"` // Cascade delete links if parent is deleted

	// If Student: Links to their parents
	ParentLinks []ParentChildLink `gorm:"foreignKey:ChildUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"parent_links,omitempty"` // Cascade delete links if child is deleted

	// If Grader: Grades they have given
	GivenGrades []Grade `gorm:"foreignKey:GraderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"given_grades,omitempty"` // Set GraderID to NULL if grader user is deleted
}

// --- User Methods (Examples) ---

// FullName returns the user's full name.
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// HasRole checks if the user has a specific role.
func (u *User) HasRole(roleName string) bool {
	return u.Role == roleName
}