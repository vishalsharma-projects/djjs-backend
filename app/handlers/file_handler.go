package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
	"github.com/followCode/djjs-event-reporting-backend/app/services"
	"github.com/followCode/djjs-event-reporting-backend/config"
	"github.com/gin-gonic/gin"
)

// UploadFileHandler handles file uploads to S3
// @Summary Upload file to S3
// @Description Upload image, video, audio, or PDF file to S3 and associate with event media
// @Tags Files
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload (image, video, audio, or PDF)"
// @Param event_id formData int true "Event ID"
// @Param media_id formData int false "Media ID (if updating existing media)"
// @Param category formData string false "File category (Event Photos, Video Coverage, Testimonials, Press Release)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/files/upload [post]
func UploadFileHandler(c *gin.Context) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Get event ID
	eventIDStr := c.PostForm("event_id")
	if eventIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event_id is required"})
		return
	}
	eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	// Get media ID if provided (for updating existing media)
	var mediaID uint
	mediaIDStr := c.PostForm("media_id")
	if mediaIDStr != "" {
		id, err := strconv.ParseUint(mediaIDStr, 10, 64)
		if err == nil {
			mediaID = uint(id)
		}
	}

	// Get category
	category := c.PostForm("category")
	if category == "" {
		category = "Event Photos"
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	// Read file data (handle large files efficiently)
	fileData := make([]byte, file.Size)
	n, err := src.Read(fileData)
	if err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to read file: %v", err),
		})
		return
	}
	if int64(n) != file.Size {
		// Adjust slice if read less than expected
		fileData = fileData[:n]
	}

	// Get content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		// Try to determine from extension
		ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".webp":
			contentType = "image/webp"
		case ".bmp":
			contentType = "image/bmp"
		case ".svg":
			contentType = "image/svg+xml"
		case ".mp4":
			contentType = "video/mp4"
		case ".mov":
			contentType = "video/quicktime"
		case ".avi":
			contentType = "video/x-msvideo"
		case ".wmv":
			contentType = "video/x-ms-wmv"
		case ".webm":
			contentType = "video/webm"
		case ".mkv":
			contentType = "video/x-matroska"
		case ".mp3":
			contentType = "audio/mpeg"
		case ".wav":
			contentType = "audio/wav"
		case ".ogg":
			contentType = "audio/ogg"
		case ".aac":
			contentType = "audio/aac"
		case ".m4a":
			contentType = "audio/x-m4a"
		case ".flac":
			contentType = "audio/flac"
		case ".pdf":
			contentType = "application/pdf"
		case ".doc":
			contentType = "application/msword"
		case ".docx":
			contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		case ".xls":
			contentType = "application/vnd.ms-excel"
		case ".xlsx":
			contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		case ".ppt":
			contentType = "application/vnd.ms-powerpoint"
		case ".pptx":
			contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
		default:
			contentType = "application/octet-stream"
		}
	}

	// Determine file type category first (needed for size validation)
	fileType := services.GetFileTypeFromContentType(contentType)

	// Validate file size
	if err := services.ValidateFileSize(file.Size, fileType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate file type
	if !services.ValidateFileType(contentType) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file type not allowed. Allowed types: " +
				"Images (JPEG, PNG, GIF, WebP, BMP, SVG), " +
				"Videos (MP4, MOV, AVI, WMV, WebM, MKV), " +
				"Audio (MP3, WAV, OGG, AAC, M4A, FLAC), " +
				"Documents (PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX)",
		})
		return
	}

	folder := services.GetFolderFromFileType(fileType)

	// Upload to S3
	fileURL, err := services.UploadFile(c.Request.Context(), fileData, file.Filename, contentType, folder)
	if err != nil {
		// Log detailed error for debugging
		fmt.Printf("S3 Upload Error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("failed to upload file to S3: %v", err),
		})
		return
	}

	// Update or create EventMedia record
	if mediaID > 0 {
		// Update existing media
		var media models.EventMedia
		if err := config.DB.First(&media, mediaID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
			return
		}

		media.FileURL = fileURL
		media.FileType = fileType
		if err := config.DB.Save(&media).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update media record"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded and media updated successfully",
			"data": gin.H{
				"media_id":  media.ID,
				"file_url":  fileURL,
				"file_type": fileType,
			},
		})
	} else {
		// Create new media record (minimal record, can be updated later)
		media := models.EventMedia{
			EventID:     uint(eventID),
			FileURL:     fileURL,
			FileType:    fileType,
			CompanyName: file.Filename,
			FirstName:   "Uploaded",
			LastName:    "File",
		}

		// Try to get a default media coverage type
		var mediaType models.MediaCoverageType
		if err := config.DB.First(&mediaType).Error; err == nil {
			media.MediaCoverageTypeID = mediaType.ID
		}

		if err := config.DB.Create(&media).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create media record"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "File uploaded successfully",
			"data": gin.H{
				"media_id":  media.ID,
				"file_url":  fileURL,
				"file_type": fileType,
			},
		})
	}
}

