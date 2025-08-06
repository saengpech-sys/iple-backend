package dto

import (
	"time"

	"github.com/saengepch-sys/iple-backend/internal/models"
)

// --- Authentication Request DTOs ---

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Role      string `json:"role" validate:"required,oneof=student teacher admin parent"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// --- User Request DTOs ---

type UpdateUserProfileRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  string `json:"last_name" validate:"omitempty,min=2,max=100"`
	Email     string `json:"email" validate:"omitempty,email"`
}

type UpdateUserPasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	Role      string `json:"role" validate:"required,oneof=student teacher admin parent"`
}

// --- Response DTOs ---

type AuthResponse struct {
	Token     string                `json:"token"`
	ExpiresIn int64                 `json:"expires_in"` // seconds
	User      UserBasicInfoResponse `json:"user"`
}

type UserBasicInfoResponse struct {
	ID        uint64     `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

type UserDetailResponse struct {
	UserBasicInfoResponse
	PdpaConsentTimestamp *time.Time `json:"pdpa_consent_timestamp,omitempty"`
	// Add other detailed fields as needed
}

// --- Helper Functions ---

func BuildUserBasicInfoResponse(user *models.User) UserBasicInfoResponse {
	return UserBasicInfoResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		LastLogin: user.LastLogin,
	}
}

func BuildUserDetailResponse(user *models.User) UserDetailResponse {
	return UserDetailResponse{
		UserBasicInfoResponse: BuildUserBasicInfoResponse(user),
		PdpaConsentTimestamp:  user.PdpaConsentTimestamp,
	}
}

func BuildUserBasicInfoResponses(users []models.User) []UserBasicInfoResponse {
	responses := make([]UserBasicInfoResponse, len(users))
	for i, user := range users {
		responses[i] = BuildUserBasicInfoResponse(&user)
	}
	return responses
}

func BuildAuthResponse(token string, expiresIn int64, user *models.User) AuthResponse {
	return AuthResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		User:      BuildUserBasicInfoResponse(user),
	}
}