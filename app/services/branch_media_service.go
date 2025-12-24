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

// CreateBranchMedia creates a new BranchMedia record
func CreateBranchMedia(media *models.BranchMedia) error {
	return config.DB.Create(media).Error
}

// GetAllBranchMedia retrieves all BranchMedia records
func GetAllBranchMedia() ([]models.BranchMedia, error) {
	var medias []models.BranchMedia
	if err := config.DB.
		Preload("Branch").
		Find(&medias).Error; err != nil {
		return nil, err
	}
	return medias, nil
}

// GetBranchMediaByBranchID retrieves all BranchMedia records by BranchID
func GetBranchMediaByBranchID(branchID uint, isChildBranch bool) ([]models.BranchMedia, error) {
	var mediaList []models.BranchMedia
	if err := config.DB.
		Preload("Branch").
		Where("branch_id = ? AND is_child_branch = ?", branchID, isChildBranch).
		Find(&mediaList).Error; err != nil {
		return nil, errors.New("no branch media found for the given branch ID")
	}
	return mediaList, nil
}

// UpdateBranchMedia updates an existing BranchMedia record
func UpdateBranchMedia(media *models.BranchMedia) error {
	return config.DB.Save(media).Error
}

// DeleteBranchMedia deletes a BranchMedia record
func DeleteBranchMedia(mediaID uint) error {
	return config.DB.Delete(&models.BranchMedia{}, mediaID).Error
}

// GetBranchMediaByID retrieves a BranchMedia record by ID
func GetBranchMediaByID(mediaID uint) (*models.BranchMedia, error) {
	var media models.BranchMedia
	if err := config.DB.First(&media, mediaID).Error; err != nil {
		return nil, errors.New("branch media not found")
	}
	return &media, nil
}

// ConvertBranchMediaToPresignedURLs converts BranchMedia items to include presigned URLs
// This function takes a slice of BranchMedia and returns a new slice with presigned URLs
// All media access uses short-lived pre-signed URLs for security
// REQUIRES: media.S3Key must be set - will fail if S3Key is empty
func ConvertBranchMediaToPresignedURLs(ctx context.Context, mediaList []models.BranchMedia) ([]models.BranchMedia, error) {
	result := make([]models.BranchMedia, len(mediaList))
	
	for i, media := range mediaList {
		result[i] = media
		
		// ENFORCE S3Key-only access - DO NOT fall back to FileURL
		if media.S3Key == "" {
			// Return error if S3Key is missing - fail fast
			return nil, fmt.Errorf("branch media item ID %d has empty S3Key - cannot generate presigned URL. Run backfill migration to populate s3_key from file_url", media.ID)
		}
		
		// Generate short-lived presigned URL (15 minutes for gallery listing)
		presignedURL, err := GetPresignedURL(ctx, media.S3Key, 15*time.Minute)
		if err != nil {
			// Fail fast on presign errors - return HTTP 500
			return nil, fmt.Errorf("failed to generate presigned URL for branch media ID %d (s3_key: %s): %w", media.ID, media.S3Key, err)
		}
		
		// Defensive check: ensure URL is presigned (contains X-Amz-Signature)
		if !strings.Contains(presignedURL, "X-Amz-Signature") && !strings.Contains(presignedURL, "Signature=") {
			return nil, fmt.Errorf("CRITICAL: generated URL for branch media ID %d does not contain presigned signature: %s", media.ID, presignedURL)
		}
		
		// Store presigned URL in URL field (for JSON serialization)
		// FileURL is internal and not serialized
		result[i].FileURL = presignedURL // Internal storage
		result[i].URL = presignedURL     // JSON response field
	}
	
	return result, nil
}


