package models

import (
	"time"
)

// MediaCoverageType represents types of media coverage
type MediaCoverageType struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	MediaType string `gorm:"not null" json:"media_type"`
}

func (MediaCoverageType) TableName() string {
	return "media_coverage_type"
}

// EventMedia represents media coverage for a specific event
type EventMedia struct {
	ID                  uint              `gorm:"primaryKey" json:"id"`
	EventID             uint              `gorm:"not null" json:"event_id"`
	MediaCoverageTypeID uint              `gorm:"not null" json:"media_coverage_type_id"`
	CompanyName         string            `gorm:"not null" json:"company_name"`
	CompanyEmail        string            `json:"company_email,omitempty"`
	CompanyWebsite      string            `json:"company_website,omitempty"`
	Gender              string            `json:"gender,omitempty"`
	Prefix              string            `json:"prefix,omitempty"`
	FirstName           string            `gorm:"not null" json:"first_name"`
	MiddleName          string            `json:"middle_name,omitempty"`
	LastName            string            `gorm:"not null" json:"last_name"`
	Designation         string            `json:"designation,omitempty"`
	Contact             string            `json:"contact,omitempty"`
	Email               string            `json:"email,omitempty"`
	FileURL             string            `json:"file_url,omitempty" gorm:"column:file_url"`
	FileType            string            `json:"file_type,omitempty" gorm:"column:file_type"` // image, video, audio, file
	CreatedOn           time.Time         `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn           time.Time         `gorm:"autoUpdateTime" json:"updated_on"`
	CreatedBy           string            `json:"created_by,omitempty" gorm:"<-:create"` // only set on create
	UpdatedBy           string            `json:"updated_by,omitempty"`
	MediaCoverageType   MediaCoverageType `gorm:"foreignKey:MediaCoverageTypeID;references:ID" json:"media_coverage_type,omitempty"`
	Event               Event             `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
}

func (EventMedia) TableName() string {
	return "event_media"
}
