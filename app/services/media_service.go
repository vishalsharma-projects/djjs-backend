package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
// Deprecated: Use GetEventMediaByEventIDPaginated for cursor-based pagination
func GetEventMediaByEventID(eventID uint) ([]models.EventMedia, error) {
	var mediaList []models.EventMedia
	if err := config.DB.
		Preload("Event").
		Preload("MediaCoverageType").
		Where("event_id = ?", eventID).
		Order("created_on DESC, id DESC").
		Find(&mediaList).Error; err != nil {
		return nil, errors.New("no event media found for the given event ID")
	}
	return mediaList, nil
}

// PaginationCursor represents a cursor for pagination
type PaginationCursor struct {
	CreatedAt time.Time
	ID        uint
}

// PaginatedEventMediaResult contains paginated results
type PaginatedEventMediaResult struct {
	Data       []models.EventMedia `json:"data"`
	NextCursor *PaginationCursor   `json:"next_cursor,omitempty"`
	HasMore    bool                `json:"has_more"`
}

// GetEventMediaByEventIDPaginated retrieves EventMedia records with cursor-based pagination
// Uses (created_at, id) as the cursor to avoid OFFSET pagination issues
func GetEventMediaByEventIDPaginated(eventID uint, limit int, cursor *PaginationCursor) (*PaginatedEventMediaResult, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	var mediaList []models.EventMedia
	query := config.DB.
		Preload("Event").
		Preload("MediaCoverageType").
		Where("event_id = ?", eventID)

	// Apply cursor if provided
	if cursor != nil {
		// Use (created_at, id) tuple for cursor-based pagination
		// This ensures stable ordering even with duplicate timestamps
		query = query.Where(
			"(created_on, id) < (?, ?)",
			cursor.CreatedAt,
			cursor.ID,
		)
	}

	// Order by created_at DESC, id DESC for consistent pagination
	// Fetch one extra to check if there's more
	err := query.
		Order("created_on DESC, id DESC").
		Limit(limit + 1).
		Find(&mediaList).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch event media: %w", err)
	}

	hasMore := len(mediaList) > limit
	if hasMore {
		mediaList = mediaList[:limit] // Remove the extra item
	}

	var nextCursor *PaginationCursor
	if hasMore && len(mediaList) > 0 {
		lastItem := mediaList[len(mediaList)-1]
		nextCursor = &PaginationCursor{
			CreatedAt: lastItem.CreatedOn,
			ID:        lastItem.ID,
		}
	}

	return &PaginatedEventMediaResult{
		Data:       mediaList,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
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
	if media.S3Key != "" {
		updates["s3_key"] = media.S3Key
	}
	if media.OriginalFilename != "" {
		updates["original_filename"] = media.OriginalFilename
	}
	if media.ThumbnailS3Key != nil {
		updates["thumbnail_s3_key"] = media.ThumbnailS3Key
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

// ConvertEventMediaToPresignedURLs converts EventMedia items to include presigned URLs
// This function takes a slice of EventMedia and returns a new slice with presigned URLs
// All media access uses short-lived pre-signed URLs for security
// REQUIRES: media.S3Key must be set - will fail if S3Key is empty
func ConvertEventMediaToPresignedURLs(ctx context.Context, mediaList []models.EventMedia) ([]models.EventMedia, error) {
	result := make([]models.EventMedia, len(mediaList))
	
	for i, media := range mediaList {
		result[i] = media
		
		// ENFORCE S3Key-only access - DO NOT fall back to FileURL
		if media.S3Key == "" {
			// Return error if S3Key is missing - fail fast
			return nil, fmt.Errorf("media item ID %d has empty S3Key - cannot generate presigned URL. Run backfill migration to populate s3_key from file_url", media.ID)
		}
		
		// Generate short-lived presigned URL (15 minutes for gallery listing)
		presignedURL, err := GetPresignedURL(ctx, media.S3Key, 15*time.Minute)
		if err != nil {
			// Fail fast on presign errors - return HTTP 500
			return nil, fmt.Errorf("failed to generate presigned URL for media ID %d (s3_key: %s): %w", media.ID, media.S3Key, err)
		}
		
		// Defensive check: ensure URL is presigned (contains X-Amz-Signature)
		if !strings.Contains(presignedURL, "X-Amz-Signature") && !strings.Contains(presignedURL, "Signature=") {
			return nil, fmt.Errorf("CRITICAL: generated URL for media ID %d does not contain presigned signature: %s", media.ID, presignedURL)
		}
		
		// Store presigned URL in URL field (for JSON serialization)
		// FileURL is internal and not serialized
		result[i].FileURL = presignedURL // Internal storage
		result[i].URL = presignedURL     // JSON response field
		
		// Generate thumbnail presigned URL if thumbnail exists
		if media.ThumbnailS3Key != nil && *media.ThumbnailS3Key != "" {
			thumbnailURL, err := GetPresignedURL(ctx, *media.ThumbnailS3Key, 15*time.Minute)
			if err != nil {
				// Fail fast on thumbnail presign errors
				return nil, fmt.Errorf("failed to generate presigned URL for thumbnail of media ID %d (thumbnail_s3_key: %s): %w", media.ID, *media.ThumbnailS3Key, err)
			}
			// Store thumbnail URL in FileURL field temporarily (frontend can use this)
			// Note: In future, consider adding separate thumbnail_url field to response DTO
			_ = thumbnailURL // Placeholder for future thumbnail URL handling
		}
	}
	
	return result, nil
}
