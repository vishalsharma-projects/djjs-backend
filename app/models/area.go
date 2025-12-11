package models

import (
	"time"

	"github.com/google/uuid"
)

// type Branch struct {
// 	// swagger:model Branch
// 	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
// 	Email           string     `gorm:"unique" json:"email,omitempty"`
// 	Name            string     `gorm:"not null" json:"name"`
// 	CoordinatorName string     `json:"coordinator_name,omitempty"`
// 	ContactNumber   string     `gorm:"unique;not null" json:"contact_number"`
// 	EstablishedDate *time.Time `json:"established_date,omitempty"`
// 	AashramArea     float64    `json:"aashram_area,omitempty"`
// 	CreatedOn       time.Time  `json:"created_on,omitempty"`
// 	UpdatedOn       *time.Time `json:"updated_on,omitempty"`
// 	CreatedBy       string     `json:"created_by,omitempty"`
// 	UpdatedBy       string     `json:"updated_by,omitempty"`
// }

type Area struct {
	// swagger:model Area
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchID         uint       `gorm:"not null" json:"branch_id"`
	Branch           Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	DistrictID       uuid.UUID  `gorm:"type:uuid;not null" json:"district_id"`
	DistrictCoverage float64    `json:"district_coverage,omitempty"`
	AreaName         string     `json:"area_name,omitempty"`
	AreaCoverage     float64    `json:"area_coverage,omitempty"`
	CreatedOn        time.Time  `json:"created_on,omitempty"`
	UpdatedOn        *time.Time `json:"updated_on,omitempty"`
	CreatedBy        string     `json:"created_by,omitempty"`
	UpdatedBy        string     `json:"updated_by,omitempty"`
}
