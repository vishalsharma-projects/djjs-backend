package validators

import (
	"errors"
	"regexp"
	"strings"
)

// ValidateEventMediaInput validates event media creation data
func ValidateEventMediaInput(eventID, mediaCoverageTypeID uint, companyName, firstName, lastName string) error {
	// Validate Event ID
	if eventID == 0 {
		return errors.New("event_id is required and must be greater than 0")
	}

	// Validate Media Coverage Type ID
	if mediaCoverageTypeID == 0 {
		return errors.New("media_coverage_type_id is required and must be greater than 0")
	}

	// Validate Company Name
	if strings.TrimSpace(companyName) == "" {
		return errors.New("company_name is required")
	}
	if len(companyName) < 2 || len(companyName) > 255 {
		return errors.New("company_name must be between 2 and 255 characters")
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

// ValidateEventMediaUpdateFields validates event media update request
func ValidateEventMediaUpdateFields(updateData map[string]interface{}) error {
	// List of fields that should not be updated
	immutableFields := map[string]bool{
		"id":                      true,
		"created_on":              true,
		"created_by":              true,
		"event_id":                true, // event should not be changed
		"media_coverage_type_id":  true, // media type should not be changed
	}

	for field := range updateData {
		if immutableFields[field] {
			return errors.New("field '" + field + "' cannot be updated")
		}
	}

	// Validate specific fields if present
	if companyName, ok := updateData["company_name"]; ok {
		nameStr := strings.TrimSpace(companyName.(string))
		if nameStr == "" {
			return errors.New("company_name cannot be empty")
		}
		if len(nameStr) < 2 || len(nameStr) > 255 {
			return errors.New("company_name must be between 2 and 255 characters")
		}
	}

	if firstName, ok := updateData["first_name"]; ok {
		nameStr := strings.TrimSpace(firstName.(string))
		if nameStr == "" {
			return errors.New("first_name cannot be empty")
		}
		if len(nameStr) < 2 || len(nameStr) > 255 {
			return errors.New("first_name must be between 2 and 255 characters")
		}
	}

	if lastName, ok := updateData["last_name"]; ok {
		nameStr := strings.TrimSpace(lastName.(string))
		if nameStr == "" {
			return errors.New("last_name cannot be empty")
		}
		if len(nameStr) < 2 || len(nameStr) > 255 {
			return errors.New("last_name must be between 2 and 255 characters")
		}
	}

	if middleName, ok := updateData["middle_name"]; ok {
		nameStr := strings.TrimSpace(middleName.(string))
		if nameStr != "" && (len(nameStr) < 2 || len(nameStr) > 255) {
			return errors.New("middle_name must be between 2 and 255 characters")
		}
	}

	if companyEmail, ok := updateData["company_email"]; ok {
		emailStr := strings.TrimSpace(companyEmail.(string))
		if emailStr != "" {
			// Validate length (VARCHAR(500) in database)
			if len(emailStr) > 500 {
				return errors.New("company_email must not exceed 500 characters")
			}
			// Validate comma-separated emails
			if !isValidEmailOrEmails(emailStr) {
				return errors.New("invalid company_email format (must be valid email(s) separated by commas)")
			}
		}
	}

	if email, ok := updateData["email"]; ok {
		emailStr := strings.TrimSpace(email.(string))
		if emailStr != "" {
			// Validate length (VARCHAR(500) in database)
			if len(emailStr) > 500 {
				return errors.New("email must not exceed 500 characters")
			}
			// Validate comma-separated emails
			if !isValidEmailOrEmails(emailStr) {
				return errors.New("invalid email format (must be valid email(s) separated by commas)")
			}
		}
	}

	if contact, ok := updateData["contact"]; ok {
		contactStr := strings.TrimSpace(contact.(string))
		if contactStr != "" {
			if !isValidPhoneNumber(contactStr) {
				return errors.New("invalid contact number format")
			}
		}
	}

	if designation, ok := updateData["designation"]; ok {
		designStr := strings.TrimSpace(designation.(string))
		if designStr != "" && (len(designStr) < 2 || len(designStr) > 100) {
			return errors.New("designation must be between 2 and 100 characters")
		}
	}

	return nil
}

// isValidEmailOrEmails validates a single email or comma-separated emails
func isValidEmailOrEmails(emailStr string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	
	// Split by comma and validate each email
	emails := strings.Split(emailStr, ",")
	for _, email := range emails {
		trimmedEmail := strings.TrimSpace(email)
		if trimmedEmail != "" && !re.MatchString(trimmedEmail) {
			return false
		}
	}
	return true
}
