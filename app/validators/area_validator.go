package validators

import (
	"errors"
	"strings"
)

// ValidateAreaInput validates area creation data
func ValidateAreaInput(branchID uint, districtID string, areaName string, areaCoverage float64) error {
	// Validate Branch ID
	if branchID == 0 {
		return errors.New("branch_id is required and must be greater than 0")
	}

	// Validate District ID
	if strings.TrimSpace(districtID) == "" {
		return errors.New("district_id is required")
	}

	// Validate Area Name (optional but if provided, validate length)
	if areaName != "" {
		if len(areaName) < 2 || len(areaName) > 255 {
			return errors.New("area_name must be between 2 and 255 characters")
		}
	}

	// Validate Area Coverage (optional but if provided, must be positive)
	if areaCoverage < 0 {
		return errors.New("area_coverage must be a non-negative number")
	}

	return nil
}

// ValidateAreaUpdateFields validates area update request
func ValidateAreaUpdateFields(updateData map[string]interface{}) error {
	// List of fields that should not be updated
	immutableFields := map[string]bool{
		"id":         true,
		"created_on": true,
		"created_by": true,
		"branch_id":  true,   // branch should not be changed via update
		"district_id": true,  // district should not be changed via update
	}

	for field := range updateData {
		if immutableFields[field] {
			return errors.New("field '" + field + "' cannot be updated")
		}
	}

	// Validate specific fields if present
	if areaName, ok := updateData["area_name"]; ok {
		nameStr := strings.TrimSpace(areaName.(string))
		if nameStr != "" && (len(nameStr) < 2 || len(nameStr) > 255) {
			return errors.New("area_name must be between 2 and 255 characters")
		}
	}

	if areaCoverage, ok := updateData["area_coverage"]; ok {
		coverage, _ := areaCoverage.(float64)
		if coverage < 0 {
			return errors.New("area_coverage must be a non-negative number")
		}
	}

	if districtCoverage, ok := updateData["district_coverage"]; ok {
		coverage, _ := districtCoverage.(float64)
		if coverage < 0 {
			return errors.New("district_coverage must be a non-negative number")
		}
	}

	return nil
}
