package services

import (
	"errors"
	"strings"
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

	// Check branch_code uniqueness if provided
	if branch.BranchCode != "" {
		var existingBranch models.Branch
		if err := config.DB.Where("branch_code = ?", branch.BranchCode).First(&existingBranch).Error; err == nil {
			return errors.New("branch code already exists")
		}
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
	
	// Ensure status is set to true when creating a branch
	// If status is not explicitly set, default to true
	if !branch.Status {
		branch.Status = true
	}

	if err := config.DB.Create(branch).Error; err != nil {
		return err
	}
	return nil
}

// GetAllBranches fetches all parent branches only (branches with parent_branch_id IS NULL)
// Child branches are stored in the same table but should only be shown when expanding parent branches
func GetAllBranches() ([]models.Branch, error) {
	var branches []models.Branch
	if err := config.DB.
		Select("id", "name", "email", "coordinator_name", "contact_number", "established_on", "aashram_area",
			"country_id", "state_id", "district_id", "city_id", "parent_branch_id",
			"address", "pincode", "post_office", "police_station", "open_days",
			"daily_start_time", "daily_end_time", "status", "ncr", "region_id", "branch_code",
			"created_on", "updated_on", "created_by", "updated_by").
		Where("parent_branch_id IS NULL"). // Only return parent branches
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City").
		Preload("Children"). // Preload child branches for expand functionality
		Order("id DESC"). // Order by ID descending to show newest first
		Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}

// GetBranch fetches a branch by ID
func GetBranch(branchID uint) (*models.Branch, error) {
	var branch models.Branch
	if err := config.DB.
		Select("id", "name", "email", "coordinator_name", "contact_number", "established_on", "aashram_area",
			"country_id", "state_id", "district_id", "city_id", "parent_branch_id",
			"address", "pincode", "post_office", "police_station", "open_days",
			"daily_start_time", "daily_end_time", "status", "ncr", "region_id", "branch_code",
			"created_on", "updated_on", "created_by", "updated_by").
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City").
		Preload("Parent").
		Preload("Children").
		Preload("Infrastructures").
		Preload("Members").
		First(&branch, branchID).Error; err != nil {
		return nil, errors.New("branch not found")
	}
	return &branch, nil
}

// GetChildBranches fetches all child branches of a parent branch
func GetChildBranches(parentBranchID uint) ([]models.Branch, error) {
	var branches []models.Branch
	if err := config.DB.
		Select("id", "name", "email", "coordinator_name", "contact_number", "established_on", "aashram_area",
			"country_id", "state_id", "district_id", "city_id", "parent_branch_id",
			"address", "pincode", "post_office", "police_station", "open_days",
			"daily_start_time", "daily_end_time", "status", "ncr", "region_id", "branch_code",
			"created_on", "updated_on", "created_by", "updated_by").
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City").
		Where("parent_branch_id = ?", parentBranchID).
		Order("id DESC").
		Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}

// GetBranchSearch fetches parent branches by name and/or coordinator name
// Only returns parent branches (parent_branch_id IS NULL) to match GetAllBranches behavior
func GetBranchSearch(branchName, coordinator string) ([]models.Branch, error) {
	var branches []models.Branch
	db := config.DB.
		Select("id", "name", "email", "coordinator_name", "contact_number", "established_on", "aashram_area",
			"country_id", "state_id", "district_id", "city_id", "parent_branch_id",
			"address", "pincode", "post_office", "police_station", "open_days",
			"daily_start_time", "daily_end_time", "status", "ncr", "region_id", "branch_code",
			"created_on", "updated_on", "created_by", "updated_by").
		Where("parent_branch_id IS NULL"). // Only search parent branches
		Preload("Country").
		Preload("State").
		Preload("District").
		Preload("City")

	// Apply filters dynamically - use OR logic if both are provided
	if branchName != "" && coordinator != "" {
		// If both are provided, search for branches matching name OR coordinator
		db = db.Where("LOWER(name) LIKE LOWER(?) OR LOWER(coordinator_name) LIKE LOWER(?)",
			"%"+branchName+"%", "%"+coordinator+"%")
	} else if branchName != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?)", "%"+branchName+"%")
	} else if coordinator != "" {
		db = db.Where("LOWER(coordinator_name) LIKE LOWER(?)", "%"+coordinator+"%")
	}

	// Order by ID descending to show newest first
	db = db.Order("id DESC")

	if err := db.Find(&branches).Error; err != nil {
		return nil, errors.New("error fetching branches")
	}

	// Return empty array instead of error when no branches found
	// This is a valid state, not an error
	return branches, nil
}

