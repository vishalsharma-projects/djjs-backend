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
func GetOratorDropdownService() ([]models.BranchMember, error) {
	var list []models.BranchMember

	err := config.DB.
		Model(&models.BranchMember{}).
		Select("id, name").
		Where("branch_role IN ?", []string{"Coordinator", "Preacher"}).
		Order("name ASC").
		Find(&list).Error

	if err != nil {
		return nil, err
	}

	return list, nil
}
