package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/app/validators"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
)

// ----------------------------------------------------
// Create Event
// ----------------------------------------------------

// CreateEventHandler godoc
// @Summary Create a new event
// @Description Creates a new event from frontend payload structure. Accepts generalDetails, mediaPromotion, involvedParticipants, donationTypes, materialTypes, specialGuests, volunteers, uploadedFiles, and optional draftId. If draftId is provided, the draft will be automatically deleted from event_drafts table after successful event creation.
// @Tags Events
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param event body object true "Frontend event payload" example({"generalDetails":{"eventType":"Spiritual","scale":"Large (L)","theme":"Devotional"},"mediaPromotion":{},"involvedParticipants":{"beneficiariesMen":50},"donationTypes":[],"materialTypes":[],"specialGuests":[],"volunteers":[],"uploadedFiles":{},"draftId":1})
// @Success 201 {object} map[string]interface{} "Event created successfully" example({"message":"Event created successfully","event":{"id":1,"event_type_id":1,"event_category_id":1}})
// @Failure 400 {object} map[string]string "Bad Request" example({"error":"Invalid event data"})
// @Failure 500 {object} map[string]string "Internal Server Error" example({"error":"Failed to create event"})
// @Router /api/events [post]
func CreateEventHandler(c *gin.Context) {
	// Accept frontend payload structure
	var frontendPayload struct {
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
	}

	// Try to bind frontend payload structure first
	if err := c.ShouldBindJSON(&frontendPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload format: " + err.Error()})
		return
	}

	// Process frontend payload - map to EventDetails with status support
	event, err := services.MapFrontendPayloadToEventWithStatus(frontendPayload.GeneralDetails, frontendPayload.InvolvedParticipants, frontendPayload.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate event
	if err := validators.ValidateEventInput(event.EventTypeID, event.EventCategoryID, event.StartDate, event.EndDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create event in main table
	if err := services.CreateEvent(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	// Create related records (media, special guests, volunteers, donations, etc.)
	if err := services.CreateEventRelatedData(event.ID, frontendPayload); err != nil {
		// Log error but don't fail event creation
		// The event is already created, related data can be added later
		log.Printf("Warning: Failed to create related data: %v", err)
	}

	// Delete draft ONLY after successful event creation with status='complete' (submit)
	// This ensures draft is kept if user just saves as draft, and deleted only when submitting
	if frontendPayload.DraftID != nil && *frontendPayload.DraftID > 0 && frontendPayload.Status == "complete" {
		_ = services.DeleteDraft(*frontendPayload.DraftID)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Event created successfully",
		"event":   event,
	})
}

// ----------------------------------------------------
// Get All Events
// ----------------------------------------------------

// GetAllEventsHandler godoc
// @Summary Get all events
// @Description Get all events, optionally filtered by status (complete/incomplete)
// @Tags Events
// @Security ApiKeyAuth
// @Produce json
// @Param status query string false "Filter by status: complete or incomplete"
// @Success 200 {array} models.EventDetails
// @Failure 500 {object} map[string]string
// @Router /api/events [get]
func GetAllEventsHandler(c *gin.Context) {
	statusFilter := c.Query("status")
	events, err := services.GetAllEvents(statusFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}

	// Add counts for related data to each event
	eventsWithCounts := make([]gin.H, 0, len(events))
	for _, event := range events {
		// Get counts for related data (return empty arrays if not found)
		specialGuests, errSG := services.GetSpecialGuestByEventID(event.ID)
		if errSG != nil {
			specialGuests = []models.SpecialGuest{}
		}

		volunteers, errVol := services.GetVolunteerByEventID(event.ID)
		if errVol != nil {
			volunteers = []models.Volunteer{}
		}

		mediaList, errMedia := services.GetEventMediaByEventID(event.ID)
		if errMedia != nil {
			mediaList = []models.EventMedia{}
		}
		// Convert to presigned URLs - HARD GUARD: fail fast if S3Key is empty
		mediaListWithPresignedURLs, err := services.ConvertEventMediaToPresignedURLs(c.Request.Context(), mediaList)
		if err != nil {
			// Log the error for debugging
			log.Printf("ERROR: Failed to generate presigned URLs for event %d: %v", event.ID, err)
			// Fail fast - return HTTP 500 with structured error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to generate presigned URLs for event media",
				"details": err.Error(),
			})
			return
		}
		mediaList = mediaListWithPresignedURLs

		// Get promotion materials count
		promotionMaterials, errPromo := services.GetPromotionMaterialDetailsByEventID(event.ID)
		if errPromo != nil {
			promotionMaterials = []models.PromotionMaterialDetails{}
		}

		// Get donations count
		donations, errDonations := services.GetDonationsByEvent(event.ID)
		if errDonations != nil {
			donations = []models.Donation{}
		}

		// Get branch from first volunteer or donation
		var branchName string
		var branchID uint
		if len(volunteers) > 0 && volunteers[0].BranchID > 0 {
			// Try to get branch from first volunteer
			var branch models.Branch
			if err := config.DB.First(&branch, volunteers[0].BranchID).Error; err == nil {
				branchName = branch.Name
				branchID = branch.ID
			}
		} else if len(donations) > 0 && donations[0].BranchID > 0 {
			// Fallback to first donation's branch
			var branch models.Branch
			if err := config.DB.First(&branch, donations[0].BranchID).Error; err == nil {
				branchName = branch.Name
				branchID = branch.ID
			}
		}

		// Convert event to map and add counts
		eventMap := gin.H{
			"id":                       event.ID,
			"event_type_id":            event.EventTypeID,
			"event_category_id":        event.EventCategoryID,
			"scale":                    event.Scale,
			"theme":                    event.Theme,
			"start_date":               event.StartDate,
			"end_date":                 event.EndDate,
			"daily_start_time":         event.DailyStartTime,
			"daily_end_time":           event.DailyEndTime,
			"spiritual_orator":         event.SpiritualOrator,
			"language":                 event.Language,
			"branch":                   branchName,
			"branch_id":                branchID,
			"country":                  event.Country,
			"state":                    event.State,
			"city":                     event.City,
			"district":                 event.District,
			"post_office":              event.PostOffice,
			"pincode":                  event.Pincode,
			"address":                  event.Address,
			"beneficiary_men":          event.BeneficiaryMen,
			"beneficiary_women":        event.BeneficiaryWomen,
			"beneficiary_child":        event.BeneficiaryChild,
			"initiation_men":           event.InitiationMen,
			"initiation_women":         event.InitiationWomen,
			"initiation_child":         event.InitiationChild,
			"status":                   event.Status,
			"created_on":               event.CreatedOn,
			"updated_on":               event.UpdatedOn,
			"created_by":               event.CreatedBy,
			"updated_by":               event.UpdatedBy,
			"event_type":               event.EventType,
			"event_category":           event.EventCategory,
			"special_guests_count":     len(specialGuests),
			"volunteers_count":         len(volunteers),
			"media_count":              len(mediaList),
			"promotion_materials_count": len(promotionMaterials),
			"donations_count":          len(donations),
		}
		eventsWithCounts = append(eventsWithCounts, eventMap)
	}

	c.JSON(http.StatusOK, eventsWithCounts)
}

