package validators

import (
	"errors"
	"strings"
)
func ValidateSpecialGuestInput(eventID uint, prefix, firstName, lastName string) error {
	// Validate Event ID
	if eventID == 0 {
		return errors.New("event_id is required and must be greater than 0")
	}

	// Validate Prefix
	if strings.TrimSpace(prefix) == "" {
		return errors.New("prefix is required")
	}
	if len(prefix) < 1 || len(prefix) > 50 {
		return errors.New("prefix must be between 1 and 50 characters")
	}

	// Validate First Name
	if strings.TrimSpace(firstName) == "" {
		return errors.New("first_name is required")
	}
	if len(firstName) < 2 || len(firstName) > 255 {
		return errors.New("first_name must be between 2 and 255 characters")
	}

	// Validate Last Name
	if strings.TrimSpace(lastName) == "" {
		return errors.New("last_name is required")
	}
	if len(lastName) < 2 || len(lastName) > 255 {
		return errors.New("last_name must be between 2 and 255 characters")
	}

	return nil
}

// ValidateSpecialGuestUpdateFields validates special guest update request
func ValidateSpecialGuestUpdateFields(updateData map[string]interface{}) error {
	// List of fields that should not be updated
	immutableFields := map[string]bool{
		"id":         true,
		"created_on": true,
		"created_by": true,
		"event_id":   true, // event should not be changed
	}

	for field := range updateData {
		if immutableFields[field] {
			return errors.New("field '" + field + "' cannot be updated")
		}
	}

	// Validate specific fields if present
	if prefix, ok := updateData["prefix"]; ok {
		prefixStr := strings.TrimSpace(prefix.(string))
		if prefixStr == "" {
			return errors.New("prefix cannot be empty")
		}
		if len(prefixStr) < 1 || len(prefixStr) > 50 {
			return errors.New("prefix must be between 1 and 50 characters")
		}
	}

	if firstName, ok := updateData["first_name"]; ok {
		nameStr := strings.TrimSpace(firstName.(string))
		if nameStr != "" && (len(nameStr) < 2 || len(nameStr) > 255) {
			return errors.New("first_name must be between 2 and 255 characters")
		}
	}

	if lastName, ok := updateData["last_name"]; ok {
		nameStr := strings.TrimSpace(lastName.(string))
		if nameStr != "" && (len(nameStr) < 2 || len(nameStr) > 255) {
			return errors.New("last_name must be between 2 and 255 characters")
		}
	}

	if middleName, ok := updateData["middle_name"]; ok {
		nameStr := strings.TrimSpace(middleName.(string))
		if nameStr != "" && (len(nameStr) < 2 || len(nameStr) > 255) {
			return errors.New("middle_name must be between 2 and 255 characters")
		}
	}

	if designation, ok := updateData["designation"]; ok {
		designStr := strings.TrimSpace(designation.(string))
		if designStr != "" && (len(designStr) < 2 || len(designStr) > 255) {
			return errors.New("designation must be between 2 and 255 characters")
		}
	}

	if organization, ok := updateData["organization"]; ok {
		orgStr := strings.TrimSpace(organization.(string))
		if orgStr != "" && (len(orgStr) < 2 || len(orgStr) > 255) {
			return errors.New("organization must be between 2 and 255 characters")
		}
	}

	if email, ok := updateData["email"]; ok {
		emailStr := strings.TrimSpace(email.(string))
		if emailStr != "" {
			if !isValidEmail(emailStr) {
				return errors.New("invalid email format")
			}
		}
	}

	if city, ok := updateData["city"]; ok {
		cityStr := strings.TrimSpace(city.(string))
		if cityStr != "" && (len(cityStr) < 2 || len(cityStr) > 100) {
			return errors.New("city must be between 2 and 100 characters")
		}
	}

	if state, ok := updateData["state"]; ok {
		stateStr := strings.TrimSpace(state.(string))
		if stateStr != "" && (len(stateStr) < 2 || len(stateStr) > 100) {
			return errors.New("state must be between 2 and 100 characters")
		}
	}

	if personalNumber, ok := updateData["personal_number"]; ok {
		numStr := strings.TrimSpace(personalNumber.(string))
		if numStr != "" {
			if !isValidPhoneNumber(numStr) {
				return errors.New("invalid personal_number format")
			}
		}
	}

	if contactPersonNumber, ok := updateData["contact_person_number"]; ok {
		numStr := strings.TrimSpace(contactPersonNumber.(string))
		if numStr != "" {
			if !isValidPhoneNumber(numStr) {
				return errors.New("invalid contact_person_number format")
			}
		}
	}

	return nil
}
