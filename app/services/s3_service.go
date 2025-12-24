package services

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

var (
	S3Client     *s3.Client
	S3Uploader   *manager.Uploader
	S3BucketName string
	S3Region     string
)

// UploadResult contains the result of an S3 upload
type UploadResult struct {
	S3Key          string // Opaque S3 object key (UUID-based)
	OriginalFilename string // Original filename from upload
}

// InitializeS3 initializes the S3 client and uploader with credentials
func InitializeS3() error {
	// Get credentials from environment variables
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")

	// Validate required environment variables
	if accessKeyID == "" {
		return fmt.Errorf("AWS_ACCESS_KEY_ID environment variable is required")
	}
	if secretAccessKey == "" {
		return fmt.Errorf("AWS_SECRET_ACCESS_KEY environment variable is required")
	}
	if bucketName == "" {
		return fmt.Errorf("AWS_S3_BUCKET_NAME environment variable is required")
	}
	if region == "" {
		return fmt.Errorf("AWS_REGION environment variable is required")
	}

	// Create AWS config with static credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID,
			secretAccessKey,
			"",
		)),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	S3Client = s3.NewFromConfig(cfg)
	S3Uploader = manager.NewUploader(S3Client)
	S3BucketName = bucketName
	S3Region = region

	return nil
}

// UploadFile uploads a file to S3 and returns the S3 key and original filename
// S3 keys are opaque UUID-based to decouple from original filenames
func UploadFile(ctx context.Context, fileData []byte, fileName string, contentType string, folder string) (*UploadResult, error) {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return nil, fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	// Generate opaque, collision-safe S3 key using UUID
	// Format: {folder}/{uuid}.{ext}
	ext := filepath.Ext(fileName)
	s3Key := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Upload file to S3 with Standard storage class for immediate access
	storageClass := types.StorageClassStandard
	putInput := &s3.PutObjectInput{
		Bucket:       aws.String(S3BucketName),
		Key:          aws.String(s3Key),
		Body:         bytes.NewReader(fileData),
		ContentType:  aws.String(contentType),
		StorageClass: storageClass,
		Metadata: map[string]string{
			"original-filename": fileName,
			"upload-date":       time.Now().Format(time.RFC3339),
		},
	}

	// Note: ACL is not set because the bucket has ACLs disabled
	// Public access should be configured via bucket policy instead
	// All access should use presigned URLs for security

	_, err := S3Uploader.Upload(ctx, putInput)
	if err != nil {
		// Return detailed error for debugging
		return nil, fmt.Errorf("S3 upload failed (bucket: %s, key: %s): %w", S3BucketName, s3Key, err)
	}

	return &UploadResult{
		S3Key:           s3Key,
		OriginalFilename: fileName,
	}, nil
}

// UploadFileLegacy uploads a file to S3 and returns the S3 URL (legacy compatibility)
// Deprecated: Use UploadFile() instead which returns S3 key separately
func UploadFileLegacy(ctx context.Context, fileData []byte, fileName string, contentType string, folder string) (string, error) {
	result, err := UploadFile(ctx, fileData, fileName, contentType, folder)
	if err != nil {
		return "", err
	}
	// Return legacy URL format for backward compatibility
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", S3BucketName, S3Region, result.S3Key)
	return url, nil
}

// GetPresignedURL generates a presigned URL for downloading a file
func GetPresignedURL(ctx context.Context, s3Key string, expiration time.Duration) (string, error) {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return "", fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	// Validate S3 key
	if s3Key == "" {
		return "", fmt.Errorf("S3 key cannot be empty")
	}

	// Verify object exists (optional check - can be removed if it causes performance issues)
	// This helps identify permission issues early
	// Note: We don't fail - presigned URL might still work even if HeadObject fails
	_, err := S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		// Presigned URL generation might still succeed even if HeadObject fails
		// Continue without logging to avoid noise
	}

	presignClient := s3.NewPresignClient(S3Client)
	
	// Generate presigned URL with response headers for CORS support
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(s3Key),
		// Add response headers for CORS support
		ResponseCacheControl:       aws.String("public, max-age=3600"),
		ResponseContentDisposition: nil, // Let browser handle disposition
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL (bucket: %s, key: %s): %w. Check AWS IAM permissions for s3:GetObject", S3BucketName, s3Key, err)
	}

	return request.URL, nil
}

