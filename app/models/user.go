package models

import (
	"time"
	//"gorm.io/gorm"
)

type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null" json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedOn   time.Time `json:"created_on,omitempty"`
	UpdatedOn   time.Time `json:"updated_on,omitempty"`
}

// User model represents the users table in PostgreSQL
// swagger:model User
type User struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"not null" json:"name" validate:"required,min=2,max=255"`
	Email         string     `gorm:"unique;not null" json:"email" validate:"required,email,max=255"`
	ContactNumber string     `json:"contact_number,omitempty" validate:"omitempty,max=20"`
	Password      string     `gorm:"not null" json:"password,omitempty"`
	RoleID        uint       `gorm:"not null" json:"role_id" validate:"required"`
	BranchID      *uint      `gorm:"column:branch_id" json:"branch_id,omitempty"`
	Role          Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Branch        *Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Token         string     `json:"token,omitempty"`
	ExpiredOn     *time.Time `json:"expired_on,omitempty"`
	LastLoginOn   *time.Time `json:"last_login_on,omitempty"`
	FirstLoginOn  *time.Time `json:"first_login_on,omitempty"`
	IsDeleted     bool       `gorm:"default:false" json:"is_deleted"`
	CreatedOn     time.Time  `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn     *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
	CreatedBy     string     `json:"created_by,omitempty"`
	UpdatedBy     string     `json:"updated_by,omitempty"`
}
