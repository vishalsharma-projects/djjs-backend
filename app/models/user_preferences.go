package models

import (
	"time"
)

// UserPreferences represents user-specific preferences stored in the database
type UserPreferences struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// User email to track which user owns these preferences
	UserEmail string `gorm:"column:user_email;not null;index" json:"user_email"`

	// Preference type (e.g., "events_list_columns", "dashboard_layout", etc.)
	PreferenceType string `gorm:"column:preference_type;not null;index" json:"preference_type"`

	// Preference data stored as JSONB for flexibility
	PreferenceData JSONB `gorm:"type:jsonb" json:"preference_data"`

	CreatedOn time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
}

func (UserPreferences) TableName() string {
	return "user_preferences"
}

// EventsListColumnPreferences represents the structure for events list column preferences
type EventsListColumnPreferences struct {
	HiddenColumns []string `json:"hidden_columns"`
	PinnedColumns []string `json:"pinned_columns"`
	ColumnOrder   []string `json:"column_order"`
}

