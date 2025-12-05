package validators

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// ValidateBranchInput validates branch creation data
func ValidateBranchInput(name, email, contactNumber, coordinatorName string) error {
	// Validate Name
	if strings.TrimSpace(name) == "" {
		return errors.New("branch name is required and cannot be empty")
	}
	if len(name) < 2 {
		return errors.New("branch name must be at least 2 characters long")
	}
	if len(name) > 255 {
		return errors.New("branch name must not exceed 255 characters")
	}

	// Validate Email (optional)
	if email != "" {
		if err := ValidateEmailFormat(email); err != nil {
			return err
		}
	}

	// Validate Contact Number
	if strings.TrimSpace(contactNumber) == "" {
		return errors.New("contact number is required")
	}
	if err := ValidateContactNumber(contactNumber); err != nil {
		return err
	}

	// Validate Coordinator Name (optional)
	if coordinatorName != "" {
		if len(coordinatorName) < 2 || len(coordinatorName) > 255 {
			return errors.New("coordinator name must be between 2 and 255 characters")
		}
	}

	return nil
}

// ValidateBranchUpdateFields validates branch update request
func ValidateBranchUpdateFields(updateData map[string]interface{}) error {
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
	if name, ok := updateData["name"]; ok {
		nameStr := strings.TrimSpace(name.(string))
		if nameStr == "" {
			return errors.New("branch name cannot be empty")
		}
		if len(nameStr) < 2 || len(nameStr) > 255 {
			return errors.New("branch name must be between 2 and 255 characters")
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

	if coordinatorName, ok := updateData["coordinator_name"]; ok {
		coordinatorStr := strings.TrimSpace(coordinatorName.(string))
		if coordinatorStr != "" && (len(coordinatorStr) < 2 || len(coordinatorStr) > 255) {
			return errors.New("coordinator name must be between 2 and 255 characters")
		}
	}

	if aashramArea, ok := updateData["aashram_area"]; ok {
		area, _ := aashramArea.(float64)
		if area < 0 {
			return errors.New("aashram area must be a positive number")
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

	if establishedOn, ok := updateData["established_on"]; ok {
		if establishedOn != nil {
			if dateStr, ok := establishedOn.(string); ok {
				if _, err := time.Parse("2006-01-02", dateStr); err != nil {
					return errors.New("invalid established_on date format (use YYYY-MM-DD)")
				}
			}
		}
	}

	return nil
}

// ValidateBranchInfrastructure validates branch infrastructure data
func ValidateBranchInfrastructure(branchID uint, infraType string, count int) error {
	if branchID == 0 {
		return errors.New("branch_id is required and must be greater than 0")
	}

	if strings.TrimSpace(infraType) == "" {
		return errors.New("infrastructure type is required")
	}

	if len(infraType) < 2 || len(infraType) > 100 {
		return errors.New("infrastructure type must be between 2 and 100 characters")
	}

	if count < 0 {
		return errors.New("count must be a non-negative number")
	}

	return nil
}

// ValidateBranchMember validates branch member data
func ValidateBranchMember(name, memberType string, branchID uint) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("member name is required")
	}

	if len(name) < 2 || len(name) > 255 {
		return errors.New("member name must be between 2 and 255 characters")
	}

	if strings.TrimSpace(memberType) == "" {
		return errors.New("member type is required")
	}

	if branchID == 0 {
		return errors.New("branch_id is required and must be greater than 0")
	}

	return nil
}

// ValidateBranchInfrastructureUpdateFields validates branch infrastructure update
func ValidateBranchInfrastructureUpdateFields(updateData map[string]interface{}) error {
	// List of fields that should not be updated
	immutableFields := map[string]bool{
		"id":         true,
		"created_on": true,
		"created_by": true,
		"branch_id":  true, // branch should not be changed via update
		"branch":     true, // branch relation should not be changed
	}

	for field := range updateData {
		if immutableFields[field] {
			return errors.New("field '" + field + "' cannot be updated")
		}
	}

	// Validate specific fields if present
	if infraType, ok := updateData["type"]; ok {
		typeStr := strings.TrimSpace(infraType.(string))
		if typeStr == "" {
			return errors.New("infrastructure type cannot be empty")
		}
		if len(typeStr) < 2 || len(typeStr) > 100 {
			return errors.New("infrastructure type must be between 2 and 100 characters")
		}
	}

	if count, ok := updateData["count"]; ok {
		countVal, _ := count.(float64)
		if countVal < 0 {
			return errors.New("count must be a non-negative number")
		}
	}

	return nil
}

// ValidateBranchMemberUpdateFields validates branch member update
func ValidateBranchMemberUpdateFields(updateData map[string]interface{}) error {
	immutableFields := map[string]bool{
		"id":         true,
		"created_on": true,
		"created_by": true,
		"branch_id":  true, // branch should not be changed via update
	}

	for field := range updateData {
		if immutableFields[field] {
			return errors.New("field '" + field + "' cannot be updated")
		}
	}

	// Validate specific fields if present
	if name, ok := updateData["name"]; ok {
		nameStr := strings.TrimSpace(name.(string))
		if nameStr == "" {
			return errors.New("member name cannot be empty")
		}
		if len(nameStr) < 2 || len(nameStr) > 255 {
			return errors.New("member name must be between 2 and 255 characters")
		}
	}

	if age, ok := updateData["age"]; ok {
		ageVal, _ := age.(float64)
		if ageVal < 0 || ageVal > 150 {
			return errors.New("age must be between 0 and 150")
		}
	}

	return nil
}

// ValidateTimeFormat validates time format (HH:MM)
func ValidateTimeFormat(timeStr string) error {
	if strings.TrimSpace(timeStr) == "" {
		return nil // optional field
	}

	if !regexp.MustCompile(`^([0-1][0-9]|2[0-3]):[0-5][0-9]$`).MatchString(timeStr) {
		return errors.New("time must be in HH:MM format (24-hour)")
	}

	return nil
}
