package models

import "time"

// swagger:model Branch
type Branch struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"not null" json:"name" validate:"required,min=2,max=255"`
	Email           string     `gorm:"unique" json:"email,omitempty" validate:"omitempty,email,max=255"`
	CoordinatorName string     `json:"coordinator_name,omitempty" validate:"omitempty,min=2,max=255"`
	ContactNumber   string     `gorm:"unique;not null" json:"contact_number" validate:"required,max=20"`
	EstablishedOn   *time.Time `json:"established_on,omitempty"`
	AashramArea     float64    `json:"aashram_area,omitempty" validate:"omitempty,min=0"`
	CountryID       *uint      `gorm:"column:country_id" json:"country_id" validate:"omitempty,min=1"`
	Country         Country    `gorm:"foreignKey:CountryID" json:"country,omitempty"`
	StateID         *uint      `gorm:"column:state_id" json:"state_id" validate:"omitempty,min=1"`
	State           State      `gorm:"foreignKey:StateID" json:"state,omitempty"`
	DistrictID      *uint      `gorm:"column:district_id" json:"district_id" validate:"omitempty,min=1"`
	District        District   `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	CityID          *uint      `gorm:"column:city_id" json:"city_id" validate:"omitempty,min=1"`
	City            City       `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Address         string     `json:"address,omitempty" validate:"omitempty,max=500"`
	Pincode         string     `json:"pincode,omitempty" validate:"omitempty,numeric,len=5|len=6"`
	PostOffice      string     `json:"post_office,omitempty" validate:"omitempty,max=100"`
	PoliceStation   string     `json:"police_station,omitempty" validate:"omitempty,max=100"`
	OpenDays        string     `json:"open_days,omitempty" validate:"omitempty,max=100"`
	DailyStartTime  string     `json:"daily_start_time,omitempty" validate:"omitempty"`
	DailyEndTime    string     `json:"daily_end_time,omitempty" validate:"omitempty"`
	ParentBranchID  *uint      `gorm:"column:parent_branch_id" json:"parent_branch_id,omitempty"`
	Parent          *Branch    `gorm:"foreignKey:ParentBranchID" json:"parent,omitempty"`
	Children        []Branch   `gorm:"foreignKey:ParentBranchID" json:"children,omitempty"`
	Infrastructures []BranchInfrastructure `gorm:"foreignKey:BranchID" json:"infrastructure,omitempty"`
	Members         []BranchMember         `gorm:"foreignKey:BranchID" json:"branch_members,omitempty"`
	CreatedOn       time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn       *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
	CreatedBy       string     `json:"created_by,omitempty"`
	UpdatedBy       string     `json:"updated_by,omitempty"`
}

// swagger:model BranchInfrastructure
type BranchInfrastructure struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchID  uint       `gorm:"not null" json:"branch_id" validate:"required,min=1"`
	Branch    Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Type      string     `gorm:"not null" json:"type" validate:"required,min=2,max=100"`
	Count     int        `gorm:"not null" json:"count" validate:"required,min=0"`
	CreatedOn time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
	CreatedBy string     `json:"created_by,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`
}

func (BranchInfrastructure) TableName() string {
	return "branch_infrastructure"
}

type BranchMember struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	MemberType     string     `gorm:"not null" json:"member_type" validate:"required,min=2,max=100"`
	Name           string     `gorm:"not null" json:"name" validate:"required,min=2,max=255"`
	BranchRole     string     `json:"branch_role,omitempty" validate:"omitempty,max=100"`
	Responsibility string     `json:"responsibility,omitempty" validate:"omitempty,max=500"`
	Age            int        `json:"age,omitempty" validate:"omitempty,min=0,max=150"`
	DateOfSamarpan *time.Time `json:"date_of_samarpan,omitempty"`
	Qualification  string     `json:"qualification,omitempty" validate:"omitempty,max=255"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty"`
	BranchID       uint       `gorm:"not null" json:"branch_id" validate:"required,min=1"`
	Branch         Branch     `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	CreatedOn      time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn      *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
	CreatedBy      string     `json:"created_by,omitempty"`
	UpdatedBy      string     `json:"updated_by,omitempty"`
}

// Optional: override table name if GORM pluralizes incorrectly
func (BranchMember) TableName() string {
	return "branch_member"
}
