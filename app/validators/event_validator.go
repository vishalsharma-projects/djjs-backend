package validators

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// ValidateEventInput validates event creation data
func ValidateEventInput(eventTypeID, eventCategoryID uint, startDate, endDate time.Time) error {
	// Validate Event Type ID
	if eventTypeID == 0 {
		return errors.New("event_type_id is required and must be greater than 0")
	}

	// Validate Event Category ID
	if eventCategoryID == 0 {
		return errors.New("event_category_id is required and must be greater than 0")
	}

	// Validate Start Date
	if startDate.IsZero() {
		return errors.New("start_date is required")
	}

	// Validate End Date
	if endDate.IsZero() {
		return errors.New("end_date is required")
	}

	// Validate End Date is after Start Date
	if endDate.Before(startDate) {
		return errors.New("end_date must be after or equal to start_date")
	}

	return nil
}

// ValidateEventUpdateFields validates event update request
func ValidateEventUpdateFields(updateData map[string]interface{}) error {
	// List of fields that should not be updated
	immutableFields := map[string]bool{
		"id":         true,
		"created_on": true,
		"created_by": true,
	}

	for field := range updateData {
		if immutableFields[field] {
			return errors.New("field '" + field + "' cannot be updated")
		}
	}

	// Validate specific fields if present
	if scale, ok := updateData["scale"]; ok {
		scaleStr := strings.TrimSpace(scale.(string))
		if scaleStr != "" && (len(scaleStr) < 2 || len(scaleStr) > 100) {
			return errors.New("scale must be between 2 and 100 characters")
		}
	}

	if theme, ok := updateData["theme"]; ok {
		themeStr := strings.TrimSpace(theme.(string))
		if themeStr != "" && (len(themeStr) < 2 || len(themeStr) > 500) {
			return errors.New("theme must be between 2 and 500 characters")
		}
	}

	if startDate, ok := updateData["start_date"]; ok {
		if dateStr, ok := startDate.(string); ok {
			if _, err := time.Parse("2006-01-02", dateStr); err != nil {
				return errors.New("invalid start_date format (use YYYY-MM-DD)")
			}
		}
	}

	if endDate, ok := updateData["end_date"]; ok {
		if dateStr, ok := endDate.(string); ok {
			if _, err := time.Parse("2006-01-02", dateStr); err != nil {
				return errors.New("invalid end_date format (use YYYY-MM-DD)")
			}
		}
	}

	if pincode, ok := updateData["pincode"]; ok {
		pincodeStr := strings.TrimSpace(pincode.(string))
		if pincodeStr != "" {
			if !regexp.MustCompile(`^\d{5,6}$`).MatchString(pincodeStr) {
				return errors.New("pincode must be 5-6 digits")
			}
		}
	}

	// Validate beneficiary counts
	if benMen, ok := updateData["beneficiary_men"]; ok {
		val, _ := benMen.(float64)
		if val < 0 {
			return errors.New("beneficiary_men must be a non-negative number")
		}
	}

	if benWomen, ok := updateData["beneficiary_women"]; ok {
		val, _ := benWomen.(float64)
		if val < 0 {
			return errors.New("beneficiary_women must be a non-negative number")
		}
	}

	if benChild, ok := updateData["beneficiary_child"]; ok {
		val, _ := benChild.(float64)
		if val < 0 {
			return errors.New("beneficiary_child must be a non-negative number")
		}
	}

	// Validate initiation counts
	if initMen, ok := updateData["initiation_men"]; ok {
		val, _ := initMen.(float64)
		if val < 0 {
			return errors.New("initiation_men must be a non-negative number")
		}
	}

	if initWomen, ok := updateData["initiation_women"]; ok {
		val, _ := initWomen.(float64)
		if val < 0 {
			return errors.New("initiation_women must be a non-negative number")
		}
	}

	if initChild, ok := updateData["initiation_child"]; ok {
		val, _ := initChild.(float64)
		if val < 0 {
			return errors.New("initiation_child must be a non-negative number")
		}
	}

	// Validate status if present
	if status, ok := updateData["status"]; ok {
		statusStr := strings.TrimSpace(status.(string))
		if statusStr != "" && statusStr != "complete" && statusStr != "incomplete" {
			return errors.New("status must be either 'complete' or 'incomplete'")
		}
	}

	// Validate spiritual_orator if present (can be comma-separated for multiple orators)
	if spiritualOrator, ok := updateData["spiritual_orator"]; ok {
		if spiritualOratorStr, ok := spiritualOrator.(string); ok {
			if len(spiritualOratorStr) > 200 {
				return errors.New("spiritual_orator must not exceed 200 characters")
			}
		}
	}

	return nil
}
