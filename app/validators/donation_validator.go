package validators

import (
	"errors"
	"strings"
)

// ValidateDonationInput validates donation creation data
func ValidateDonationInput(eventID, branchID uint, donationType string, amount float64) error {
	// Validate Event ID
	if eventID == 0 {
		return errors.New("event_id is required and must be greater than 0")
	}

	// Validate Branch ID
	if branchID == 0 {
		return errors.New("branch_id is required and must be greater than 0")
	}

	// Validate Donation Type (optional)
	if donationType != "" {
		if len(donationType) < 2 || len(donationType) > 100 {
			return errors.New("donation_type must be between 2 and 100 characters")
		}
	}

	// Validate Amount (optional but if provided, must be positive)
	if amount < 0 {
		return errors.New("amount must be a non-negative number")
	}

	return nil
}

// ValidateDonationUpdateFields validates donation update request
func ValidateDonationUpdateFields(updateData map[string]interface{}) error {
	// List of fields that should not be updated
	immutableFields := map[string]bool{
		"id":         true,
		"created_on": true,
		"created_by": true,
		"event_id":   true,   // event should not be changed after creation
		"branch_id":  true,   // branch should not be changed after creation
	}

	for field := range updateData {
		if immutableFields[field] {
			return errors.New("field '" + field + "' cannot be updated")
		}
	}

	// Validate specific fields if present
	if donationType, ok := updateData["donation_type"]; ok {
		typeStr := strings.TrimSpace(donationType.(string))
		if typeStr != "" && (len(typeStr) < 2 || len(typeStr) > 100) {
			return errors.New("donation_type must be between 2 and 100 characters")
		}
	}

	if amount, ok := updateData["amount"]; ok {
		amountVal, _ := amount.(float64)
		if amountVal < 0 {
			return errors.New("amount must be a non-negative number")
		}
	}

	if kindType, ok := updateData["kindtype"]; ok {
		kindStr := strings.TrimSpace(kindType.(string))
		if kindStr != "" && (len(kindStr) < 2 || len(kindStr) > 255) {
			return errors.New("kindtype must be between 2 and 255 characters")
		}
	}

	return nil
}
