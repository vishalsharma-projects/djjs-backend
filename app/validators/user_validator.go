package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidateUserInput performs comprehensive validation on user data
func ValidateUserInput(name, email, contactNumber string, roleID uint) error {
	// Validate Name
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required and cannot be empty")
	}
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	if len(name) > 255 {
		return errors.New("name must not exceed 255 characters")
	}

	// Validate Email
	if strings.TrimSpace(email) == "" {
		return errors.New("email is required and cannot be empty")
	}
	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}
	if len(email) > 255 {
		return errors.New("email must not exceed 255 characters")
	}

	// Validate Contact Number (optional)
	if contactNumber != "" {
		if !isValidPhoneNumber(contactNumber) {
			return errors.New("invalid contact number format (expected 10 digits or +91 format)")
		}
	}

	return nil
}

// ValidateEmailFormat validates email format
func ValidateEmailFormat(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// ValidateContactNumber validates phone number format
func ValidateContactNumber(contactNumber string) error {
	if contactNumber == "" {
		return nil // optional field
	}
	if !isValidPhoneNumber(contactNumber) {
		return errors.New("invalid contact number format (expected 10 digits or +91XXXXXXXXXX format)")
	}
	return nil
}

// ValidateUpdateFields validates update request fields
func ValidateUpdateFields(updateData map[string]interface{}) error {
	// List of fields that should not be updated
	immutableFields := map[string]bool{
		"id":         true,
		"created_on": true,
		"created_by": true,
		"password":   true, // password should be updated via separate endpoint
		"role_id":    true, // role should be updated via separate endpoint
	}

	for field := range updateData {
		if immutableFields[field] {
			return fmt.Errorf("field '%s' cannot be updated", field)
		}
	}

	// Validate specific fields if present
	if name, ok := updateData["name"]; ok {
		nameStr := strings.TrimSpace(name.(string))
		if nameStr == "" {
			return errors.New("name cannot be empty")
		}
		if len(nameStr) < 2 || len(nameStr) > 255 {
			return errors.New("name must be between 2 and 255 characters")
		}
	}

	if email, ok := updateData["email"]; ok {
		if err := ValidateEmailFormat(email.(string)); err != nil {
			return err
		}
	}

	if contactNumber, ok := updateData["contact_number"]; ok {
		if err := ValidateContactNumber(contactNumber.(string)); err != nil {
			return err
		}
	}

	return nil
}

// ValidateSearchInput validates search parameters
func ValidateSearchInput(email, contact string) error {
	// At least one search parameter must be provided
	if strings.TrimSpace(email) == "" && strings.TrimSpace(contact) == "" {
		return errors.New("at least one search parameter (email or contact_number) is required")
	}

	// Validate email if provided
	if email != "" {
		if err := ValidateEmailFormat(email); err != nil {
			return err
		}
	}

	// Validate contact number if provided
	if contact != "" {
		if err := ValidateContactNumber(contact); err != nil {
			return err
		}
	}

	return nil
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ValidatePasswordChange validates password change request
func ValidatePasswordChange(oldPassword, newPassword, confirmPassword string) error {
	if strings.TrimSpace(oldPassword) == "" {
		return errors.New("old password is required")
	}

	if err := ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	if newPassword != confirmPassword {
		return errors.New("new password and confirm password do not match")
	}

	if oldPassword == newPassword {
		return errors.New("new password must be different from old password")
	}

	return nil
}

// ValidatePasswordStrength validates password strength
func ValidatePasswordStrength(password string) error {
	if strings.TrimSpace(password) == "" {
		return errors.New("password is required and cannot be empty")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 255 {
		return errors.New("password must not exceed 255 characters")
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check for at least one digit
	if !regexp.MustCompile(`\d`).MatchString(password) {
		return errors.New("password must contain at least one digit")
	}

	// Check for at least one special character
	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password) {
		return errors.New("password must contain at least one special character (!@#$%^&*)")
	}

	return nil
}

// Helper function to validate phone number format
// Accepts: 10 digits, +91XXXXXXXXXX, or (XXX) XXX-XXXX
func isValidPhoneNumber(phone string) bool {
	// Remove common formatting characters
	cleaned := regexp.MustCompile(`[\s\-\(\).]`).ReplaceAllString(phone, "")

	// Check for +91 followed by 10 digits
	if regexp.MustCompile(`^\+91\d{10}$`).MatchString(cleaned) {
		return true
	}

	// Check for 10 digits
	if regexp.MustCompile(`^\d{10}$`).MatchString(cleaned) {
		return true
	}

	return false
}
