package models

import (
	"time"
)

// EventDraft represents draft data for event creation
type EventDraft struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// Draft data for each step (JSONB)
	GeneralDetailsDraft JSONB `gorm:"type:jsonb" json:"general_details_draft,omitempty"`
	MediaPromotionDraft JSONB `gorm:"type:jsonb" json:"media_promotion_draft,omitempty"`
	SpecialGuestsDraft  JSONB `gorm:"type:jsonb" json:"special_guests_draft,omitempty"`
	VolunteersDraft     JSONB `gorm:"type:jsonb" json:"volunteers_draft,omitempty"`

	// Optional: Link to event if draft is associated with an existing event
	EventID *uint `json:"event_id,omitempty"`

	// User email to track which user created the draft
	UserEmail string `gorm:"column:user_email" json:"user_email,omitempty"`

	CreatedOn time.Time  `json:"created_on,omitempty"`
	UpdatedOn *time.Time `json:"updated_on,omitempty"`
}

func (EventDraft) TableName() string {
	return "event_drafts"
}