// UpdateBranch updates branch fields
func UpdateBranch(branchID uint, updatedData map[string]interface{}) error {
	var branch models.Branch
	if err := config.DB.First(&branch, branchID).Error; err != nil {
		return errors.New("branch not found")
	}

	// Check email uniqueness if email is being updated (skip if empty or nil)
	if email, ok := updatedData["email"]; ok && email != nil {
		if emailStr, ok := email.(string); ok && strings.TrimSpace(emailStr) != "" {
			var existingBranch models.Branch
			if err := config.DB.Where("email = ? AND id != ?", emailStr, branchID).First(&existingBranch).Error; err == nil {
				return errors.New("email already exists")
			}
		}
	}

	// Check contact number uniqueness if being updated (skip if empty or nil)
	if contactNumber, ok := updatedData["contact_number"]; ok && contactNumber != nil {
		if contactStr, ok := contactNumber.(string); ok && strings.TrimSpace(contactStr) != "" {
			var existingBranch models.Branch
			if err := config.DB.Where("contact_number = ? AND id != ?", contactStr, branchID).First(&existingBranch).Error; err == nil {
				return errors.New("contact number already exists")
			}
		}
	}

	// Check branch_code uniqueness if being updated
	if branchCode, ok := updatedData["branch_code"]; ok {
		if branchCodeStr, ok := branchCode.(string); ok && branchCodeStr != "" {
			var existingBranch models.Branch
			if err := config.DB.Where("branch_code = ? AND id != ?", branchCodeStr, branchID).First(&existingBranch).Error; err == nil {
				return errors.New("branch code already exists")
			}
		}
	}

	// Validate Country ID if being updated
	if countryID, ok := updatedData["country_id"]; ok {
		// Allow nil to clear the country_id
		if countryID == nil {
			// Set to nil to clear the relationship
			updatedData["country_id"] = nil
		} else {
			var countryIDVal uint
			switch v := countryID.(type) {
			case float64:
				countryIDVal = uint(v)
			case uint:
				countryIDVal = v
			case int:
				countryIDVal = uint(v)
			case *uint:
				if v == nil {
					updatedData["country_id"] = nil
					countryIDVal = 0
				} else {
					countryIDVal = *v
				}
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
	}

	// Validate State ID if being updated
	if stateID, ok := updatedData["state_id"]; ok {
		// Allow nil to clear the state_id
		if stateID == nil {
			// Set to nil to clear the relationship
			updatedData["state_id"] = nil
		} else {
			var stateIDVal uint
			switch v := stateID.(type) {
			case float64:
				stateIDVal = uint(v)
			case uint:
				stateIDVal = v
			case int:
				stateIDVal = uint(v)
			case *uint:
				if v == nil {
					updatedData["state_id"] = nil
					stateIDVal = 0
				} else {
					stateIDVal = *v
				}
			default:
				return errors.New("invalid state_id type")
			}
			if stateIDVal > 0 {
				var state models.State
				if err := config.DB.First(&state, stateIDVal).Error; err != nil {
					return errors.New("invalid state_id")
				}
				// Validate state belongs to country if country_id is also being updated or already set
				if countryID, ok := updatedData["country_id"]; ok && countryID != nil {
					var countryIDVal uint
					switch v := countryID.(type) {
					case float64:
						countryIDVal = uint(v)
					case uint:
						countryIDVal = v
					case int:
						countryIDVal = uint(v)
					case *uint:
						if v != nil {
							countryIDVal = *v
						}
					}
					if countryIDVal > 0 && state.CountryID != countryIDVal {
						return errors.New("state does not belong to the specified country")
					}
				} else if branch.CountryID != nil && *branch.CountryID > 0 && state.CountryID != *branch.CountryID {
					return errors.New("state does not belong to the branch's country")
				}
			}
		}
	}

	// Validate District ID if being updated
	if districtID, ok := updatedData["district_id"]; ok {
		// Allow nil to clear the district_id
		if districtID == nil {
			// Set to nil to clear the relationship
			updatedData["district_id"] = nil
		} else {
			var districtIDVal uint
			switch v := districtID.(type) {
			case float64:
				districtIDVal = uint(v)
			case uint:
				districtIDVal = v
			case int:
				districtIDVal = uint(v)
			case *uint:
				if v == nil {
					updatedData["district_id"] = nil
					districtIDVal = 0
				} else {
					districtIDVal = *v
				}
			default:
				return errors.New("invalid district_id type")
			}
			if districtIDVal > 0 {
				var district models.District
				if err := config.DB.First(&district, districtIDVal).Error; err != nil {
					return errors.New("invalid district_id")
				}
				// Validate district belongs to state if state_id is also being updated or already set
				if stateID, ok := updatedData["state_id"]; ok && stateID != nil {
					var stateIDVal uint
					switch v := stateID.(type) {
					case float64:
						stateIDVal = uint(v)
					case uint:
						stateIDVal = v
					case int:
						stateIDVal = uint(v)
					case *uint:
						if v != nil {
							stateIDVal = *v
						}
					}
					if stateIDVal > 0 && district.StateID != stateIDVal {
						return errors.New("district does not belong to the specified state")
					}
				} else if branch.StateID != nil && *branch.StateID > 0 && district.StateID != *branch.StateID {
					return errors.New("district does not belong to the branch's state")
				}
				// Validate district belongs to country
				if countryID, ok := updatedData["country_id"]; ok && countryID != nil {
					var countryIDVal uint
					switch v := countryID.(type) {
					case float64:
						countryIDVal = uint(v)
					case uint:
						countryIDVal = v
					case int:
						countryIDVal = uint(v)
					case *uint:
						if v != nil {
							countryIDVal = *v
						}
					}
					if countryIDVal > 0 && district.CountryID != countryIDVal {
						return errors.New("district does not belong to the specified country")
					}
				} else if branch.CountryID != nil && *branch.CountryID > 0 && district.CountryID != *branch.CountryID {
					return errors.New("district does not belong to the branch's country")
				}
			}
		}
	}

	// Validate City ID if being updated
	if cityID, ok := updatedData["city_id"]; ok {
		// Allow nil to clear the city_id
		if cityID == nil {
			// Set to nil to clear the relationship
			updatedData["city_id"] = nil
		} else {
			var cityIDVal uint
			switch v := cityID.(type) {
			case float64:
				cityIDVal = uint(v)
			case uint:
				cityIDVal = v
			case int:
				cityIDVal = uint(v)
			case *uint:
				if v == nil {
					updatedData["city_id"] = nil
					cityIDVal = 0
				} else {
					cityIDVal = *v
				}
			default:
				return errors.New("invalid city_id type")
			}
			if cityIDVal > 0 {
				var city models.City
				if err := config.DB.First(&city, cityIDVal).Error; err != nil {
					return errors.New("invalid city_id")
				}
				// Validate city belongs to state if state_id is also being updated or already set
				if stateID, ok := updatedData["state_id"]; ok && stateID != nil {
					var stateIDVal uint
					switch v := stateID.(type) {
					case float64:
						stateIDVal = uint(v)
					case uint:
						stateIDVal = v
					case int:
						stateIDVal = uint(v)
					case *uint:
						if v != nil {
							stateIDVal = *v
						}
					}
					if stateIDVal > 0 && city.StateID != stateIDVal {
						return errors.New("city does not belong to the specified state")
					}
				} else if branch.StateID != nil && *branch.StateID > 0 && city.StateID != *branch.StateID {
					return errors.New("city does not belong to the branch's state")
				}
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

// GetBranchMembersWithFilters fetches branch members with optional filters
func GetBranchMembersWithFilters(search string, memberType string, branchType string) ([]models.BranchMember, error) {
	var members []models.BranchMember
	db := config.DB.Preload("Branch")

	// Apply search filter (searches in name, branch_role, responsibility, branch name)
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		db = db.Where(`
			LOWER(name) LIKE ? OR 
			LOWER(branch_role) LIKE ? OR 
			LOWER(responsibility) LIKE ? OR
			EXISTS (
				SELECT 1 FROM branches 
				WHERE branches.id = branch_member.branch_id 
				AND LOWER(branches.name) LIKE ?
			)
		`, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply member type filter
	if memberType != "" && memberType != "all" {
		db = db.Where("LOWER(member_type) = LOWER(?)", memberType)
	}

	// Apply branch type filter
	if branchType != "" && branchType != "all" {
		if branchType == "branch" {
			// Only parent branches (parent_branch_id IS NULL)
			db = db.Where("EXISTS (SELECT 1 FROM branches WHERE branches.id = branch_member.branch_id AND branches.parent_branch_id IS NULL)")
		} else if branchType == "child_branch" {
			// Only child branches (parent_branch_id IS NOT NULL)
			db = db.Where("EXISTS (SELECT 1 FROM branches WHERE branches.id = branch_member.branch_id AND branches.parent_branch_id IS NOT NULL)")
		}
	}

	if err := db.Find(&members).Error; err != nil {
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

	// Parse date_of_birth if provided as string
	if dob, ok := updatedData["date_of_birth"]; ok && dob != nil {
		if dobStr, ok := dob.(string); ok && dobStr != "" {
			if parsedTime, err := time.Parse("2006-01-02", dobStr); err == nil {
				updatedData["date_of_birth"] = &parsedTime
			} else if parsedTime, err := time.Parse(time.RFC3339, dobStr); err == nil {
				updatedData["date_of_birth"] = &parsedTime
			} else {
				// If parsing fails, remove it to avoid database error
				delete(updatedData, "date_of_birth")
			}
		}
	}

	// Parse date_of_samarpan if provided as string
	if dos, ok := updatedData["date_of_samarpan"]; ok && dos != nil {
		if dosStr, ok := dos.(string); ok && dosStr != "" {
			if parsedTime, err := time.Parse("2006-01-02", dosStr); err == nil {
				updatedData["date_of_samarpan"] = &parsedTime
			} else if parsedTime, err := time.Parse(time.RFC3339, dosStr); err == nil {
				updatedData["date_of_samarpan"] = &parsedTime
			} else {
				// If parsing fails, remove it to avoid database error
				delete(updatedData, "date_of_samarpan")
			}
		}
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
