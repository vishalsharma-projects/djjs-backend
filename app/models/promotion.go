package models

import (
	"time"
)

// PromotionMaterial represents types of promotion materials
type PromotionMaterial struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	MaterialType string `gorm:"not null" json:"material_type"`
}

func (PromotionMaterial) TableName() string {
	return "promotion_material_type"
}

// Event represents an event
type Event struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	EventName string `gorm:"not null" json:"event_name"`
}

func (Event) TableName() string {
	return "event_details"
}

// PromotionMaterialDetails represents detailed info about promotion materials for an event
type PromotionMaterialDetails struct {
	ID                  uint              `gorm:"primaryKey" json:"id"`
	PromotionMaterialID uint              `gorm:"not null" json:"promotion_material_id"`
	PromotionMaterial   PromotionMaterial `gorm:"foreignKey:PromotionMaterialID;references:ID" json:"promotion_material,omitempty"`
	EventID             uint              `json:"event_id"`
	Event               Event             `gorm:"foreignKey:EventID;references:ID" json:"event,omitempty"`
	Quantity            int               `gorm:"not null" json:"quantity"`
	Size                string            `json:"size,omitempty"`
	DimensionHeight     float64           `json:"dimension_height,omitempty"`
	DimensionWidth      float64           `json:"dimension_width,omitempty"`
	CreatedOn           time.Time         `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn           time.Time         `gorm:"autoUpdateTime" json:"updated_on"`
	CreatedBy           string            `json:"created_by,omitempty"`
	UpdatedBy           string            `json:"updated_by,omitempty"`
}

func (PromotionMaterialDetails) TableName() string {
	return "promotion_material_details"
}
