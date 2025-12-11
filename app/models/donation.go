package models

import "time"

// Donation represents donation details for an event
type Donation struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	EventID  uint `gorm:"not null" json:"event_id"`
	BranchID uint `gorm:"not null" json:"branch_id"`

	DonationType string  `json:"donation_type,omitempty"`
	Amount       float64 `json:"amount,omitempty"`
	KindType     string  `json:"kindtype,omitempty"`

	CreatedOn time.Time `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn time.Time `gorm:"autoUpdateTime" json:"updated_on"`

	CreatedBy string `json:"created_by,omitempty" gorm:"<-:create"` // only set on create
	UpdatedBy string `json:"updated_by,omitempty"`

	// Relations
	Event  Event  `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
	Branch Branch `gorm:"foreignKey:BranchID;references:ID" json:"branch,omitempty"`
}
