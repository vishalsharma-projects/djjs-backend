package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// Create a new event
func CreateEvent(event *models.EventDetails) error {
	event.CreatedOn = time.Now()
	event.UpdatedOn = nil

	if err := config.DB.Create(event).Error; err != nil {
		return err
	}
	return nil
}

// Get all events with type + category
// statusFilter can be "complete", "incomplete", or empty string for all
func GetAllEvents(statusFilter string) ([]models.EventDetails, error) {
	var events []models.EventDetails

	db := config.DB.
		Preload("EventType").
		Preload("EventCategory")

	// Apply status filter if provided
	if statusFilter != "" {
		db = db.Where("status = ?", statusFilter)
	}

	if err := db.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

// Search events by type, category, or theme
func SearchEvents(search string) ([]models.EventDetails, error) {
	var events []models.EventDetails

	db := config.DB.Preload("EventType").Preload("EventCategory")

	if search != "" {
		db = db.Where(`
			LOWER(theme) LIKE LOWER(?) OR
			LOWER(scale) LIKE LOWER(?)`,
			"%"+search+"%", "%"+search+"%",
		)
	}

	if err := db.Find(&events).Error; err != nil {
		return nil, errors.New("error fetching events")
	}

	if len(events) == 0 {
		return nil, errors.New("no events found")
	}

	return events, nil
}

// Update event
func UpdateEvent(eventID uint, updatedData map[string]interface{}) error {
	var event models.EventDetails

	if err := config.DB.First(&event, eventID).Error; err != nil {
		return errors.New("event not found")
	}

	now := time.Now()
	updatedData["updated_on"] = &now

	if err := config.DB.Model(&event).Updates(updatedData).Error; err != nil {
		return err
	}

	return nil
}

// Delete event and all related data
func DeleteEvent(eventID uint) error {
	// Start a transaction to ensure all deletions succeed or none do
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Delete all special guests for this event
	if err := tx.Where("event_id = ?", eventID).Delete(&models.SpecialGuest{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete special guests: " + err.Error())
	}

	// Delete all volunteers for this event
	if err := tx.Where("event_id = ?", eventID).Delete(&models.Volunteer{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete volunteers: " + err.Error())
	}

	// Delete all event media for this event
	if err := tx.Where("event_id = ?", eventID).Delete(&models.EventMedia{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete event media: " + err.Error())
	}

	// Delete all donations for this event
	if err := tx.Where("event_id = ?", eventID).Delete(&models.Donation{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete donations: " + err.Error())
	}

	// Delete all promotion material details for this event
	if err := tx.Where("event_id = ?", eventID).Delete(&models.PromotionMaterialDetails{}).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete promotion materials: " + err.Error())
	}

	// Delete the event itself
	if err := tx.Delete(&models.EventDetails{}, eventID).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to delete event: " + err.Error())
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	return nil
}

// GetEventByID retrieves an event by ID with all related data
func GetEventByID(eventID uint) (*models.EventDetails, error) {
	var event models.EventDetails

	if err := config.DB.
		Preload("EventType").
		Preload("EventCategory").
		First(&event, eventID).Error; err != nil {
		return nil, errors.New("event not found")
	}

	return &event, nil
}

// UpdateEventStatus updates the status of an event
func UpdateEventStatus(eventID uint, status string) error {
	var event models.EventDetails

	if err := config.DB.First(&event, eventID).Error; err != nil {
		return errors.New("event not found")
	}

	now := time.Now()
	updateData := map[string]interface{}{
		"status":     status,
		"updated_on": &now,
	}

	if err := config.DB.Model(&event).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}
