package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// CreateBranch inserts a new branch record
func CreateBranch(branch *models.Branch) error {
	// Check email uniqueness if provided
	if branch.Email != "" {
		var existingBranch models.Branch
		if err := config.DB.Where("email = ?", branch.Email).First(&existingBranch).Error; err == nil {
			return errors.New("email already exists")
		}
	}

	// Check contact number uniqueness
	var existingBranch models.Branch
	if err := config.DB.Where("contact_number = ?", branch.ContactNumber).First(&existingBranch).Error; err == nil {
		return errors.New("contact number already exists")
	}

	branch.CreatedOn = time.Now()
	branch.UpdatedOn = nil

	if err := config.DB.Create(branch).Error; err != nil {
		return err
	}
	return nil
}

// GetAllBranches fetches all branches
func GetAllBranches() ([]models.Branch, error) {
	var branches []models.Branch
	if err := config.DB.Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}

// GetBranchSearch fetches branches by name and/or coordinator name
func GetBranchSearch(branchName, coordinator string) ([]models.Branch, error) {
	var branches []models.Branch
	db := config.DB

	// Apply filters dynamically
	if branchName != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?)", "%"+branchName+"%")
	}
	if coordinator != "" {
		db = db.Where("LOWER(coordinator_name) LIKE LOWER(?)", "%"+coordinator+"%")
	}

	if err := db.Find(&branches).Error; err != nil {
		return nil, errors.New("error fetching branches")
	}

	if len(branches) == 0 {
		return nil, errors.New("no branches found")
	}

	return branches, nil
}

// UpdateBranch updates branch fields
func UpdateBranch(branchID uint, updatedData map[string]interface{}) error {
	var branch models.Branch
	if err := config.DB.First(&branch, branchID).Error; err != nil {
		return errors.New("branch not found")
	}

	// Check email uniqueness if email is being updated
	if email, ok := updatedData["email"]; ok {
		var existingBranch models.Branch
		if err := config.DB.Where("email = ? AND id != ?", email, branchID).First(&existingBranch).Error; err == nil {
			return errors.New("email already exists")
		}
	}

	// Check contact number uniqueness if being updated
	if contactNumber, ok := updatedData["contact_number"]; ok {
		var existingBranch models.Branch
		if err := config.DB.Where("contact_number = ? AND id != ?", contactNumber, branchID).First(&existingBranch).Error; err == nil {
			return errors.New("contact number already exists")
		}
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&branch).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteBranch deletes a branch by ID
func DeleteBranch(branchID uint) error {
	if err := config.DB.Delete(&models.Branch{}, branchID).Error; err != nil {
		return err
	}
	return nil
}

// *************************************** Branch Infrastructure ****************************************************** //

// CreateBranchInfrastructure inserts a new record
func CreateBranchInfrastructure(infra *models.BranchInfrastructure) error {
	infra.CreatedOn = time.Now()
	infra.UpdatedOn = nil

	if err := config.DB.Create(infra).Error; err != nil {
		return err
	}
	return nil
}

// GetAllBranchInfrastructure fetches all records
func GetAllBranchInfrastructure() ([]models.BranchInfrastructure, error) {
	var infra []models.BranchInfrastructure
	if err := config.DB.Preload("Branch").Find(&infra).Error; err != nil {
		return nil, err
	}
	return infra, nil
}

// GetInfrastructureByBranch fetches records by branch ID
func GetInfrastructureByBranch(branchID uint) ([]models.BranchInfrastructure, error) {
	var infra []models.BranchInfrastructure
	if err := config.DB.Where("branch_id = ?", branchID).Preload("Branch").Find(&infra).Error; err != nil {
		return nil, err
	}
	return infra, nil
}

// UpdateBranchInfrastructure updates a record by ID
func UpdateBranchInfrastructure(id uint, updatedData map[string]interface{}) error {
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

// DeleteBranchInfrastructure deletes a record by ID
func DeleteBranchInfrastructure(id uint) error {
	if err := config.DB.Delete(&models.BranchInfrastructure{}, id).Error; err != nil {
		return err
	}
	return nil
}

// *************************************** Branch Member ****************************************************** //

// CreateBranchMember inserts a new branch member
func CreateBranchMember(member *models.BranchMember) error {
	member.CreatedOn = time.Now()
	member.UpdatedOn = nil
	if err := config.DB.Create(member).Error; err != nil {
		return err
	}
	return nil
}

// GetAllBranchMembers fetches all branch members
func GetAllBranchMembers() ([]models.BranchMember, error) {
	var members []models.BranchMember
	if err := config.DB.Preload("Branch").Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// GetMembersByBranch fetches members for a specific branch
func GetMembersByBranch(branchID uint) ([]models.BranchMember, error) {
	var members []models.BranchMember
	if err := config.DB.Where("branch_id = ?", branchID).Preload("Branch").Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// UpdateBranchMember updates a member by ID
func UpdateBranchMember(id uint, updatedData map[string]interface{}) error {
	var member models.BranchMember
	if err := config.DB.First(&member, id).Error; err != nil {
		return errors.New("member not found")
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&member).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteBranchMember deletes a member by ID
func DeleteBranchMember(id uint) error {
	if err := config.DB.Delete(&models.BranchMember{}, id).Error; err != nil {
		return err
	}
	return nil
}
