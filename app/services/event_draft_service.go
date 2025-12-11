package services

import (
	"errors"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// SaveDraft saves or updates a draft for a specific step
// Returns the draft ID
func SaveDraft(draftID *uint, step string, data map[string]interface{}, userEmail string) (uint, error) {
	var draft models.EventDraft

	if draftID != nil && *draftID > 0 {
		// Update existing draft
		if err := config.DB.First(&draft, *draftID).Error; err != nil {
			return 0, errors.New("draft not found")
		}
	} else {
		// Create new draft
		draft = models.EventDraft{
			UserEmail: userEmail,
			CreatedOn: time.Now(),
		}
	}

	// Update the appropriate step field based on the step name
	switch step {
	case "generalDetails":
		draft.GeneralDetailsDraft = models.JSONB(data)
	case "mediaPromotion":
		draft.MediaPromotionDraft = models.JSONB(data)
	case "specialGuests":
		draft.SpecialGuestsDraft = models.JSONB(data)
	case "volunteers":
		draft.VolunteersDraft = models.JSONB(data)
	default:
		return 0, errors.New("invalid step name")
	}

	now := time.Now()
	draft.UpdatedOn = &now

	if draftID != nil && *draftID > 0 {
		// Update existing draft
		if err := config.DB.Save(&draft).Error; err != nil {
			return 0, err
		}
		return draft.ID, nil
	} else {
		// Create new draft
		if err := config.DB.Create(&draft).Error; err != nil {
			return 0, err
		}
		return draft.ID, nil
	}
}

// GetDraft retrieves a draft by ID
func GetDraft(draftID uint) (*models.EventDraft, error) {
	var draft models.EventDraft
	if err := config.DB.First(&draft, draftID).Error; err != nil {
		return nil, errors.New("draft not found")
	}
	return &draft, nil
}

// GetLatestDraftByUserEmail retrieves the latest draft for a user by email
func GetLatestDraftByUserEmail(userEmail string) (*models.EventDraft, error) {
	var draft models.EventDraft
	if err := config.DB.Where("user_email = ?", userEmail).
		Order("updated_on DESC, created_on DESC").
		First(&draft).Error; err != nil {
		return nil, errors.New("draft not found")
	}
	return &draft, nil
}

// DeleteDraft deletes a draft by ID
func DeleteDraft(draftID uint) error {
	if err := config.DB.Delete(&models.EventDraft{}, draftID).Error; err != nil {
		return errors.New("failed to delete draft")
	}
	return nil
}