// DownloadFileHandler generates a presigned URL for downloading a file
// @Summary Get download URL for file
// @Description Generates a presigned URL for downloading a file from S3
// @Tags Files
// @Security ApiKeyAuth
// @Produce json
// @Param media_id path int true "Media ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/files/{media_id}/download [get]
func DownloadFileHandler(c *gin.Context) {
	mediaIDStr := c.Param("media_id")
	mediaID, err := strconv.ParseUint(mediaIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid media_id"})
		return
	}

	// Get media record
	var media models.EventMedia
	if err := config.DB.First(&media, mediaID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	if media.FileURL == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "file URL not found for this media"})
		return
	}

	// Extract S3 key from URL
	s3Key := services.GetS3KeyFromURL(media.FileURL)
	if s3Key == "" {
		// If URL is already a presigned URL or direct URL, return it
		c.JSON(http.StatusOK, gin.H{
			"download_url": media.FileURL,
			"file_type":    media.FileType,
			"file_name":    media.CompanyName,
		})
		return
	}

	// Try to get original filename from S3 metadata
	originalFilename := services.GetOriginalFilename(c.Request.Context(), s3Key)
	if originalFilename == "" {
		// Fallback to company name or generate from URL
		originalFilename = media.CompanyName
		if originalFilename == "" {
			// Extract filename from S3 key
			parts := strings.Split(s3Key, "/")
			if len(parts) > 0 {
				originalFilename = parts[len(parts)-1]
			}
		}
	}

	// Generate presigned URL (valid for 1 hour)
	// Note: Glacier Instant Retrieval provides instant access, so no special handling needed
	presignedURL, err := services.GetPresignedURL(c.Request.Context(), s3Key, time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate download URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"download_url": presignedURL,
		"file_type":    media.FileType,
		"file_name":    originalFilename,
	})
}

// DeleteFileHandler deletes a file from S3 and the media record
// @Summary Delete file from S3
// @Description Deletes a file from S3 and optionally the media record. Optionally validates event_id to ensure file belongs to specific event.
// @Tags Files
// @Security ApiKeyAuth
// @Produce json
// @Param media_id path int true "Media ID"
// @Param event_id query int false "Event ID (optional, for validation)"
// @Param delete_record query bool false "Delete media record from database (default: true)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/files/{media_id} [delete]
func DeleteFileHandler(c *gin.Context) {
	mediaIDStr := c.Param("media_id")
	mediaID, err := strconv.ParseUint(mediaIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid media_id"})
		return
	}

	// Get media record
	var media models.EventMedia
	if err := config.DB.First(&media, mediaID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	// Validate event_id if provided (ensures file belongs to the event)
	eventIDStr := c.Query("event_id")
	if eventIDStr != "" {
		eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
		if err == nil {
			if media.EventID != uint(eventID) {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "file does not belong to the specified event",
				})
				return
			}
		}
	}

	// Delete from S3 if file URL exists
	if media.FileURL != "" {
		s3Key := services.GetS3KeyFromURL(media.FileURL)
		if s3Key != "" {
			if err := services.DeleteFile(c.Request.Context(), s3Key); err != nil {
				// Log error but continue to delete database record
				fmt.Printf("Warning: failed to delete file from S3: %v\n", err)
			}
		}
	}

	// Delete media record if requested (default: true)
	deleteRecord := c.DefaultQuery("delete_record", "true")
	if deleteRecord == "true" {
		if err := config.DB.Delete(&media).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete media record"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "File and media record deleted successfully"})
	} else {
		// Just clear the file URL
		media.FileURL = ""
		media.FileType = ""
		if err := config.DB.Save(&media).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update media record"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully, media record kept"})
	}
}

