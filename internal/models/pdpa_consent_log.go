package models

import "time"

// PDPAConsentLog records the history of consent actions for PDPA compliance and auditability.
// Requirements reference: PDPA compliance [cite: 3, 5]
type PDPAConsentLog struct {
	BaseModel // Includes ID, CreatedAt, UpdatedAt, DeletedAt (Soft delete might not be typical here, but BaseModel provides it)

	UserID uint64 `gorm:"not null;index:idx_pdpa_log_user_type_version" json:"user_id"` // FK to User

	// Type of consent (e.g., terms, privacy, specific AI feature). Use constants.
	ConsentType string `gorm:"size:100;not null;index:idx_pdpa_log_user_type_version" json:"consent_type"`

	// Action taken: GRANTED or WITHDRAWN. Use constants.
	ConsentAction string `gorm:"size:20;not null;index" json:"consent_action"`

	// Timestamp of the action. Not null, defaults to creation time.
	Timestamp time.Time `gorm:"autoCreateTime;not null" json:"timestamp"`

	// Version of the policy/terms consented to (e.g., '1.2', '2025-APR'). Not null.
	Version string `gorm:"size:50;not null;index:idx_pdpa_log_user_type_version" json:"version"`

	// Optional details about the consent action (e.g., specific scope granted/withdrawn).
	Details string `gorm:"type:text" json:"details,omitempty"`

	// Optional: Where the consent action took place (e.g., 'REGISTRATION', 'SETTINGS', 'FEATURE_PROMPT').
	Source string `gorm:"size:100;index" json:"source,omitempty"`

	// Relationship (exclude from JSON)
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Log deleted if user deleted
}

// Define constants for ConsentType and ConsentAction for consistency and type safety.
const (
	ConsentActionGranted   = "GRANTED"
	ConsentActionWithdrawn = "WITHDRAWN"

	ConsentTypeTerms             = "TERMS"             // General Terms of Service
	ConsentTypePrivacy           = "PRIVACY"           // Privacy Policy
	ConsentTypeAIProcessing      = "AI_PROCESSING"     // General consent for AI data use
	ConsentTypeAIRecommendation  = "AI_RECOMMENDATION" // Specific consent for personalized recommendations
	ConsentTypeAIAssessment      = "AI_ASSESSMENT"     // Specific consent for AI grading
	// Add other specific consent types as needed
)

// Composite index `idx_pdpa_log_user_type_version` on (user_id, consent_type, version) can be useful for querying latest consent status.