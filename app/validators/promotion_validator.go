package validators

import (
	"errors"
	"strings"
)

// ValidatePromotionMaterialDetailsInput validates promotion material details creation data
func ValidatePromotionMaterialDetailsInput(eventID, promotionMaterialID uint, quantity int) error {
	// Validate Event ID
	if eventID == 0 {
		return errors.New("event_id is required and must be greater than 0")
	}

	// Validate Promotion Material ID
	if promotionMaterialID == 0 {
		return errors.New("promotion_material_id is required and must be greater than 0")
	}

	// Validate Quantity
	if quantity < 0 {
		return errors.New("quantity must be a non-negative number")
	}

	return nil
}

// ValidatePromotionMaterialDetailsUpdateFields validates promotion material details update request
func ValidatePromotionMaterialDetailsUpdateFields(updateData map[string]interface{}) error {
	// List of fields that should not be updated
	immutableFields := map[string]bool{
		"id":                      true,
		"created_on":              true,
		"created_by":              true,
		"event_id":                true, // event should not be changed
		"promotion_material_id":   true, // material should not be changed
	}

	for field := range updateData {
		if immutableFields[field] {
			return errors.New("field '" + field + "' cannot be updated")
		}
	}

	// Validate specific fields if present
	if quantity, ok := updateData["quantity"]; ok {
		quantityVal, _ := quantity.(float64)
		if quantityVal < 0 {
			return errors.New("quantity must be a non-negative number")
		}
	}

	if size, ok := updateData["size"]; ok {
		sizeStr := strings.TrimSpace(size.(string))
		if sizeStr != "" && (len(sizeStr) < 1 || len(sizeStr) > 50) {
			return errors.New("size must be between 1 and 50 characters")
		}
	}

	if height, ok := updateData["dimension_height"]; ok {
		heightVal, _ := height.(float64)
		if heightVal < 0 {
			return errors.New("dimension_height must be a non-negative number")
		}
	}

	if width, ok := updateData["dimension_width"]; ok {
		widthVal, _ := width.(float64)
		if widthVal < 0 {
			return errors.New("dimension_width must be a non-negative number")
		}
	}

	return nil
}
