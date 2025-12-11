package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// CreateEventMedia creates a new EventMedia record
func CreateEventMedia(media *models.EventMedia) error {
	return config.DB.Create(media).Error
}

// GetAllEventMedia retrieves all EventMedia records with related Event and MediaCoverageType
func GetAllEventMedia() ([]models.EventMedia, error) {
	var medias []models.EventMedia
	if err := config.DB.
		Preload("Event").
		Preload("MediaCoverageType").
		Find(&medias).Error; err != nil {
		return nil, err
	}
	return medias, nil
}

// GetEventMediaByEventID retrieves all EventMedia records by EventID
func GetEventMediaByEventID(eventID uint) ([]models.EventMedia, error) {
	var mediaList []models.EventMedia
	if err := config.DB.
		Preload("Event").
		Preload("MediaCoverageType").
		Where("event_id = ?", eventID).
		Find(&mediaList).Error; err != nil {
		return nil, errors.New("no event media found for the given event ID")
	}
	return mediaList, nil
}

// UpdateEventMedia updates an existing EventMedia record
func UpdateEventMedia(media *models.EventMedia) error {
	var existing models.EventMedia

	// Check if record exists
	if err := config.DB.First(&existing, media.ID).Error; err != nil {
		return errors.New("record not found")
	}

	// Prepare dynamic updates
	updates := map[string]interface{}{
		"updated_on": time.Now(),
	}

	if media.CompanyName != "" {
		updates["company_name"] = media.CompanyName
	}
	if media.CompanyEmail != "" {
		updates["company_email"] = media.CompanyEmail
	}
	if media.CompanyWebsite != "" {
		updates["company_website"] = media.CompanyWebsite
	}
	if media.Gender != "" {
		updates["gender"] = media.Gender
	}
	if media.Prefix != "" {
		updates["prefix"] = media.Prefix
	}
	if media.FirstName != "" {
		updates["first_name"] = media.FirstName
	}
	if media.MiddleName != "" {
		updates["middle_name"] = media.MiddleName
	}
	if media.LastName != "" {
		updates["last_name"] = media.LastName
	}
	if media.Designation != "" {
		updates["designation"] = media.Designation
	}
	if media.Contact != "" {
		updates["contact"] = media.Contact
	}
	if media.Email != "" {
		updates["email"] = media.Email
	}
	if media.EventID != 0 {
		updates["event_id"] = media.EventID
	}
	if media.MediaCoverageTypeID != 0 {
		updates["media_coverage_type_id"] = media.MediaCoverageTypeID
	}
	if media.FileURL != "" {
		updates["file_url"] = media.FileURL
	}
	if media.FileType != "" {
		updates["file_type"] = media.FileType
	}
	if media.UpdatedBy != "" {
		updates["updated_by"] = media.UpdatedBy
	}

	// Apply updates
	return config.DB.Model(&existing).Updates(updates).Error
}

// DeleteEventMedia deletes an EventMedia record by ID
func DeleteEventMedia(id uint) error {
	result := config.DB.Delete(&models.EventMedia{}, id)
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return result.Error
}
