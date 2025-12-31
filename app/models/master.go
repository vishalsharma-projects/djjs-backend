package models

import "time"

type Country struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type State struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	CountryID uint   `json:"country_id"`
}

type City struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `json:"name"`
	StateID uint   `json:"state_id"`
}

type District struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	StateID   uint   `json:"state_id"`
	CountryID uint   `json:"country_id"`
}

type PromotionMaterialType struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	MaterialType string `json:"material_type"`
}

type Language struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	CreatedOn time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
}

type SevaType struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CreatedOn   time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn   *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
}

type Theme struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	CreatedOn time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
}

type InfrastructureType struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"not null;unique" json:"name"`
	CreatedOn time.Time  `gorm:"autoCreateTime" json:"created_on,omitempty"`
	UpdatedOn *time.Time `gorm:"autoUpdateTime" json:"updated_on,omitempty"`
}

func (InfrastructureType) TableName() string {
	return "infrastructure_types"
}