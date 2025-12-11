package models

import "time"

// Volunteer represents volunteer details captured from UI
// swagger:model Volunteer
type Volunteer struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchID      uint       `gorm:"not null" json:"branch_id" validate:"required,min=1"`
	Branch        Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	VolunteerName string     `gorm:"not null" json:"volunteer_name" validate:"required,min=2,max=255"`
	Contact       string     `gorm:"column:contact" json:"contact,omitempty" validate:"omitempty,max=20"`
	NumberOfDays  int        `gorm:"column:number_of_days" json:"number_of_days,omitempty" validate:"omitempty,min=0,max=365"`
	SevaInvolved  string     `json:"seva_involved,omitempty" validate:"omitempty,min=2,max=500"`
	MentionSeva   string     `gorm:"column:mention_seva" json:"mention_seva,omitempty" validate:"omitempty,min=2,max=500"`
	EventID       uint       `json:"event_id" validate:"required,min=1"`
	Event         Event      `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
	CreatedOn     time.Time  `json:"created_on,omitempty"`
	UpdatedOn     *time.Time `json:"updated_on,omitempty"`
	CreatedBy     string     `json:"created_by,omitempty"`
	UpdatedBy     string     `json:"updated_by,omitempty"`
}
