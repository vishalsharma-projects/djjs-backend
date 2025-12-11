package validators

import (
	"errors"
	"strings"
)

// ValidateVolunteerInput validates volunteer creation data
func ValidateVolunteerInput(volunteerName string, branchID, eventID uint, numberOfDays int) error {
	// Validate Volunteer Name
	if strings.TrimSpace(volunteerName) == "" {
		return errors.New("volunteer name is required and cannot be empty")
	}
	if len(volunteerName) < 2 {
		return errors.New("volunteer name must be at least 2 characters long")
	}
	if len(volunteerName) > 255 {
		return errors.New("volunteer name must not exceed 255 characters")
	}

	// Validate Branch ID
	if branchID == 0 {
		return errors.New("branch_id is required and must be greater than 0")
	}

	// Validate Event ID
	if eventID == 0 {
		return errors.New("event_id is required and must be greater than 0")
	}

	// Validate Number of Days
	if numberOfDays < 0 {
		return errors.New("number_of_days must be a non-negative number")
	}
	if numberOfDays > 365 {
		return errors.New("number_of_days cannot exceed 365")
	}

	return nil
}

// ValidateVolunteerUpdateFields validates volunteer update request
func ValidateVolunteerUpdateFields(updateData map[string]interface{}) error {
	// List of fields that should not be updated
	immutableFields := map[string]bool{
		"id":         true,
		"created_on": true,
		"created_by": true,
		"event_id":   true,   // event should not be changed after creation
		"branch_id":  true,   // branch should not be changed after creation
		"branch":     true,   // branch relation should not be changed
		"event":      true,   // event relation should not be changed
	}

	for field := range updateData {
		if immutableFields[field] {
			return errors.New("field '" + field + "' cannot be updated")
		}
	}

	// Validate specific fields if present
	if volunteerName, ok := updateData["volunteer_name"]; ok {
		nameStr := strings.TrimSpace(volunteerName.(string))
		if nameStr == "" {
			return errors.New("volunteer name cannot be empty")
		}
		if len(nameStr) < 2 || len(nameStr) > 255 {
			return errors.New("volunteer name must be between 2 and 255 characters")
		}
	}

	if numberOfDays, ok := updateData["number_of_days"]; ok {
		daysVal, _ := numberOfDays.(float64)
		if daysVal < 0 {
			return errors.New("number_of_days must be a non-negative number")
		}
		if daysVal > 365 {
			return errors.New("number_of_days cannot exceed 365")
		}
	}

	if sevaInvolved, ok := updateData["seva_involved"]; ok {
		sevaStr := strings.TrimSpace(sevaInvolved.(string))
		if sevaStr != "" && (len(sevaStr) < 2 || len(sevaStr) > 500) {
			return errors.New("seva_involved must be between 2 and 500 characters")
		}
	}

	if mentionSeva, ok := updateData["mention_seva"]; ok {
		mentionStr := strings.TrimSpace(mentionSeva.(string))
		if mentionStr != "" && (len(mentionStr) < 2 || len(mentionStr) > 500) {
			return errors.New("mention_seva must be between 2 and 500 characters")
		}
	}

	return nil
}

// ValidateSevaInvolved validates seva involved field
func ValidateSevaInvolved(sevaInvolved string) error {
	if sevaInvolved == "" {
		return nil // optional field
	}

	if len(sevaInvolved) < 2 || len(sevaInvolved) > 500 {
		return errors.New("seva_involved must be between 2 and 500 characters")
	}

	return nil
}

// ValidateMentionSeva validates mention seva field
func ValidateMentionSeva(mentionSeva string) error {
	if mentionSeva == "" {
		return nil // optional field
	}

	if len(mentionSeva) < 2 || len(mentionSeva) > 500 {
		return errors.New("mention_seva must be between 2 and 500 characters")
	}

	return nil
}

// ValidateNumberOfDays validates the number of days
func ValidateNumberOfDays(numberOfDays int) error {
	if numberOfDays < 0 {
		return errors.New("number_of_days must be a non-negative number")
	}
	if numberOfDays > 365 {
		return errors.New("number_of_days cannot exceed 365")
	}

	return nil
}
