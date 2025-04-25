package repositories

import (
	"context"
	"iple-backend/internal/models"
)

// UserListOptions defines options for listing users (pagination, filtering, sorting).
type UserListOptions struct {
	Limit  int
	Offset int
	// Add other filtering/sorting fields as needed, e.g., Role, IsActive, SortBy
	Role     *string // Use pointer to distinguish between no filter and filter by empty string
	IsActive *bool
	// ... other filters
}

// UserGetOptions defines options for getting a user (e.g., preloading).
type UserGetOptions struct {
	PreloadEnrollments bool
	PreloadCourses     bool // If teacher
	PreloadLinks       bool // Parent/Child links
	// Add other preload flags as needed
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	// GetUserByID retrieves a user by ID with optional preloading.
	GetUserByID(ctx context.Context, id uint64, opts *UserGetOptions) (*models.User, error)
	// ListUsers retrieves a list of users with pagination and filtering options.
	ListUsers(ctx context.Context, opts UserListOptions) ([]models.User, int64, error) // Returns users, total count, error
	// UpdateUser updates specific fields of an existing user.
	// Pass only the fields to update in the 'user' parameter map.
	UpdateUser(ctx context.Context, id uint64, updates map[string]interface{}) error
	// DeleteUser performs a soft delete on a user.
	DeleteUser(ctx context.Context, id uint64) error
	// FindOrCreateUser (Optional, useful sometimes)
	// ... other specific query methods if needed
}
