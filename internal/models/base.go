package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// BaseModel provides common fields for GORM models: ID, CreatedAt, UpdatedAt, DeletedAt (for soft deletes).
// Explicitly defined here to use uint64 for ID (GORM default `gorm.Model` uses uint).
// Requirements reference: General ORM best practices, Soft Delete pattern. [cite: 1]
type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `gorm:"index;autoCreateTime;not null" json:"created_at"` // Ensure not null
	UpdatedAt time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`       // Ensure not null
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                                  // Soft Delete, exclude from JSON responses
}

// ContentData defines a flexible structure for details stored in JSONB.
// Used for Lesson content, Submission answers, Quiz definitions, etc. [cite: 1]
type ContentData map[string]interface{}

// Value implements the driver.Valuer interface for ContentData.
// Ensures that even an empty or nil map is saved as a valid JSON object '{}'.
func (cd ContentData) Value() (driver.Value, error) {
	if cd == nil {
		// Marshal nil map as an empty JSON object
		return json.Marshal(map[string]interface{}{})
	}
	// Marshal the actual map (could be empty -> '{}')
	return json.Marshal(cd)
}

// Scan implements the sql.Scanner interface for ContentData.
// Handles reading JSONB data (including NULL or '{}') from the database.
func (cd *ContentData) Scan(value interface{}) error {
	source, ok := value.([]byte)
	if !ok {
		// Handle non-byte source (e.g., actual NULL from DB)
		if value == nil {
			*cd = nil // Represent DB NULL as nil map in Go
			return nil
		}
		return errors.New("type assertion to []byte failed for ContentData Scan")
	}

	// Handle empty byte slice (e.g., from empty JSON object '{}' in DB)
	if len(source) == 0 || string(source) == "{}" {
		*cd = make(ContentData) // Represent empty JSON object as empty map
		return nil
	}

	// Handle JSON null literal
	if string(source) == "null" {
		*cd = nil // Represent JSON null as nil map
		return nil
	}

	// Unmarshal actual JSON data
	// Allocate map if target pointer is nil before unmarshaling
	if *cd == nil {
		*cd = make(ContentData)
	}
	err := json.Unmarshal(source, cd)
	if err != nil {
		return fmt.Errorf("failed to unmarshal ContentData: %w", err)
	}

	return nil
}