// ----------------------------------------------------
// Get Event By ID
// ----------------------------------------------------

// GetEventByIdHandler godoc
// @Summary Get event by ID
// @Description Get a single event by its ID with related data (special guests, volunteers, media)
// @Tags Events
// @Security ApiKeyAuth
// @Produce json
// @Param event_id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "Event with related data"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events/{event_id} [get]
func GetEventByIdHandler(c *gin.Context) {
	idParam := c.Param("event_id")
	eventID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := services.GetEventByID(uint(eventID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Fetch related data (return empty arrays if not found)
	specialGuests, errSG := services.GetSpecialGuestByEventID(uint(eventID))
	if errSG != nil {
		// Special guests service returns error only on DB error, not on empty result
		specialGuests = []models.SpecialGuest{}
	}

	volunteers, errVol := services.GetVolunteerByEventID(uint(eventID))
	if errVol != nil {
		// Volunteers service returns ErrVolunteerNotFound if empty, treat as empty array
		volunteers = []models.Volunteer{}
	}

	mediaList, errMedia := services.GetEventMediaByEventID(uint(eventID))
	if errMedia != nil {
		// Media service returns error if not found, treat as empty array
		mediaList = []models.EventMedia{}
	}
		// Convert to presigned URLs - HARD GUARD: fail fast if S3Key is empty
		mediaListWithPresignedURLs, err := services.ConvertEventMediaToPresignedURLs(c.Request.Context(), mediaList)
		if err != nil {
			// Fail fast - return HTTP 500 with structured error
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to generate presigned URLs for event media",
				"details": err.Error(),
			})
			return
		}
		mediaList = mediaListWithPresignedURLs

	// Fetch promotion materials
	promotionMaterials, errPromo := services.GetPromotionMaterialDetailsByEventID(uint(eventID))
	if errPromo != nil {
		// Return empty array if not found (consistent with other related data)
		promotionMaterials = []models.PromotionMaterialDetails{}
	}

	// Fetch donations
	donations, errDonations := services.GetDonationsByEvent(uint(eventID))
	if errDonations != nil {
		donations = []models.Donation{}
	}

	// Get branch from first volunteer or donation
	var branchName string
	var branchID uint
	if len(volunteers) > 0 && volunteers[0].BranchID > 0 {
		// Try to get branch from first volunteer
		var branch models.Branch
		if err := config.DB.First(&branch, volunteers[0].BranchID).Error; err == nil {
			branchName = branch.Name
			branchID = branch.ID
		}
	} else if len(donations) > 0 && donations[0].BranchID > 0 {
		// Fallback to first donation's branch
		var branch models.Branch
		if err := config.DB.First(&branch, donations[0].BranchID).Error; err == nil {
			branchName = branch.Name
			branchID = branch.ID
		}
	}

	// Build response with event and related data
	response := gin.H{
		"event":                  event,
		"branch":                 branchName,
		"branch_id":              branchID,
		"specialGuests":          specialGuests,
		"volunteers":             volunteers,
		"media":                  mediaList,
		"promotionMaterials":     promotionMaterials,
		"donations":              donations,
		"specialGuestsCount":     len(specialGuests),
		"volunteersCount":        len(volunteers),
		"mediaCount":             len(mediaList),
		"promotionMaterialsCount": len(promotionMaterials),
		"donationsCount":         len(donations),
	}

	c.JSON(http.StatusOK, response)
}

