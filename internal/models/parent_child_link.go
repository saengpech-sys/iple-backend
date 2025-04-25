package models

// ParentChildLink links a parent user to a child (student) user.
// Business logic in the Service layer must validate User roles (parent, student) before creating links.
// Requirements reference: Parent Portal functionality [cite: 1]
type ParentChildLink struct {
	BaseModel // Includes ID, CreatedAt, UpdatedAt, DeletedAt

	// Composite unique index ensures a parent can only be linked to a child once (for non-deleted links)
	ParentUserID uint64 `gorm:"not null;uniqueIndex:idx_parent_child_unique_nondeleted,where:deleted_at IS NULL" json:"parent_user_id"` // FK to User (Role=parent)
	ChildUserID  uint64 `gorm:"not null;uniqueIndex:idx_parent_child_unique_nondeleted,where:deleted_at IS NULL" json:"child_user_id"`  // FK to User (Role=student)

	// Relationships (exclude from default JSON to prevent cycles and unnecessary data)
	// Belongs To Parent (User)
	Parent User `gorm:"foreignKey:ParentUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Cascade delete link if parent is deleted

	// Belongs To Child (User)
	Child User `gorm:"foreignKey:ChildUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Cascade delete link if child is deleted
}