package services

import (
	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// ===================== Services =====================

func GetAllEventTypesService() ([]models.EventType, error) {
	var list []models.EventType
	if err := config.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func GetAllEventCategoriesService() ([]models.EventCategory, error) {
	var list []models.EventCategory
	if err := config.DB.Preload("EventType").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetAllCountriesService returns all countries
func GetAllCountriesService() ([]models.Country, error) {
	var countries []models.Country
	if err := config.DB.Find(&countries).Error; err != nil {
		return nil, err
	}
	return countries, nil
}

// GetAllStatesService returns all states
func GetAllStatesService() ([]models.State, error) {
	var states []models.State
	if err := config.DB.Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

// GetStatesByCountryService returns states filtered by country ID
func GetStatesByCountryService(countryID uint) ([]models.State, error) {
	var states []models.State
	if err := config.DB.Where("country_id = ?", countryID).Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

// GetAllCitiesService returns all cities
func GetAllCitiesService() ([]models.City, error) {
	var cities []models.City
	if err := config.DB.Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

// GetCitiesByStateService returns cities filtered by state only
func GetCitiesByStateService(stateID uint) ([]models.City, error) {
	var cities []models.City
	db := config.DB

	if stateID != 0 {
		db = db.Where("state_id = ?", stateID)
	}

	if err := db.Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

// GetAllDistricts fetches all districts without filter
func GetAllDistricts() ([]models.District, error) {
	var districts []models.District

	if err := config.DB.Find(&districts).Error; err != nil {
		return nil, err
	}

	return districts, nil
}

// GetDistrictsByStateCountry fetches districts filtered by state and/or country
func GetDistrictsByStateCountry(stateID, countryID uint) ([]models.District, error) {
	var districts []models.District

	db := config.DB

	if stateID != 0 {
		db = db.Where("state_id = ?", stateID)
	}
	if countryID != 0 {
		db = db.Where("country_id = ?", countryID)
	}

	if err := db.Find(&districts).Error; err != nil {
		return nil, err
	}

	return districts, nil
}

func GetAllPromotionMaterialTypesService() ([]models.PromotionMaterial, error) {
	var list []models.PromotionMaterial
	if err := config.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetCoordinatorsDropdownService
func GetCoordinatorDropdownService() ([]models.BranchMember, error) {
	var list []models.BranchMember

	err := config.DB.
		Model(&models.BranchMember{}).
		Select("id, name").
		Where("branch_role = ?", "Coordinator").
		Order("name ASC").
		Find(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetOratorDropdownService
func GetOratorDropdownService() ([]models.Orator, error) {
	var list []models.Orator

	err := config.DB.
		Order("name ASC").
		Find(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetAllLanguagesService returns all languages
func GetAllLanguagesService() ([]models.Language, error) {
	var languages []models.Language
	if err := config.DB.Order("name ASC").Find(&languages).Error; err != nil {
		return nil, err
	}
	return languages, nil
}

// GetAllSevaTypesService returns all seva types
func GetAllSevaTypesService() ([]models.SevaType, error) {
	var sevaTypes []models.SevaType
	if err := config.DB.Order("name ASC").Find(&sevaTypes).Error; err != nil {
		return nil, err
	}
	return sevaTypes, nil
}

// GetAllPrefixesService returns all prefixes
func GetAllPrefixesService() ([]models.Prefix, error) {
	var prefixes []models.Prefix
	if err := config.DB.Order("name ASC").Find(&prefixes).Error; err != nil {
		return nil, err
	}
	return prefixes, nil
}

// GetAllEventSubCategoriesService returns all event sub categories
func GetAllEventSubCategoriesService() ([]models.EventSubCategory, error) {
	var subCategories []models.EventSubCategory
	if err := config.DB.Preload("EventCategory").Preload("EventCategory.EventType").Order("name ASC").Find(&subCategories).Error; err != nil {
		return nil, err
	}
	return subCategories, nil
}

// GetEventSubCategoriesByCategoryService returns event sub categories filtered by category ID
func GetEventSubCategoriesByCategoryService(categoryID uint) ([]models.EventSubCategory, error) {
	var subCategories []models.EventSubCategory
	if err := config.DB.Where("event_category_id = ?", categoryID).Preload("EventCategory").Preload("EventCategory.EventType").Order("name ASC").Find(&subCategories).Error; err != nil {
		return nil, err
	}
	return subCategories, nil
}

// GetAllRolesService returns all roles
func GetAllRolesService() ([]models.Role, error) {
	var roles []models.Role
	if err := config.DB.Order("name ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetAllThemesService returns all themes
func GetAllThemesService() ([]models.Theme, error) {
	var themes []models.Theme
	if err := config.DB.Order("name ASC").Find(&themes).Error; err != nil {
		return nil, err
	}
	return themes, nil
}

// GetAllInfrastructureTypesService returns all infrastructure types
func GetAllInfrastructureTypesService() ([]models.InfrastructureType, error) {
	var types []models.InfrastructureType
	if err := config.DB.Order("name ASC").Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}