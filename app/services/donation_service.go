package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// CreateDonation creates a new donation
func CreateDonation(donation *models.Donation) error {
	donation.CreatedOn = time.Now()

	if err := config.DB.Create(donation).Error; err != nil {
		return err
	}
	return nil
}

// GetAllDonations retrieves all donation entries
func GetAllDonations() ([]models.Donation, error) {
	var donations []models.Donation
	if err := config.DB.Find(&donations).Error; err != nil {
		return nil, err
	}
	return donations, nil
}

// GetDonationsByEvent retrieves donations for a specific event
func GetDonationsByEvent(eventID uint) ([]models.Donation, error) {
	var donations []models.Donation

	if err := config.DB.Where("event_id = ?", eventID).Find(&donations).Error; err != nil {
		return nil, errors.New("error fetching donations")
	}

	return donations, nil
}

// UpdateDonation updates donation fields
func UpdateDonation(id uint, updateData map[string]interface{}) error {
	var donation models.Donation

	if err := config.DB.First(&donation, id).Error; err != nil {
		return errors.New("donation not found")
	}

	now := time.Now()
	updateData["updated_on"] = &now

	if err := config.DB.Model(&donation).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

// DeleteDonation deletes a donation
func DeleteDonation(id uint) error {
	if err := config.DB.Delete(&models.Donation{}, id).Error; err != nil {
		return err
	}
	return nil
}
