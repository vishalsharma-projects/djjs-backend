package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// Create a new PromotionMaterialDetails record
func CreatePromotionMaterialDetails(detail *models.PromotionMaterialDetails) error {
	return config.DB.Create(detail).Error
}

// Get all PromotionMaterialDetails records
func GetAllPromotionMaterialDetails() ([]models.PromotionMaterialDetails, error) {
	var details []models.PromotionMaterialDetails
	if err := config.DB.
		Preload("Event").
		Find(&details).Error; err != nil {
		return nil, err
	}
	return details, nil
}

// GetPromotionMaterialDetailsByEventID retrieves all PromotionMaterialDetails records by EventID
func GetPromotionMaterialDetailsByEventID(eventID uint) ([]models.PromotionMaterialDetails, error) {
	var details []models.PromotionMaterialDetails
	if err := config.DB.
		Preload("Event").
		Where("event_id = ?", eventID).
		Find(&details).Error; err != nil {
		return nil, errors.New("no promotion material details found for the given event ID")
	}
	return details, nil
}

// Update a record
// UpdatePromotionMaterialDetails updates an existing PromotionMaterialDetails record
func UpdatePromotionMaterialDetails(detail *models.PromotionMaterialDetails) error {
	var existing models.PromotionMaterialDetails

	// Check if record exists
	if err := config.DB.First(&existing, detail.ID).Error; err != nil {
		return errors.New("record not found")
	}

	// Prepare dynamic updates
	updates := map[string]interface{}{
		"updated_on": time.Now(),
	}

	if detail.PromotionMaterialID != 0 {
		updates["promotion_material_id"] = detail.PromotionMaterialID
	}
	if detail.EventID != 0 {
		updates["event_id"] = detail.EventID
	}
	if detail.Quantity != 0 {
		updates["quantity"] = detail.Quantity
	}
	if detail.Size != "" {
		updates["size"] = detail.Size
	}
	if detail.DimensionHeight != 0 {
		updates["dimension_height"] = detail.DimensionHeight
	}
	if detail.DimensionWidth != 0 {
		updates["dimension_width"] = detail.DimensionWidth
	}
	if detail.UpdatedBy != "" {
		updates["updated_by"] = detail.UpdatedBy
	}

	// Apply updates
	return config.DB.Model(&existing).Updates(updates).Error
}

// Delete a record
func DeletePromotionMaterialDetails(id uint) error {
	result := config.DB.Delete(&models.PromotionMaterialDetails{}, id)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return result.Error
}
