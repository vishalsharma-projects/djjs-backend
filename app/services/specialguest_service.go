package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"gorm.io/gorm"
)

var ErrSpecialGuestNotFound = errors.New("special guest not found")

// CreateSpecialGuest inserts a new special guest record
func CreateSpecialGuest(sg *models.SpecialGuest) error {
	now := time.Now()
	sg.CreatedOn = now
	sg.UpdatedOn = nil

	if err := config.DB.Create(sg).Error; err != nil {
		return err
	}
	return nil
}

// GetAllSpecialGuests fetches all special guests
func GetAllSpecialGuests() ([]models.SpecialGuest, error) {
	var guests []models.SpecialGuest
	if err := config.DB.Find(&guests).Error; err != nil {
		return nil, err
	}
	return guests, nil
}

// GetSpecialGuestByEventID fetches all special guests for a given eventID
func GetSpecialGuestByEventID(eventID uint) ([]models.SpecialGuest, error) {
	var guests []models.SpecialGuest

	if err := config.DB.Where("event_id = ?", eventID).Find(&guests).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSpecialGuestNotFound
		}
		return nil, err
	}

	return guests, nil
}

// UpdateSpecialGuest updates a special guest by ID
func UpdateSpecialGuest(sgID uint, updatedData map[string]interface{}) error {
	var guest models.SpecialGuest
	if err := config.DB.First(&guest, sgID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSpecialGuestNotFound
		}
		return err
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&guest).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// DeleteSpecialGuest deletes a special guest
func DeleteSpecialGuest(sgID uint) error {
	result := config.DB.Delete(&models.SpecialGuest{}, sgID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrSpecialGuestNotFound
	}
	return nil
}
