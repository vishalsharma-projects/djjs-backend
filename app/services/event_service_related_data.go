package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// CreateEventRelatedData creates related data for an event (media, guests, volunteers, donations)
func CreateEventRelatedData(eventID uint, payload struct {
	GeneralDetails       map[string]interface{} `json:"generalDetails"`
	MediaPromotion       map[string]interface{} `json:"mediaPromotion"`
	InvolvedParticipants map[string]interface{} `json:"involvedParticipants"`
	DonationTypes        []interface{}          `json:"donationTypes"`
	MaterialTypes        []interface{}          `json:"materialTypes"`
	SpecialGuests        []interface{}          `json:"specialGuests"`
	Volunteers           []interface{}          `json:"volunteers"`
	UploadedFiles        map[string]interface{} `json:"uploadedFiles"`
	DraftID              *uint                  `json:"draftId,omitempty"`
	Status               string                 `json:"status,omitempty"`
}) error {
	// Create Event Media records
	// Check both "eventMediaList" (from frontend) and "eventMedia" (legacy)
	var eventMediaList []interface{}
	if list, ok := payload.MediaPromotion["eventMediaList"].([]interface{}); ok {
		eventMediaList = list
	} else if list, ok := payload.MediaPromotion["eventMedia"].([]interface{}); ok {
		eventMediaList = list
	}

	if len(eventMediaList) > 0 {
		for _, mediaItem := range eventMediaList {
			if mediaMap, ok := mediaItem.(map[string]interface{}); ok {
				media := models.EventMedia{
					EventID: eventID,
				}

				// Get media coverage type
				if mediaTypeName, ok := mediaMap["mediaCoverageType"].(string); ok && mediaTypeName != "" {
					var mediaType models.MediaCoverageType
					if err := config.DB.Where("media_type = ?", mediaTypeName).First(&mediaType).Error; err == nil {
						media.MediaCoverageTypeID = mediaType.ID
					}
				}

				// Company details
				if val, ok := mediaMap["companyName"].(string); ok {
					media.CompanyName = val
				}
				if val, ok := mediaMap["companyEmail"].(string); ok {
					media.CompanyEmail = val
				}
				if val, ok := mediaMap["companyWebsite"].(string); ok {
					media.CompanyWebsite = val
				}

				// Media person details
				if mediaPerson, ok := mediaMap["mediaPerson"].(map[string]interface{}); ok {
					if val, ok := mediaPerson["gender"].(string); ok {
						media.Gender = val
					}
					if val, ok := mediaPerson["prefix"].(string); ok {
						media.Prefix = val
					}
					if val, ok := mediaPerson["firstName"].(string); ok {
						media.FirstName = val
					}
					if val, ok := mediaPerson["middleName"].(string); ok {
						media.MiddleName = val
					}
					if val, ok := mediaPerson["lastName"].(string); ok {
						media.LastName = val
					}
					if val, ok := mediaPerson["designation"].(string); ok {
						media.Designation = val
					}
					if val, ok := mediaPerson["contact"].(string); ok {
						media.Contact = val
					}
					if val, ok := mediaPerson["email"].(string); ok {
						media.Email = val
					}
				} else {
					// Fallback: check direct fields (for backward compatibility)
					if val, ok := mediaMap["gender"].(string); ok {
						media.Gender = val
					}
					if val, ok := mediaMap["prefix"].(string); ok {
						media.Prefix = val
					}
					if val, ok := mediaMap["firstName"].(string); ok {
						media.FirstName = val
					}
					if val, ok := mediaMap["middleName"].(string); ok {
						media.MiddleName = val
					}
					if val, ok := mediaMap["lastName"].(string); ok {
						media.LastName = val
					}
					if val, ok := mediaMap["designation"].(string); ok {
						media.Designation = val
					}
					if val, ok := mediaMap["contact"].(string); ok {
						media.Contact = val
					}
					if val, ok := mediaMap["email"].(string); ok {
						media.Email = val
					}
				}

				if media.CompanyName != "" && media.FirstName != "" && media.LastName != "" {
					if err := config.DB.Create(&media).Error; err != nil {
						// Log error but continue with other records
						fmt.Printf("Error creating event media: %v\n", err)
					}
				}
			}
		}
	}

	// Create Promotion Material Details
	if materialTypes, ok := payload.MediaPromotion["promotionalMaterials"].([]interface{}); ok {
		for _, materialItem := range materialTypes {
			if materialMap, ok := materialItem.(map[string]interface{}); ok {
				material := models.PromotionMaterialDetails{
					EventID: eventID,
				}

				// Get promotion material type
				if materialTypeName, ok := materialMap["materialType"].(string); ok && materialTypeName != "" {
					var promoType models.PromotionMaterial
					if err := config.DB.Where("material_type = ?", materialTypeName).First(&promoType).Error; err == nil {
						material.PromotionMaterialID = promoType.ID
					}
				}

				// Quantity
				if val, ok := materialMap["quantity"].(string); ok {
					if qty, err := strconv.Atoi(val); err == nil {
						material.Quantity = qty
					}
				} else if val, ok := materialMap["quantity"].(float64); ok {
					material.Quantity = int(val)
				}

				// Size
				if val, ok := materialMap["size"].(string); ok {
					material.Size = val
				}

				// Custom dimensions
				if val, ok := materialMap["customHeight"].(float64); ok {
					material.DimensionHeight = val
				}
				if val, ok := materialMap["customWidth"].(float64); ok {
					material.DimensionWidth = val
				}

				if material.PromotionMaterialID > 0 && material.Quantity > 0 {
					if err := config.DB.Create(&material).Error; err != nil {
						fmt.Printf("Error creating promotion material: %v\n", err)
					}
				}
			}
		}
	}

	// Create Special Guests
	for _, guestItem := range payload.SpecialGuests {
		if guestMap, ok := guestItem.(map[string]interface{}); ok {
			guest := models.SpecialGuest{
				EventID: eventID,
			}

			if val, ok := guestMap["gender"].(string); ok {
				guest.Gender = val
			}
			if val, ok := guestMap["prefix"].(string); ok {
				guest.Prefix = val
			}
			if val, ok := guestMap["firstName"].(string); ok {
				guest.FirstName = val
			}
			if val, ok := guestMap["middleName"].(string); ok {
				guest.MiddleName = val
			}
			if val, ok := guestMap["lastName"].(string); ok {
				guest.LastName = val
			}
			if val, ok := guestMap["designation"].(string); ok {
				guest.Designation = val
			}
			if val, ok := guestMap["organization"].(string); ok {
				guest.Organization = val
			}
			if val, ok := guestMap["email"].(string); ok {
				guest.Email = val
			}
			if val, ok := guestMap["city"].(string); ok {
				guest.City = val
			}
			if val, ok := guestMap["state"].(string); ok {
				guest.State = val
			}
			if val, ok := guestMap["personalNumber"].(string); ok {
				guest.PersonalNumber = val
			}
			if val, ok := guestMap["contactPerson"].(string); ok {
				guest.ContactPerson = val
			}
			if val, ok := guestMap["contactPersonNumber"].(string); ok {
				guest.ContactPersonNumber = val
			}
			if val, ok := guestMap["referenceBranchId"].(string); ok {
				guest.ReferenceBranchID = val
			}
			if val, ok := guestMap["referenceVolunteerId"].(string); ok {
				guest.ReferenceVolunteerID = val
			}
			if val, ok := guestMap["referencePersonName"].(string); ok {
				guest.ReferencePersonName = val
			}

			if guest.Prefix != "" {
				if err := config.DB.Create(&guest).Error; err != nil {
					fmt.Printf("Error creating special guest: %v\n", err)
				}
			}
		}
	}

	// Create Volunteers
	for _, volunteerItem := range payload.Volunteers {
		if volMap, ok := volunteerItem.(map[string]interface{}); ok {
			volunteer := models.Volunteer{
				EventID: eventID,
			}

			// Branch ID - try to parse from string or number
			if val, ok := volMap["branchId"].(string); ok && val != "" {
				if branchID, err := strconv.ParseUint(val, 10, 64); err == nil {
					volunteer.BranchID = uint(branchID)
				}
			} else if val, ok := volMap["branchId"].(float64); ok {
				volunteer.BranchID = uint(val)
			}

			if val, ok := volMap["name"].(string); ok {
				volunteer.VolunteerName = val
			}
			if val, ok := volMap["contact"].(string); ok {
				volunteer.Contact = val
			}
			if val, ok := volMap["days"].(float64); ok {
				volunteer.NumberOfDays = int(val)
			}
			if val, ok := volMap["seva"].(string); ok {
				volunteer.SevaInvolved = val
			}
			if val, ok := volMap["mentionSeva"].(string); ok {
				volunteer.MentionSeva = val
			}

			if volunteer.BranchID > 0 && volunteer.VolunteerName != "" {
				if err := config.DB.Create(&volunteer).Error; err != nil {
					fmt.Printf("Error creating volunteer: %v\n", err)
				}
			}
		}
	}

	// Create Donations
	// Handle donations from generalDetails.donations or donationTypes
	var donations []interface{}
	if len(payload.DonationTypes) > 0 {
		donations = payload.DonationTypes
	} else if donationList, ok := payload.GeneralDetails["donations"].([]interface{}); ok {
		donations = donationList
	}

	for _, donationItem := range donations {
		if donationMap, ok := donationItem.(map[string]interface{}); ok {
			donation := models.Donation{
				EventID: eventID,
			}

			// Get branch ID from generalDetails or use a default (you may need to adjust this)
			// For now, we'll set it to 0 or get from context if available
			donation.BranchID = 0 // TODO: Get from context or form data

			if val, ok := donationMap["type"].(string); ok {
				donation.DonationType = val
			}

			if donation.DonationType == "cash" {
				if val, ok := donationMap["amount"].(float64); ok {
					donation.Amount = val
				}
			} else if donation.DonationType == "in-kind" {
				// Store tags as JSON in KindType
				if tags, ok := donationMap["tags"].([]interface{}); ok {
					if tagsJSON, err := json.Marshal(tags); err == nil {
						donation.KindType = string(tagsJSON)
					}
				}
				if val, ok := donationMap["materialValue"].(float64); ok {
					donation.Amount = val
				}
			}

			if donation.DonationType != "" {
				if err := config.DB.Create(&donation).Error; err != nil {
					fmt.Printf("Error creating donation: %v\n", err)
				}
			}
		}
	}

	// Process uploaded files from frontend
	// Note: Files are uploaded to S3 via separate API call after event creation
	// This section can be used to create EventMedia records for files uploaded during event creation
	// The uploadedFiles map contains file references that will be processed by the frontend
	// after event creation, or we can process them here if needed

	return nil
}
