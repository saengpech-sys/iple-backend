package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

const (
	// SaltSize is the size of the salt in bytes
	SaltSize = 32
	// HashSize is the size of the hash in bytes
	HashSize = 64
	// ScryptN is the CPU/memory cost parameter for scrypt
	ScryptN = 32768
	// ScryptR is the block size parameter for scrypt
	ScryptR = 8
	// ScryptP is the parallelization parameter for scrypt
	ScryptP = 1
)

// HashPassword hashes a password using scrypt
func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, SaltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash the password using scrypt
	hash, err := scrypt.Key([]byte(password), salt, ScryptN, ScryptR, ScryptP, HashSize)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Combine salt and hash, then encode to base64
	combined := append(salt, hash...)
	encoded := base64.StdEncoding.EncodeToString(combined)

	return encoded, nil
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(password, hashedPassword string) (bool, error) {
	// Decode the base64 encoded hash
	combined, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Check if the combined hash has the correct length
	if len(combined) != SaltSize+HashSize {
		return false, errors.New("invalid hash format")
	}

	// Split salt and hash
	salt := combined[:SaltSize]
	hash := combined[SaltSize:]

	// Hash the provided password with the same salt
	newHash, err := scrypt.Key([]byte(password), salt, ScryptN, ScryptR, ScryptP, HashSize)
	if err != nil {
		return false, fmt.Errorf("failed to hash password for verification: %w", err)
	}

	// Compare hashes using constant-time comparison
	return subtle.ConstantTimeCompare(hash, newHash) == 1, nil
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}