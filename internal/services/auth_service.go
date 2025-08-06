package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/saengepch-sys/iple-backend/internal/config"
	"github.com/saengepch-sys/iple-backend/internal/dto"
	"github.com/saengepch-sys/iple-backend/internal/models"
	"github.com/saengepch-sys/iple-backend/internal/repositories"
	"github.com/saengepch-sys/iple-backend/internal/utils"
)

type AuthService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Validate role
	validRoles := []string{"student", "teacher", "admin", "parent"}
	isValidRole := false
	for _, role := range validRoles {
		if req.Role == role {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return nil, errors.New("invalid role")
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, repositories.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     strings.ToLower(req.Email),
		Password:  hashedPassword,
		Role:      req.Role,
		IsActive:  true,
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(fmt.Sprintf("%d", user.ID), user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Build response
	authResponse := dto.BuildAuthResponse(token, int64(config.Env.JWTExpiresIn.Seconds()), user)
	return &authResponse, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Verify password
	isValid, err := utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}
	if !isValid {
		return nil, errors.New("invalid email or password")
	}

	// Update last login timestamp
	now := time.Now()
	err = s.userRepo.UpdateUser(ctx, user.ID, map[string]interface{}{
		"last_login": now,
	})
	if err != nil {
		// Log error but don't fail login
		// slog.Warn("Failed to update last login timestamp", "user_id", user.ID, "error", err)
	} else {
		user.LastLogin = &now
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(fmt.Sprintf("%d", user.ID), user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Build response
	authResponse := dto.BuildAuthResponse(token, int64(config.Env.JWTExpiresIn.Seconds()), user)
	return &authResponse, nil
}

// RefreshToken generates a new token from an existing valid token
func (s *AuthService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Validate existing token
	claims, err := utils.ValidateJWT(req.Token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Find user to ensure they still exist and are active
	userID := claims.UserID
	user, err := s.userRepo.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	// Generate new token
	token, err := utils.GenerateJWT(userID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Build response
	authResponse := dto.BuildAuthResponse(token, int64(config.Env.JWTExpiresIn.Seconds()), user)
	return &authResponse, nil
}