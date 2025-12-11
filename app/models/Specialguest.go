package models

import (
	"time"
)

// SpecialGuest represents a special guest in the system
// swagger:model SpecialGuest
type SpecialGuest struct {
	ID                   uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Gender               string     `json:"gender,omitempty"`
	Prefix               string     `gorm:"not null" json:"prefix"`
	FirstName            string     `json:"first_name,omitempty"`
	MiddleName           string     `json:"middle_name,omitempty"`
	LastName             string     `json:"last_name,omitempty"`
	Designation          string     `json:"designation,omitempty"`
	Organization         string     `json:"organization,omitempty"`
	Email                string     `gorm:"unique" json:"email,omitempty"`
	City                 string     `json:"city,omitempty"`
	State                string     `json:"state,omitempty"`
	PersonalNumber       string     `json:"personal_number,omitempty"`
	ContactPerson        string     `json:"contact_person,omitempty"`
	ContactPersonNumber  string     `json:"contact_person_number,omitempty"`
	ReferenceBranchID    string     `json:"reference_branch_id,omitempty"`
	ReferenceVolunteerID string     `json:"reference_volunteer_id,omitempty"`
	ReferencePersonName  string     `json:"reference_person_name,omitempty"`
	EventID              uint       `json:"event_id"`
	Event                Event      `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
	CreatedOn            time.Time  `json:"created_on,omitempty"`
	UpdatedOn            *time.Time `json:"updated_on,omitempty"`
	CreatedBy            string     `json:"created_by,omitempty"`
	UpdatedBy            string     `json:"updated_by,omitempty"`
}
