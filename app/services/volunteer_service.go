package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"gorm.io/gorm"
)

var ErrVolunteerNotFound = errors.New("volunteer not found")

// CreateVolunteer persists a new volunteer record
func CreateVolunteer(volunteer *models.Volunteer) error {
	// Validate that branch exists
	var branch models.Branch
	if err := config.DB.First(&branch, volunteer.BranchID).Error; err != nil {
		return errors.New("invalid branch_id: branch does not exist")
	}

	// Validate that event exists
	var event models.Event
	if err := config.DB.First(&event, volunteer.EventID).Error; err != nil {
		return errors.New("invalid event_id: event does not exist")
	}

	now := time.Now()
	volunteer.CreatedOn = now
	volunteer.UpdatedOn = nil

	if err := config.DB.Create(volunteer).Error; err != nil {
		return err
	}
	return nil
}

// GetAllVolunteers returns all volunteers
func GetAllVolunteers() ([]models.Volunteer, error) {
	var volunteers []models.Volunteer
	if err := config.DB.Preload("Branch").Find(&volunteers).Error; err != nil {
		return nil, err
	}
	return volunteers, nil
}

// GetVolunteerByEventID fetches all volunteers for a given eventID
func GetVolunteerByEventID(eventID uint) ([]models.Volunteer, error) {
	var volunteers []models.Volunteer

	if err := config.DB.Where("event_id = ?", eventID).Preload("Branch").Preload("Event").Find(&volunteers).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVolunteerNotFound
		}
		return nil, err
	}

	if len(volunteers) == 0 {
		return nil, ErrVolunteerNotFound
	}

	return volunteers, nil
}

// UpdateVolunteer updates the provided fields on a volunteer
func UpdateVolunteer(id uint, updates map[string]interface{}) error {
	var volunteer models.Volunteer
	if err := config.DB.First(&volunteer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVolunteerNotFound
		}
		return err
	}

	now := time.Now()
	updates["updated_on"] = &now

	if err := config.DB.Model(&volunteer).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// DeleteVolunteer removes a volunteer record
func DeleteVolunteer(id uint) error {
	result := config.DB.Delete(&models.Volunteer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrVolunteerNotFound
	}
	return nil
}