// DeleteFile deletes a file from S3
func DeleteFile(ctx context.Context, s3Key string) error {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	_, err := S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// GetS3KeyFromURL extracts the S3 key from a full S3 URL
func GetS3KeyFromURL(s3URL string) string {
	// Handle presigned URLs - extract key before query parameters
	// Format: https://bucket.s3.region.amazonaws.com/key?X-Amz-Algorithm=...
	if strings.Contains(s3URL, "?") {
		s3URL = strings.Split(s3URL, "?")[0]
	}
	
	// Extract key from URL like: https://bucket.s3.region.amazonaws.com/key
	parts := strings.Split(s3URL, ".amazonaws.com/")
	if len(parts) > 1 {
		key := parts[1]
		// URL decode the key in case it was encoded
		decodedKey, err := url.QueryUnescape(key)
		if err == nil {
			return decodedKey
		}
		return key
	}
	
	// Try alternative format: https://s3.region.amazonaws.com/bucket/key
	if strings.Contains(s3URL, "/"+S3BucketName+"/") {
		parts := strings.Split(s3URL, "/"+S3BucketName+"/")
		if len(parts) > 1 {
			key := parts[1]
			decodedKey, err := url.QueryUnescape(key)
			if err == nil {
				return decodedKey
			}
			return key
		}
	}
	
	return ""
}

// GetObjectMetadata retrieves metadata for an S3 object
func GetObjectMetadata(ctx context.Context, s3Key string) (map[string]string, error) {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return nil, fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	result, err := S3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}

	metadata := make(map[string]string)
	if result.Metadata != nil {
		for key, value := range result.Metadata {
			metadata[key] = value
		}
	}

	return metadata, nil
}

// GetOriginalFilename retrieves the original filename from S3 object metadata
func GetOriginalFilename(ctx context.Context, s3Key string) string {
	metadata, err := GetObjectMetadata(ctx, s3Key)
	if err != nil {
		return ""
	}
	return metadata["original-filename"]
}

// GetFileTypeFromContentType determines file type category from content type
func GetFileTypeFromContentType(contentType string) string {
	// Normalize content type
	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	contentType = strings.TrimSpace(contentType)

	if strings.HasPrefix(contentType, "image/") {
		return "image"
	} else if strings.HasPrefix(contentType, "video/") {
		return "video"
	} else if strings.HasPrefix(contentType, "audio/") {
		return "audio"
	} else if strings.Contains(contentType, "pdf") ||
		strings.Contains(contentType, "word") ||
		strings.Contains(contentType, "excel") ||
		strings.Contains(contentType, "powerpoint") ||
		strings.Contains(contentType, "spreadsheet") ||
		strings.Contains(contentType, "presentation") {
		return "file"
	}
	return "file"
}

// GetFolderFromFileType returns the S3 folder based on file type
func GetFolderFromFileType(fileType string) string {
	switch fileType {
	case "image":
		return "images"
	case "video":
		return "videos"
	case "audio":
		return "audio"
	case "file":
		return "files"
	default:
		return "files"
	}
}

// ValidateFileType checks if the file type is allowed
func ValidateFileType(contentType string) bool {
	allowedTypes := []string{
		// Images
		"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp", "image/bmp", "image/svg+xml",
		// Videos
		"video/mp4", "video/mpeg", "video/quicktime", "video/x-msvideo", "video/x-ms-wmv",
		"video/webm", "video/ogg", "video/x-matroska",
		// Audio
		"audio/mpeg", "audio/mp3", "audio/wav", "audio/ogg", "audio/webm", "audio/aac",
		"audio/x-m4a", "audio/flac", "audio/x-wav",
		// Documents
		"application/pdf",
		// Office documents (optional)
		"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint", "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	}

	// Normalize content type (remove charset, etc.)
	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	contentType = strings.TrimSpace(contentType)

	for _, allowed := range allowedTypes {
		if contentType == allowed {
			return true
		}
	}

	return false
}

// ValidateFileSize checks if the file size is within allowed limits
func ValidateFileSize(size int64, fileType string) error {
	var maxSize int64

	switch fileType {
	case "image":
		maxSize = 10 * 1024 * 1024 // 10 MB for images
	case "video":
		maxSize = 500 * 1024 * 1024 // 500 MB for videos
	case "audio":
		maxSize = 50 * 1024 * 1024 // 50 MB for audio
	case "file":
		maxSize = 100 * 1024 * 1024 // 100 MB for PDFs and other files
	default:
		maxSize = 100 * 1024 * 1024 // 100 MB default
	}

	if size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d MB", maxSize/(1024*1024))
	}

	return nil
}
