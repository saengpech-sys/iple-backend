package utils

import (
	"errors"
	"time"

	"github.com/saengepch-sys/iple-backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey []byte

// InitJWT initializes the JWT secret key from config
func InitJWT() error {
	if config.Env == nil {
		return errors.New("config not loaded")
	}
	
	if config.Env.JWTSecretKey == "" {
		return errors.New("JWT secret key not configured")
	}
	
	jwtSecretKey = []byte(config.Env.JWTSecretKey)
	return nil
}

// Claims represents the JWT claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for a user
func GenerateJWT(userID, email, role string) (string, error) {
	if len(jwtSecretKey) == 0 {
		return "", errors.New("JWT secret key not initialized")
	}

	// Create claims
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Env.JWTExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "iple-backend",
			Subject:   userID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign token with secret
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*Claims, error) {
	if len(jwtSecretKey) == 0 {
		return nil, errors.New("JWT secret key not initialized")
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshJWT creates a new token for an existing valid token (extend expiry)
func RefreshJWT(tokenString string) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same claims but extended expiry
	return GenerateJWT(claims.UserID, claims.Email, claims.Role)
}