// UploadMultipleFilesHandler handles multiple file uploads to S3 in a single request
// @Summary Upload multiple files to S3
// @Description Upload multiple image, video, audio, or PDF files to S3 and associate with event media
// @Tags Files
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Files to upload (multiple files allowed)"
// @Param event_id formData int true "Event ID"
// @Param category formData string false "File category (Event Photos, Video Coverage, Testimonials, Press Release)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/files/upload-multiple [post]
func UploadMultipleFilesHandler(c *gin.Context) {
	// Get event ID
	eventIDStr := c.PostForm("event_id")
	if eventIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event_id is required"})
		return
	}
	eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	// Get category
	category := c.PostForm("category")
	if category == "" {
		category = "Event Photos"
	}

	// Get multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
		return
	}

	// Get all files
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files provided"})
		return
	}

	// Process each file
	var results []map[string]interface{}
	var errors []string

	for _, fileHeader := range files {
		// Open file
		src, err := fileHeader.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to open file", fileHeader.Filename))
			continue
		}

		// Read file data
		fileData := make([]byte, fileHeader.Size)
		n, err := src.Read(fileData)
		if err != nil && err.Error() != "EOF" {
			src.Close()
			errors = append(errors, fmt.Sprintf("%s: failed to read file", fileHeader.Filename))
			continue
		}
		if int64(n) != fileHeader.Size {
			fileData = fileData[:n]
		}
		src.Close()

		// Get content type
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			// Try to determine from extension
			ext := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):])
			switch ext {
			case ".jpg", ".jpeg":
				contentType = "image/jpeg"
			case ".png":
				contentType = "image/png"
			case ".gif":
				contentType = "image/gif"
			case ".webp":
				contentType = "image/webp"
			case ".bmp":
				contentType = "image/bmp"
			case ".svg":
				contentType = "image/svg+xml"
			case ".mp4":
				contentType = "video/mp4"
			case ".mov":
				contentType = "video/quicktime"
			case ".avi":
				contentType = "video/x-msvideo"
			case ".wmv":
				contentType = "video/x-ms-wmv"
			case ".webm":
				contentType = "video/webm"
			case ".mkv":
				contentType = "video/x-matroska"
			case ".mp3":
				contentType = "audio/mpeg"
			case ".wav":
				contentType = "audio/wav"
			case ".ogg":
				contentType = "audio/ogg"
			case ".aac":
				contentType = "audio/aac"
			case ".m4a":
				contentType = "audio/x-m4a"
			case ".flac":
				contentType = "audio/flac"
			case ".pdf":
				contentType = "application/pdf"
			case ".doc":
				contentType = "application/msword"
			case ".docx":
				contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			case ".xls":
				contentType = "application/vnd.ms-excel"
			case ".xlsx":
				contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
			case ".ppt":
				contentType = "application/vnd.ms-powerpoint"
			case ".pptx":
				contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
			default:
				contentType = "application/octet-stream"
			}
		}

		// Determine file type category
		fileType := services.GetFileTypeFromContentType(contentType)

		// Validate file size
		if err := services.ValidateFileSize(fileHeader.Size, fileType); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", fileHeader.Filename, err))
			continue
		}

		// Validate file type
		if !services.ValidateFileType(contentType) {
			errors = append(errors, fmt.Sprintf("%s: file type not allowed", fileHeader.Filename))
			continue
		}

		folder := services.GetFolderFromFileType(fileType)

		// Upload to S3
		fileURL, err := services.UploadFile(c.Request.Context(), fileData, fileHeader.Filename, contentType, folder)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", fileHeader.Filename, err))
			continue
		}

		// Create EventMedia record
		media := models.EventMedia{
			EventID:     uint(eventID),
			FileURL:     fileURL,
			FileType:    fileType,
			CompanyName: fileHeader.Filename,
			FirstName:   "Uploaded",
			LastName:    "File",
		}

		// Try to get a default media coverage type
		var mediaType models.MediaCoverageType
		if err := config.DB.First(&mediaType).Error; err == nil {
			media.MediaCoverageTypeID = mediaType.ID
		}

		if err := config.DB.Create(&media).Error; err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to create media record", fileHeader.Filename))
			continue
		}

		results = append(results, map[string]interface{}{
			"filename":  fileHeader.Filename,
			"media_id":  media.ID,
			"file_url":  fileURL,
			"file_type": fileType,
			"status":    "success",
		})
	}

	// Return results
	response := map[string]interface{}{
		"message": fmt.Sprintf("Processed %d file(s)", len(files)),
		"success": len(results),
		"failed":  len(errors),
		"results": results,
	}

	if len(errors) > 0 {
		response["errors"] = errors
	}

	if len(results) > 0 {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, response)
	}
}