// ----------------------------------------------------
// Search Events
// ----------------------------------------------------

// SearchEventsHandler godoc
// @Summary Search events
// @Description Search events by keyword
// @Tags Events
// @Security ApiKeyAuth
// @Produce json
// @Param search query string false "Search keyword"
// @Success 200 {array} models.EventDetails
// @Failure 500 {object} map[string]string
// @Router /api/events/search [get]
func SearchEventsHandler(c *gin.Context) {
	search := c.Query("search")

	events, err := services.SearchEvents(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// ----------------------------------------------------
// Update Event
// ----------------------------------------------------

// UpdateEventHandler godoc
// @Summary Update an event
// @Description Updates an event. Accepts both flat structure (for simple updates) and nested frontend payload structure (for full updates with related data)
// @Tags Events
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param event_id path int true "Event ID"
// @Param event body object true "Updated fields (can be flat or nested frontend payload)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events/{event_id} [put]
func UpdateEventHandler(c *gin.Context) {
	idParam := c.Param("event_id")
	eventID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Try to bind as frontend payload structure first
	var frontendPayload struct {
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
	}

	// Check if it's a nested frontend payload
	if err := c.ShouldBindJSON(&frontendPayload); err == nil && frontendPayload.GeneralDetails != nil {
		// It's a nested frontend payload - map to EventDetails and update
		event, err := services.MapFrontendPayloadToEventWithStatus(frontendPayload.GeneralDetails, frontendPayload.InvolvedParticipants, frontendPayload.Status)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Convert event to update map
		updateData := make(map[string]interface{})
		if event.EventTypeID > 0 {
			updateData["event_type_id"] = event.EventTypeID
		}
		if event.EventCategoryID > 0 {
			updateData["event_category_id"] = event.EventCategoryID
		}
		if event.Scale != "" {
			updateData["scale"] = event.Scale
		}
		if event.Theme != "" {
			updateData["theme"] = event.Theme
		}
		if !event.StartDate.IsZero() {
			updateData["start_date"] = event.StartDate
		}
		if !event.EndDate.IsZero() {
			updateData["end_date"] = event.EndDate
		}
		if event.DailyStartTime != nil {
			updateData["daily_start_time"] = event.DailyStartTime
		}
		if event.DailyEndTime != nil {
			updateData["daily_end_time"] = event.DailyEndTime
		}
		if event.SpiritualOrator != "" {
			updateData["spiritual_orator"] = event.SpiritualOrator
		}
		if event.Language != "" {
			updateData["language"] = event.Language
		}
		if event.Country != "" {
			updateData["country"] = event.Country
		}
		if event.State != "" {
			updateData["state"] = event.State
		}
		if event.District != "" {
			updateData["district"] = event.District
		}
		if event.City != "" {
			updateData["city"] = event.City
		}
		if event.Pincode != "" {
			updateData["pincode"] = event.Pincode
		}
		if event.PostOffice != "" {
			updateData["post_office"] = event.PostOffice
		}
		if event.Address != "" {
			updateData["address"] = event.Address
		}
		if event.BeneficiaryMen > 0 {
			updateData["beneficiary_men"] = event.BeneficiaryMen
		}
		if event.BeneficiaryWomen > 0 {
			updateData["beneficiary_women"] = event.BeneficiaryWomen
		}
		if event.BeneficiaryChild > 0 {
			updateData["beneficiary_child"] = event.BeneficiaryChild
		}
		if event.InitiationMen > 0 {
			updateData["initiation_men"] = event.InitiationMen
		}
		if event.InitiationWomen > 0 {
			updateData["initiation_women"] = event.InitiationWomen
		}
		if event.InitiationChild > 0 {
			updateData["initiation_child"] = event.InitiationChild
		}
		if event.Status != "" {
			updateData["status"] = event.Status
		}

		// Validate update fields
		if err := validators.ValidateEventUpdateFields(updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update event
		if err := services.UpdateEvent(uint(eventID), updateData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Update related data if provided
		if err := services.CreateEventRelatedData(uint(eventID), frontendPayload); err != nil {
			log.Printf("Warning: Failed to update related data: %v", err)
		}

		// Delete draft ONLY if status is 'complete' (submit)
		// This ensures draft is kept if user just saves as draft, and deleted only when submitting
		if frontendPayload.DraftID != nil && *frontendPayload.DraftID > 0 && frontendPayload.Status == "complete" {
			_ = services.DeleteDraft(*frontendPayload.DraftID)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
		return
	}

	// Fallback: try as flat structure (for simple updates)
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract draftId and status from flat structure if present
	var draftID *uint
	var status string
	if draftIdVal, ok := updateData["draftId"]; ok && draftIdVal != nil {
		switch v := draftIdVal.(type) {
		case float64:
			id := uint(v)
			draftID = &id
		case string:
			if parsed, err := strconv.ParseUint(v, 10, 64); err == nil {
				id := uint(parsed)
				draftID = &id
			}
		case uint:
			draftID = &v
		case int:
			id := uint(v)
			draftID = &id
		}
		// Remove draftId from updateData as it's not a field in event_details table
		delete(updateData, "draftId")
	}
	if statusVal, ok := updateData["status"].(string); ok {
		status = statusVal
	}

	// Parse date strings to time.Time if present
	if startDateStr, ok := updateData["start_date"].(string); ok && startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			updateData["start_date"] = t
		} else if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			updateData["start_date"] = t
		}
	}
	if endDateStr, ok := updateData["end_date"].(string); ok && endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			updateData["end_date"] = t
		} else if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			updateData["end_date"] = t
		}
	}

	if err := validators.ValidateEventUpdateFields(updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateEvent(uint(eventID), updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete draft ONLY if status is 'complete' (submit)
	// This ensures draft is kept if user just saves as draft, and deleted only when submitting
	if draftID != nil && *draftID > 0 && status == "complete" {
		_ = services.DeleteDraft(*draftID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// ----------------------------------------------------
// Delete Event
// ----------------------------------------------------

// DeleteEventHandler godoc
// @Summary Delete an event
// @Tags Events
// @Security ApiKeyAuth
// @Produce json
// @Param event_id path int true "Event ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events/{event_id} [delete]
func DeleteEventHandler(c *gin.Context) {
	idParam := c.Param("event_id")
	eventID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	if err := services.DeleteEvent(uint(eventID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// ----------------------------------------------------
// Download Event
// ----------------------------------------------------

// DownloadEventHandler godoc
// @Summary Download event data as PDF
// @Description Downloads event data as a PDF document
// @Tags Events
// @Security ApiKeyAuth
// @Produce application/pdf
// @Param event_id path int true "Event ID"
// @Success 200 {file} file "Event data PDF file"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events/{event_id}/download [get]
func DownloadEventHandler(c *gin.Context) {
	idParam := c.Param("event_id")
	eventID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// Get event with all related data
	event, err := services.GetEventByID(uint(eventID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Fetch all related data
	specialGuests, _ := services.GetSpecialGuestByEventID(uint(eventID))
	volunteers, _ := services.GetVolunteerByEventID(uint(eventID))
	mediaList, _ := services.GetEventMediaByEventID(uint(eventID))
	// Convert to presigned URLs - HARD GUARD: fail fast if S3Key is empty
	mediaListWithPresignedURLs, err := services.ConvertEventMediaToPresignedURLs(c.Request.Context(), mediaList)
	if err != nil {
		// Fail fast - return HTTP 500 with structured error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to generate presigned URLs for event media",
			"details": err.Error(),
		})
		return
	}
	mediaList = mediaListWithPresignedURLs
	promotionMaterials, _ := services.GetPromotionMaterialDetailsByEventID(uint(eventID))
	donations, _ := services.GetDonationsByEvent(uint(eventID))

	// Generate PDF document
	pdfBytes, err := services.GenerateEventPDF(event, specialGuests, volunteers, mediaList, promotionMaterials, donations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
		return
	}

	// Set headers for PDF file download
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=event_%d_%s.pdf", eventID, time.Now().Format("20060102_150405")))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// ----------------------------------------------------
// Save Draft
// ----------------------------------------------------

// SaveDraftHandler godoc
// @Summary Save draft data for a specific step
// @Description Saves draft data for event creation. Creates a new draft if draftId is not provided, or updates existing draft. Drafts are stored in a separate event_drafts table and automatically deleted when event is submitted.
// @Tags Events
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param draft body object true "Draft payload" example({"draftId":null,"step":"generalDetails","data":{"eventType":"Spiritual","eventName":"Bhagwat Katha","scale":"Large (L)"}})
// @Success 200 {object} map[string]interface{} "Draft saved successfully" example({"draftId":1,"message":"Draft saved successfully"})
// @Failure 400 {object} map[string]string "Bad Request" example({"error":"Invalid step name. Must be one of: generalDetails, mediaPromotion, specialGuests, volunteers, donations"})
// @Failure 500 {object} map[string]string "Internal Server Error" example({"error":"Failed to save draft"})
// @Router /api/events/draft [post]
func SaveDraftHandler(c *gin.Context) {
	var draftRequest struct {
		DraftID interface{} `json:"draftId"` // Changed from eventId to draftId
		Step    string      `json:"step"`
		Data    interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&draftRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate step
	validSteps := map[string]bool{
		"generalDetails": true,
		"mediaPromotion": true,
		"specialGuests":  true,
		"volunteers":     true,
		"donations":      true,
	}

	if !validSteps[draftRequest.Step] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid step name. Must be one of: generalDetails, mediaPromotion, specialGuests, volunteers, donations"})
		return
	}

	// Convert draftId to uint pointer
	var draftID *uint
	if draftRequest.DraftID != nil {
		var id uint
		switch v := draftRequest.DraftID.(type) {
		case float64:
			id = uint(v)
			draftID = &id
		case string:
			parsed, err := strconv.ParseUint(v, 10, 64)
			if err == nil {
				id = uint(parsed)
				draftID = &id
			}
		case uint:
			draftID = &v
		case int:
			id := uint(v)
			draftID = &id
		}
	}

	// Get user email from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Get user email from database
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user information"})
		return
	}

	// Save draft (returns draftId, not eventId)
	// Convert draftRequest.Data to map[string]interface{}
	var dataMap map[string]interface{}
	if draftRequest.Data != nil {
		if data, ok := draftRequest.Data.(map[string]interface{}); ok {
			dataMap = data
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
			return
		}
	} else {
		dataMap = make(map[string]interface{})
	}

	savedDraftID, err := services.SaveDraft(draftID, draftRequest.Step, dataMap, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"draftId": savedDraftID, // Changed from eventId to draftId
		"message": "Draft saved successfully",
	})
}

// ----------------------------------------------------
// Get Draft
// ----------------------------------------------------

// GetDraftHandler godoc
// @Summary Get draft data by draft ID
// @Description Retrieves draft data for event creation. Returns all draft steps (generalDetails, mediaPromotion, specialGuests, volunteers, donations) stored in the event_drafts table.
// @Tags Events
// @Security ApiKeyAuth
// @Produce json
// @Param draftId path int true "Draft ID"
// @Success 200 {object} map[string]interface{} "Draft data" example({"draftId":1,"generalDetails":{},"mediaPromotion":{},"specialGuests":{},"volunteers":{},"donations":{}})
// @Failure 400 {object} map[string]string "Bad Request" example({"error":"Invalid draft ID"})
// @Failure 404 {object} map[string]string "Not Found" example({"error":"Draft not found"})
// @Failure 500 {object} map[string]string "Internal Server Error" example({"error":"Failed to retrieve draft"})
// @Router /api/events/draft/{draftId} [get]
func GetDraftHandler(c *gin.Context) {
	draftIDParam := c.Param("draftId")
	draftID, err := strconv.ParseUint(draftIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
		return
	}

	draft, err := services.GetDraft(uint(draftID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"draftId":        draft.ID,
		"generalDetails": draft.GeneralDetailsDraft,
		"mediaPromotion": draft.MediaPromotionDraft,
		"specialGuests":  draft.SpecialGuestsDraft,
		"volunteers":     draft.VolunteersDraft,
		"donations":      draft.DonationsDraft,
		"createdOn":      draft.CreatedOn,
		"updatedOn":      draft.UpdatedOn,
	})
}

// ----------------------------------------------------
// Get Latest Draft by User
// ----------------------------------------------------

// GetLatestDraftByUserHandler godoc
// @Summary Get latest draft for current user
// @Description Retrieves the most recent draft for the authenticated user. Used to restore draft after logout/login.
// @Tags Events
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Draft data" example({"draftId":1,"generalDetails":{},"mediaPromotion":{},"specialGuests":{},"volunteers":{},"donations":{}})
// @Failure 404 {object} map[string]string "Not Found" example({"error":"No draft found for user"})
// @Failure 500 {object} map[string]string "Internal Server Error" example({"error":"Failed to retrieve draft"})
// @Router /api/events/draft/latest [get]
func GetLatestDraftByUserHandler(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Get user email from database
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user information"})
		return
	}

	// Get latest draft for this user
	draft, err := services.GetLatestDraftByUserEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"draftId":        draft.ID,
		"generalDetails": draft.GeneralDetailsDraft,
		"mediaPromotion": draft.MediaPromotionDraft,
		"specialGuests":  draft.SpecialGuestsDraft,
		"volunteers":     draft.VolunteersDraft,
		"donations":      draft.DonationsDraft,
		"createdOn":      draft.CreatedOn,
		"updatedOn":      draft.UpdatedOn,
	})
}

// ----------------------------------------------------
// Update Event Status
// ----------------------------------------------------

// UpdateEventStatusHandler godoc
// @Summary Update event status
// @Description Update the status of an event (complete or incomplete)
// @Tags Events
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param event_id path int true "Event ID"
// @Param status body object true "Status update" example({"status":"complete"})
// @Success 200 {object} map[string]interface{} "Status updated successfully" example({"message":"Event status updated successfully","status":"complete"})
// @Failure 400 {object} map[string]string "Bad Request" example({"error":"Invalid status. Must be 'complete' or 'incomplete'"})
// @Failure 404 {object} map[string]string "Not Found" example({"error":"Event not found"})
// @Failure 500 {object} map[string]string "Internal Server Error" example({"error":"Failed to update event status"})
// @Router /api/events/{event_id}/status [patch]
func UpdateEventStatusHandler(c *gin.Context) {
	eventIDParam := c.Param("event_id")
	eventID, err := strconv.ParseUint(eventIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var request struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateEventStatus(uint(eventID), request.Status); err != nil {
		if err.Error() == "event not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event status updated successfully",
		"status":  request.Status,
	})
}

// Helper function to parse event from map (handles string dates)
func parseEventFromMap(data map[string]interface{}, event *models.EventDetails) error {
	// Parse basic fields
	if id, ok := data["id"].(float64); ok {
		event.ID = uint(id)
	}
	if eventTypeID, ok := data["event_type_id"].(float64); ok {
		event.EventTypeID = uint(eventTypeID)
	}
	if eventCategoryID, ok := data["event_category_id"].(float64); ok {
		event.EventCategoryID = uint(eventCategoryID)
	}
	if scale, ok := data["scale"].(string); ok {
		event.Scale = scale
	}
	if theme, ok := data["theme"].(string); ok {
		event.Theme = theme
	}
	if spiritualOrator, ok := data["spiritual_orator"].(string); ok {
		event.SpiritualOrator = spiritualOrator
	}
	if country, ok := data["country"].(string); ok {
		event.Country = country
	}
	if state, ok := data["state"].(string); ok {
		event.State = state
	}
	if city, ok := data["city"].(string); ok {
		event.City = city
	}
	if district, ok := data["district"].(string); ok {
		event.District = district
	}
	if postOffice, ok := data["post_office"].(string); ok {
		event.PostOffice = postOffice
	}
	if pincode, ok := data["pincode"].(string); ok {
		event.Pincode = pincode
	}
	if address, ok := data["address"].(string); ok {
		event.Address = address
	}
	if status, ok := data["status"].(string); ok {
		event.Status = status
	}

	// Parse dates from strings
	if startDateStr, ok := data["start_date"].(string); ok && startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			event.StartDate = t
		} else if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			event.StartDate = t
		}
	}
	if endDateStr, ok := data["end_date"].(string); ok && endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			event.EndDate = t
		} else if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			event.EndDate = t
		}
	}

	// Parse times from strings (HH:mm or HH:mm:ss format)
	if dailyStartTimeStr, ok := data["daily_start_time"].(string); ok && dailyStartTimeStr != "" {
		parts := strings.Split(dailyStartTimeStr, ":")
		if len(parts) >= 2 {
			if hour, err1 := strconv.Atoi(parts[0]); err1 == nil {
				if minute, err2 := strconv.Atoi(parts[1]); err2 == nil {
					// Create time with today's date
					now := time.Now()
					t := models.TimeOnly{
						Time: time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, time.UTC),
					}
					event.DailyStartTime = &t
				}
			}
		}
	}
	if dailyEndTimeStr, ok := data["daily_end_time"].(string); ok && dailyEndTimeStr != "" {
		parts := strings.Split(dailyEndTimeStr, ":")
		if len(parts) >= 2 {
			if hour, err1 := strconv.Atoi(parts[0]); err1 == nil {
				if minute, err2 := strconv.Atoi(parts[1]); err2 == nil {
					// Create time with today's date
					now := time.Now()
					t := models.TimeOnly{
						Time: time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, time.UTC),
					}
					event.DailyEndTime = &t
				}
			}
		}
	}

	// Parse numeric fields
	if benMen, ok := data["beneficiary_men"].(float64); ok {
		event.BeneficiaryMen = int(benMen)
	}
	if benWomen, ok := data["beneficiary_women"].(float64); ok {
		event.BeneficiaryWomen = int(benWomen)
	}
	if benChild, ok := data["beneficiary_child"].(float64); ok {
		event.BeneficiaryChild = int(benChild)
	}
	if initMen, ok := data["initiation_men"].(float64); ok {
		event.InitiationMen = int(initMen)
	}
	if initWomen, ok := data["initiation_women"].(float64); ok {
		event.InitiationWomen = int(initWomen)
	}
	if initChild, ok := data["initiation_child"].(float64); ok {
		event.InitiationChild = int(initChild)
	}

	return nil
}

