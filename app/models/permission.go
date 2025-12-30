package models

import (
	"time"
)

// Permission represents a specific action that can be performed in the system
type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null;index" json:"name"`
	Resource    string    `gorm:"not null;index" json:"resource"`    // e.g., "users", "events", "branches"
	Action      string    `gorm:"not null;index" json:"action"`      // e.g., "create", "read", "update", "delete"
	Description string    `json:"description,omitempty"`
	CreatedOn   time.Time `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn   time.Time `gorm:"autoUpdateTime" json:"updated_on"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       uint      `gorm:"primaryKey" json:"role_id"`
	PermissionID uint      `gorm:"primaryKey" json:"permission_id"`
	Role         Role      `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
	Permission   Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE" json:"permission,omitempty"`
	GrantedOn    time.Time `gorm:"autoCreateTime" json:"granted_on"`
	GrantedBy    string    `json:"granted_by,omitempty"`
}

// TableName overrides the table name
func (RolePermission) TableName() string {
	return "role_permissions"
}


