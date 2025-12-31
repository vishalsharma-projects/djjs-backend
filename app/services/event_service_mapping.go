package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/config"
)

// MapFrontendPayloadToEvent maps frontend payload structure to EventDetails model
func MapFrontendPayloadToEvent(generalDetails map[string]interface{}, involvedParticipants map[string]interface{}) (*models.EventDetails, error) {
	event := &models.EventDetails{}

	// Map general details - handle both "type" and "eventType" from frontend
	var eventTypeName string
	if val, ok := generalDetails["eventType"].(string); ok && val != "" {
		eventTypeName = val
	} else if val, ok := generalDetails["type"].(string); ok && val != "" {
		eventTypeName = val
	}

	if eventTypeName != "" {
		// Find event type ID by name
		var eventType models.EventType
		if err := config.DB.Where("name = ?", eventTypeName).First(&eventType).Error; err == nil {
			event.EventTypeID = eventType.ID
		} else {
			return nil, fmt.Errorf("event type '%s' not found in database", eventTypeName)
		}
	} else {
		return nil, fmt.Errorf("eventType is required in generalDetails")
	}

	// Handle eventCategory - check both "eventCategory" and "eventName" (which might be used as category)
	var eventCategoryName string
	if val, ok := generalDetails["eventCategory"].(string); ok && val != "" {
		eventCategoryName = val
	} else if val, ok := generalDetails["eventName"].(string); ok && val != "" {
		// Sometimes eventName is used as category
		eventCategoryName = val
	}

	if eventCategoryName != "" {
		// Find event category ID by name
		var eventCategory models.EventCategory
		if err := config.DB.Where("name = ?", eventCategoryName).First(&eventCategory).Error; err == nil {
			event.EventCategoryID = eventCategory.ID
		} else {
			return nil, fmt.Errorf("event category '%s' not found in database", eventCategoryName)
		}
	} else {
		return nil, fmt.Errorf("eventCategory is required in generalDetails")
	}

	// Handle eventSubCategory - optional field
	var eventSubCategoryName string
	if val, ok := generalDetails["eventSubCategory"].(string); ok && val != "" {
		eventSubCategoryName = val
	} else if val, ok := generalDetails["event_sub_category"].(string); ok && val != "" {
		eventSubCategoryName = val
	}

	if eventSubCategoryName != "" {
		// Find event sub category ID by name (must belong to the selected category)
		var eventSubCategory models.EventSubCategory
		if err := config.DB.Where("name = ? AND event_category_id = ?", eventSubCategoryName, event.EventCategoryID).First(&eventSubCategory).Error; err == nil {
			event.EventSubCategoryID = &eventSubCategory.ID
		} else {
			// If sub category not found, log warning but don't fail (it's optional)
			fmt.Printf("Warning: event sub category '%s' not found for category ID %d\n", eventSubCategoryName, event.EventCategoryID)
		}
	}

	if scale, ok := generalDetails["scale"].(string); ok {
		event.Scale = scale
	}

	if theme, ok := generalDetails["theme"].(string); ok {
		event.Theme = theme
	}

	// Parse duration string to start_date and end_date
	if duration, ok := generalDetails["duration"].(string); ok && duration != "" {
		// Check if duration contains the separator
		if !strings.Contains(duration, " - ") {
			return nil, fmt.Errorf("duration must be in format 'dd MMM yyyy - dd MMM yyyy' (e.g., '01 Jan 2024 - 05 Jan 2024'), got: %s", duration)
		}

		dates := strings.Split(duration, " - ")
		if len(dates) == 2 {
			startDateStr := strings.TrimSpace(dates[0])
			endDateStr := strings.TrimSpace(dates[1])

			if startDateStr == "" || endDateStr == "" {
				return nil, fmt.Errorf("duration must include both start and end dates")
			}

			startDate, err1 := parseDate(startDateStr)
			endDate, err2 := parseDate(endDateStr)
			if err1 != nil {
				return nil, fmt.Errorf("invalid start date format in duration: %v. Expected format: 'dd MMM yyyy' (e.g., '01 Jan 2024')", err1)
			}
			if err2 != nil {
				return nil, fmt.Errorf("invalid end date format in duration: %v. Expected format: 'dd MMM yyyy' (e.g., '05 Jan 2024')", err2)
			}

			event.StartDate = startDate
			event.EndDate = endDate
		} else {
			return nil, fmt.Errorf("duration must be in format 'dd MMM yyyy - dd MMM yyyy' (e.g., '01 Jan 2024 - 05 Jan 2024')")
		}
	} else {
		// Duration is empty or missing - only require it for complete events
		// Check if start_date and end_date are provided directly
		if startDateStr, ok := generalDetails["start_date"].(string); ok && startDateStr != "" {
			if startDate, err := parseDate(startDateStr); err == nil {
				event.StartDate = startDate
			}
		}
		if endDateStr, ok := generalDetails["end_date"].(string); ok && endDateStr != "" {
			if endDate, err := parseDate(endDateStr); err == nil {
				event.EndDate = endDate
			}
		}
		// If still no dates, return error
		if event.StartDate.IsZero() || event.EndDate.IsZero() {
			return nil, fmt.Errorf("duration is required in generalDetails. Please provide dates in format 'dd MMM yyyy - dd MMM yyyy' (e.g., '01 Jan 2024 - 05 Jan 2024')")
		}
	}

	// Parse daily times
	if dailyStartTime, ok := generalDetails["dailyStartTime"].(string); ok && dailyStartTime != "" {
		if t, err := parseTime(dailyStartTime); err == nil {
			timeOnly := models.TimeOnly{Time: t}
			event.DailyStartTime = &timeOnly
		}
	}

	if dailyEndTime, ok := generalDetails["dailyEndTime"].(string); ok && dailyEndTime != "" {
		if t, err := parseTime(dailyEndTime); err == nil {
			timeOnly := models.TimeOnly{Time: t}
			event.DailyEndTime = &timeOnly
		}
	}

	if spiritualOrator, ok := generalDetails["spiritualOrator"].(string); ok {
		event.SpiritualOrator = spiritualOrator
	}

	if language, ok := generalDetails["language"].(string); ok && language != "" {
		event.Language = language
	}

	// Location fields - handle both direct fields and nested venue structure
	var venue map[string]interface{}
	if v, ok := generalDetails["venue"].(map[string]interface{}); ok {
		venue = v
	}

	// Country - check venue first, then direct field
	if venue != nil {
		if country, ok := venue["country"].(string); ok && country != "" {
			event.Country = country
		}
	}
	if event.Country == "" {
		if country, ok := generalDetails["country"].(string); ok && country != "" {
			event.Country = country
		}
	}

	// State
	if venue != nil {
		if state, ok := venue["state"].(string); ok && state != "" {
			event.State = state
		}
	}
	if event.State == "" {
		if state, ok := generalDetails["state"].(string); ok && state != "" {
			event.State = state
		}
	}

	// District
	if venue != nil {
		if district, ok := venue["district"].(string); ok && district != "" {
			event.District = district
		}
	}
	if event.District == "" {
		if district, ok := generalDetails["district"].(string); ok && district != "" {
			event.District = district
		}
	}

	// City
	if venue != nil {
		if city, ok := venue["city"].(string); ok && city != "" {
			event.City = city
		}
	}
	if event.City == "" {
		if city, ok := generalDetails["city"].(string); ok && city != "" {
			event.City = city
		}
	}

	// Pincode
	if venue != nil {
		if pincode, ok := venue["pincode"].(string); ok && pincode != "" {
			event.Pincode = pincode
		}
	}
	if event.Pincode == "" {
		if pincode, ok := generalDetails["pincode"].(string); ok && pincode != "" {
			event.Pincode = pincode
		}
	}

	// Post Office
	if venue != nil {
		if postOffice, ok := venue["postOffice"].(string); ok && postOffice != "" {
			event.PostOffice = postOffice
		}
	}
	if event.PostOffice == "" {
		if postOffice, ok := generalDetails["postOffice"].(string); ok && postOffice != "" {
			event.PostOffice = postOffice
		}
	}

	// Address
	if venue != nil {
		if address, ok := venue["address"].(string); ok && address != "" {
			event.Address = address
		}
	}
	if event.Address == "" {
		if address, ok := generalDetails["address"].(string); ok && address != "" {
			event.Address = address
		}
	}

	// Address Type
	if addressType, ok := generalDetails["addressType"].(string); ok && addressType != "" {
		event.AddressType = addressType
	} else if addressType, ok := generalDetails["address_type"].(string); ok && addressType != "" {
		event.AddressType = addressType
	}

	// Police Station
	if policeStation, ok := generalDetails["policeStation"].(string); ok && policeStation != "" {
		event.PoliceStation = policeStation
	} else if policeStation, ok := generalDetails["police_station"].(string); ok && policeStation != "" {
		event.PoliceStation = policeStation
	}

	// Area Covered
	if areaCovered, ok := generalDetails["areaCovered"].(string); ok && areaCovered != "" {
		event.AreaCovered = areaCovered
	} else if areaCovered, ok := generalDetails["area_covered"].(string); ok && areaCovered != "" {
		event.AreaCovered = areaCovered
	}

	// Map branch_id (optional field)
	if branchId, ok := generalDetails["branchId"].(float64); ok && branchId > 0 {
		branchIDUint := uint(branchId)
		event.BranchID = &branchIDUint
	} else if branchId, ok := generalDetails["branch_id"].(float64); ok && branchId > 0 {
		// Also handle snake_case for backward compatibility
		branchIDUint := uint(branchId)
		event.BranchID = &branchIDUint
	}

	// Map involved participants
	if involvedParticipants != nil {
		if val, ok := involvedParticipants["beneficiariesMen"].(float64); ok {
			event.BeneficiaryMen = int(val)
		}
		if val, ok := involvedParticipants["beneficiariesWomen"].(float64); ok {
			event.BeneficiaryWomen = int(val)
		}
		if val, ok := involvedParticipants["beneficiariesChildren"].(float64); ok {
			event.BeneficiaryChild = int(val)
		}
		if val, ok := involvedParticipants["initiationMen"].(float64); ok {
			event.InitiationMen = int(val)
		}
		if val, ok := involvedParticipants["initiationWomen"].(float64); ok {
			event.InitiationWomen = int(val)
		}
		if val, ok := involvedParticipants["initiationChildren"].(float64); ok {
			event.InitiationChild = int(val)
		}
	}

	return event, nil
}

// MapFrontendPayloadToEventWithStatus maps frontend payload with status support
func MapFrontendPayloadToEventWithStatus(generalDetails map[string]interface{}, involvedParticipants map[string]interface{}, status string) (*models.EventDetails, error) {
	event, err := MapFrontendPayloadToEvent(generalDetails, involvedParticipants)
	if err != nil {
		return nil, err
	}

	// Set status if provided, otherwise default to incomplete
	if status != "" && (status == "complete" || status == "incomplete") {
		event.Status = status
	} else {
		event.Status = "incomplete"
	}

	return event, nil
}

// Helper functions for parsing dates and times
func parseDate(dateStr string) (time.Time, error) {
	// Try different date formats
	formats := []string{
		"02 Jan 2006",
		"02-01-2006",
		"2006-01-02",
		"01/02/2006",
	}

	dateStr = strings.TrimSpace(dateStr)
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func parseTime(timeStr string) (time.Time, error) {
	// Parse time in HH:mm format
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hour, err1 := strconv.Atoi(parts[0])
	minute, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return time.Time{}, fmt.Errorf("invalid time values: %s", timeStr)
	}

	// Create a time with today's date
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location()), nil
}