// ----------------------------------------------------
// Export Events to Excel
// ----------------------------------------------------

// ExportEventsHandler godoc
// @Summary Export events to Excel
// @Description Export events to Excel file filtered by date range (optional)
// @Tags Events
// @Security ApiKeyAuth
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param status query string false "Status filter (complete/incomplete)"
// @Success 200 {file} file "Excel file"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/events/export [get]
func ExportEventsHandler(c *gin.Context) {
	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	statusFilter := c.Query("status")
	
	// Always filter by created_on date
	dateFilterType := "created_on"

	var startDate *time.Time
	var endDate *time.Time

	// Parse start date
	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
		// Set start date to start of day (00:00:00) in UTC to ensure consistent comparison
		startOfDay := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
		startDate = &startOfDay
	}

	// Parse end date
	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
		// Set end date to end of day
		endOfDay := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 999999999, time.UTC)
		endDate = &endOfDay
	}

	// Validate date range
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before or equal to end_date"})
		return
	}

	// Get events by date range
	events, err := services.GetEventsByDateRange(startDate, endDate, statusFilter, dateFilterType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events: " + err.Error()})
		return
	}

	// Export to Excel
	excelBuffer, err := services.ExportEventsToExcel(events)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel file: " + err.Error()})
		return
	}

	// Generate filename with date range
	filename := "events_export"
	if startDate != nil && endDate != nil {
		filename = fmt.Sprintf("events_%s_to_%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	} else if startDate != nil {
		filename = fmt.Sprintf("events_from_%s", startDate.Format("2006-01-02"))
	} else if endDate != nil {
		filename = fmt.Sprintf("events_until_%s", endDate.Format("2006-01-02"))
	}
	filename += ".xlsx"

	// Set response headers
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", excelBuffer.Len()))

	// Write Excel buffer to response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBuffer.Bytes())
}
