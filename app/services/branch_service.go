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

	// Validate Country ID if provided
	if branch.CountryID != nil && *branch.CountryID > 0 {
		var country models.Country
		if err := config.DB.First(&country, *branch.CountryID).Error; err != nil {
			return errors.New("invalid country_id")
		}
	}

	// Validate State ID if provided
	if branch.StateID != nil && *branch.StateID > 0 {
		var state models.State
		if err := config.DB.First(&state, *branch.StateID).Error; err != nil {
			return errors.New("invalid state_id")
		}
		// Validate state belongs to country if country_id is also provided
		if branch.CountryID != nil && *branch.CountryID > 0 && state.CountryID != *branch.CountryID {
			return errors.New("state does not belong to the specified country")
		}
	}

	// Validate District ID if provided
	if branch.DistrictID != nil && *branch.DistrictID > 0 {
		var district models.District
		if err := config.DB.First(&district, *branch.DistrictID).Error; err != nil {
			return errors.New("invalid district_id")
		}
		// Validate district belongs to state if state_id is also provided
		if branch.StateID != nil && *branch.StateID > 0 && district.StateID != *branch.StateID {
			return errors.New("district does not belong to the specified state")
		}
		// Validate district belongs to country if country_id is also provided
		if branch.CountryID != nil && *branch.CountryID > 0 && district.CountryID != *branch.CountryID {
			return errors.New("district does not belong to the specified country")
		}
	}

	// Validate City ID if provided
	if branch.CityID != nil && *branch.CityID > 0 {
		var city models.City
		if err := config.DB.First(&city, *branch.CityID).Error; err != nil {
			return errors.New("invalid city_id")
		}
		// Validate city belongs to state if state_id is also provided
		if branch.StateID != nil && *branch.StateID > 0 && city.StateID != *branch.StateID {
			return errors.New("city does not belong to the specified state")
		}
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
	if err := config.DB.Preload("Country").Preload("State").Preload("District").Preload("City").Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}

// GetBranch fetches a branch by ID
func GetBranch(branchID uint) (*models.Branch, error) {
	var branch models.Branch
	if err := config.DB.Preload("Country").Preload("State").Preload("District").Preload("City").First(&branch, branchID).Error; err != nil {
		return nil, errors.New("branch not found")
	}
	return &branch, nil
}

// GetBranchSearch fetches branches by name and/or coordinator name
func GetBranchSearch(branchName, coordinator string) ([]models.Branch, error) {
	var branches []models.Branch
	db := config.DB.Preload("Country").Preload("State").Preload("District").Preload("City")

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

	// Validate Country ID if being updated
	if countryID, ok := updatedData["country_id"]; ok {
		var countryIDVal uint
		switch v := countryID.(type) {
		case float64:
			countryIDVal = uint(v)
		case uint:
			countryIDVal = v
		case int:
			countryIDVal = uint(v)
		default:
			return errors.New("invalid country_id type")
		}
		if countryIDVal > 0 {
			var country models.Country
			if err := config.DB.First(&country, countryIDVal).Error; err != nil {
				return errors.New("invalid country_id")
			}
		}
	}

	// Validate State ID if being updated
	if stateID, ok := updatedData["state_id"]; ok {
		var stateIDVal uint
		switch v := stateID.(type) {
		case float64:
			stateIDVal = uint(v)
		case uint:
			stateIDVal = v
		case int:
			stateIDVal = uint(v)
		default:
			return errors.New("invalid state_id type")
		}
		if stateIDVal > 0 {
			var state models.State
			if err := config.DB.First(&state, stateIDVal).Error; err != nil {
				return errors.New("invalid state_id")
			}
			// Validate state belongs to country if country_id is also being updated or already set
			if countryID, ok := updatedData["country_id"]; ok {
				var countryIDVal uint
				switch v := countryID.(type) {
				case float64:
					countryIDVal = uint(v)
				case uint:
					countryIDVal = v
				case int:
					countryIDVal = uint(v)
				}
				if countryIDVal > 0 && state.CountryID != countryIDVal {
					return errors.New("state does not belong to the specified country")
				}
			} else if branch.CountryID != nil && *branch.CountryID > 0 && state.CountryID != *branch.CountryID {
				return errors.New("state does not belong to the branch's country")
			}
		}
	}

	// Validate District ID if being updated
	if districtID, ok := updatedData["district_id"]; ok {
		var districtIDVal uint
		switch v := districtID.(type) {
		case float64:
			districtIDVal = uint(v)
		case uint:
			districtIDVal = v
		case int:
			districtIDVal = uint(v)
		default:
			return errors.New("invalid district_id type")
		}
		if districtIDVal > 0 {
			var district models.District
			if err := config.DB.First(&district, districtIDVal).Error; err != nil {
				return errors.New("invalid district_id")
			}
			// Validate district belongs to state if state_id is also being updated or already set
			if stateID, ok := updatedData["state_id"]; ok {
				var stateIDVal uint
				switch v := stateID.(type) {
				case float64:
					stateIDVal = uint(v)
				case uint:
					stateIDVal = v
				case int:
					stateIDVal = uint(v)
				}
				if stateIDVal > 0 && district.StateID != stateIDVal {
					return errors.New("district does not belong to the specified state")
				}
			} else if branch.StateID != nil && *branch.StateID > 0 && district.StateID != *branch.StateID {
				return errors.New("district does not belong to the branch's state")
			}
			// Validate district belongs to country
			if countryID, ok := updatedData["country_id"]; ok {
				var countryIDVal uint
				switch v := countryID.(type) {
				case float64:
					countryIDVal = uint(v)
				case uint:
					countryIDVal = v
				case int:
					countryIDVal = uint(v)
				}
				if countryIDVal > 0 && district.CountryID != countryIDVal {
					return errors.New("district does not belong to the specified country")
				}
			} else if branch.CountryID != nil && *branch.CountryID > 0 && district.CountryID != *branch.CountryID {
				return errors.New("district does not belong to the branch's country")
			}
		}
	}

	// Validate City ID if being updated
	if cityID, ok := updatedData["city_id"]; ok {
		var cityIDVal uint
		switch v := cityID.(type) {
		case float64:
			cityIDVal = uint(v)
		case uint:
			cityIDVal = v
		case int:
			cityIDVal = uint(v)
		default:
			return errors.New("invalid city_id type")
		}
		if cityIDVal > 0 {
			var city models.City
			if err := config.DB.First(&city, cityIDVal).Error; err != nil {
				return errors.New("invalid city_id")
			}
			// Validate city belongs to state if state_id is also being updated or already set
			if stateID, ok := updatedData["state_id"]; ok {
				var stateIDVal uint
				switch v := stateID.(type) {
				case float64:
					stateIDVal = uint(v)
				case uint:
					stateIDVal = v
				case int:
					stateIDVal = uint(v)
				}
				if stateIDVal > 0 && city.StateID != stateIDVal {
					return errors.New("city does not belong to the specified state")
				}
			} else if branch.StateID != nil && *branch.StateID > 0 && city.StateID != *branch.StateID {
				return errors.New("city does not belong to the branch's state")
			}
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
