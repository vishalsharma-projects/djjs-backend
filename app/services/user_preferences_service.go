package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// SaveUserPreference saves or updates a user preference
func SaveUserPreference(userEmail string, preferenceType string, preferenceData map[string]interface{}) (*models.UserPreferences, error) {
	var preference models.UserPreferences

	// Try to find existing preference
	err := config.DB.Where("user_email = ? AND preference_type = ?", userEmail, preferenceType).
		First(&preference).Error

	now := time.Now()

	if err != nil {
		// Create new preference
		preference = models.UserPreferences{
			UserEmail:      userEmail,
			PreferenceType: preferenceType,
			PreferenceData: models.JSONB(preferenceData),
			CreatedOn:     now,
			UpdatedOn:      &now,
		}
		if err := config.DB.Create(&preference).Error; err != nil {
			return nil, err
		}
	} else {
		// Update existing preference
		preference.PreferenceData = models.JSONB(preferenceData)
		preference.UpdatedOn = &now
		if err := config.DB.Save(&preference).Error; err != nil {
			return nil, err
		}
	}

	return &preference, nil
}

// GetUserPreference retrieves a user preference by type
func GetUserPreference(userEmail string, preferenceType string) (*models.UserPreferences, error) {
	var preference models.UserPreferences
	if err := config.DB.Where("user_email = ? AND preference_type = ?", userEmail, preferenceType).
		First(&preference).Error; err != nil {
		return nil, errors.New("preference not found")
	}
	return &preference, nil
}

// DeleteUserPreference deletes a user preference
func DeleteUserPreference(userEmail string, preferenceType string) error {
	if err := config.DB.Where("user_email = ? AND preference_type = ?", userEmail, preferenceType).
		Delete(&models.UserPreferences{}).Error; err != nil {
		return err
	}
	return nil
}

