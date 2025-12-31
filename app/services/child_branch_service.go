package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// CreateChildBranch creates a new child branch (now using Branch model with parent_branch_id)
func CreateChildBranch(childBranch *models.Branch) error {
	// Ensure parent_branch_id is set (required for child branches)
	if childBranch.ParentBranchID == nil || *childBranch.ParentBranchID == 0 {
		return errors.New("parent_branch_id is required for child branches")
	}
	
	childBranch.CreatedOn = time.Now()
	
	// Ensure status is set to true when creating a child branch
	// If status is not explicitly set, default to true
	if !childBranch.Status {
		childBranch.Status = true
	}
	
	if err := config.DB.Create(childBranch).Error; err != nil {
		return err
	}
	return nil
}

// GetAllChildBranches fetches all child branches (branches with parent_branch_id set)
func GetAllChildBranches() ([]models.Branch, error) {
	var childBranches []models.Branch
	if err := config.DB.
		Where("parent_branch_id IS NOT NULL").
		Preload("Parent").
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City").
		Preload("Infrastructures").
		Preload("Members").
		Order("id DESC").
		Find(&childBranches).Error; err != nil {
		return nil, err
	}
	return childBranches, nil
}

// GetChildBranch fetches a child branch by ID (branch with parent_branch_id set)
func GetChildBranch(childBranchID uint) (*models.Branch, error) {
	var childBranch models.Branch
	if err := config.DB.
		Where("id = ? AND parent_branch_id IS NOT NULL", childBranchID).
		Preload("Parent").
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City").
		Preload("Infrastructures").
		Preload("Members").
		First(&childBranch).Error; err != nil {
		return nil, errors.New("child branch not found")
	}
	return &childBranch, nil
}

// GetChildBranchesByParent fetches all child branches of a parent branch
func GetChildBranchesByParent(parentBranchID uint) ([]models.Branch, error) {
	var childBranches []models.Branch
	if err := config.DB.
		Where("parent_branch_id = ?", parentBranchID).
		Preload("Parent").
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City").
		Preload("Infrastructures").
		Preload("Members").
		Order("id DESC").
		Find(&childBranches).Error; err != nil {
		return nil, err
	}
	return childBranches, nil
}

// UpdateChildBranch updates a child branch
func UpdateChildBranch(childBranchID uint, updatedData map[string]interface{}) error {
	var childBranch models.Branch
	if err := config.DB.Where("id = ? AND parent_branch_id IS NOT NULL", childBranchID).First(&childBranch).Error; err != nil {
		return errors.New("child branch not found")
	}

	// Validate parent_branch_id if being updated
	if parentID, ok := updatedData["parent_branch_id"]; ok {
		var parentIDVal uint
		switch v := parentID.(type) {
		case float64:
			parentIDVal = uint(v)
		case uint:
			parentIDVal = v
		case int:
			parentIDVal = uint(v)
		}
		if parentIDVal > 0 {
			var parentBranch models.Branch
			if err := config.DB.First(&parentBranch, parentIDVal).Error; err != nil {
				return errors.New("invalid parent_branch_id")
			}
		}
	}

	// Validate location IDs if being updated
	if countryID, ok := updatedData["country_id"]; ok && countryID != nil {
		var countryIDVal uint
		switch v := countryID.(type) {
		case float64:
			countryIDVal = uint(v)
		case uint:
			countryIDVal = v
		case int:
			countryIDVal = uint(v)
		}
		if countryIDVal > 0 {
			var country models.Country
			if err := config.DB.First(&country, countryIDVal).Error; err != nil {
				return errors.New("invalid country_id")
			}
		}
	}

	if stateID, ok := updatedData["state_id"]; ok && stateID != nil {
		var stateIDVal uint
		switch v := stateID.(type) {
		case float64:
			stateIDVal = uint(v)
		case uint:
			stateIDVal = v
		case int:
			stateIDVal = uint(v)
		}
		if stateIDVal > 0 {
			var state models.State
			if err := config.DB.First(&state, stateIDVal).Error; err != nil {
				return errors.New("invalid state_id")
			}
		}
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&childBranch).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteChildBranch deletes a child branch by ID
func DeleteChildBranch(childBranchID uint) error {
	// Only delete if it's actually a child branch (has parent_branch_id)
	var childBranch models.Branch
	if err := config.DB.Where("id = ? AND parent_branch_id IS NOT NULL", childBranchID).First(&childBranch).Error; err != nil {
		return errors.New("child branch not found")
	}
	if err := config.DB.Delete(&childBranch).Error; err != nil {
		return err
	}
	return nil
}

// *************************************** Child Branch Infrastructure ****************************************************** //
// Note: Child branch infrastructure now uses BranchInfrastructure model with branch_id

// CreateChildBranchInfrastructure creates a new child branch infrastructure record
func CreateChildBranchInfrastructure(infra *models.BranchInfrastructure) error {
	// Validate required fields
	if infra.BranchID == 0 {
		return errors.New("branch_id is required")
	}
	if infra.Type == "" {
		return errors.New("type is required")
	}
	if infra.Count < 0 {
		return errors.New("count must be 0 or greater")
	}

	infra.CreatedOn = time.Now()
	if err := config.DB.Create(infra).Error; err != nil {
		return err
	}
	return nil
}

// GetInfrastructureByChildBranch fetches infrastructure records by child branch ID
func GetInfrastructureByChildBranch(childBranchID uint) ([]models.BranchInfrastructure, error) {
	var infra []models.BranchInfrastructure
	if err := config.DB.Where("branch_id = ?", childBranchID).Preload("Branch").Find(&infra).Error; err != nil {
		return nil, err
	}
	return infra, nil
}

// UpdateChildBranchInfrastructure updates a child branch infrastructure record
func UpdateChildBranchInfrastructure(id uint, updatedData map[string]interface{}) error {
	var infra models.BranchInfrastructure
	if err := config.DB.First(&infra, id).Error; err != nil {
		return errors.New("infrastructure not found")
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&infra).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteChildBranchInfrastructure deletes a child branch infrastructure record
func DeleteChildBranchInfrastructure(id uint) error {
	if err := config.DB.Delete(&models.BranchInfrastructure{}, id).Error; err != nil {
		return err
	}
	return nil
}

// *************************************** Child Branch Member ****************************************************** //
// Note: Child branch members now use BranchMember model with branch_id

// CreateChildBranchMember creates a new child branch member
func CreateChildBranchMember(member *models.BranchMember) error {
	member.CreatedOn = time.Now()
	if err := config.DB.Create(member).Error; err != nil {
		return err
	}
	return nil
}

// GetMembersByChildBranch fetches all members of a child branch
func GetMembersByChildBranch(childBranchID uint) ([]models.BranchMember, error) {
	var members []models.BranchMember
	if err := config.DB.Where("branch_id = ?", childBranchID).Preload("Branch").Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// UpdateChildBranchMember updates a child branch member
func UpdateChildBranchMember(memberID uint, updatedData map[string]interface{}) error {
	var member models.BranchMember
	if err := config.DB.First(&member, memberID).Error; err != nil {
		return errors.New("member not found")
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&member).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteChildBranchMember deletes a child branch member
func DeleteChildBranchMember(memberID uint) error {
	if err := config.DB.Delete(&models.BranchMember{}, memberID).Error; err != nil {
		return err
	}
	return nil
}


