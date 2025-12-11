package services

import (
	"bytes"
	"context"
	"fmt"
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

// UploadFile uploads a file to S3 and returns the S3 URL
func UploadFile(ctx context.Context, fileData []byte, fileName string, contentType string, folder string) (string, error) {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return "", fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	// Generate unique filename
	ext := filepath.Ext(fileName)
	uniqueFileName := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	// Upload file to S3 with Glacier Instant Retrieval storage class
	storageClass := types.StorageClassGlacierIr // Use Glacier Instant Retrieval
	putInput := &s3.PutObjectInput{
		Bucket:       aws.String(S3BucketName),
		Key:          aws.String(uniqueFileName),
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
	// If you need public access, configure the bucket policy in AWS Console

	_, err := S3Uploader.Upload(ctx, putInput)
	if err != nil {
		// Return detailed error for debugging
		return "", fmt.Errorf("S3 upload failed (bucket: %s, key: %s): %w", S3BucketName, uniqueFileName, err)
	}

	// Return the public URL (format: https://bucket-name.s3.region.amazonaws.com/key)
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", S3BucketName, S3Region, uniqueFileName)
	return url, nil
}

// GetPresignedURL generates a presigned URL for downloading a file
func GetPresignedURL(ctx context.Context, s3Key string, expiration time.Duration) (string, error) {
	if S3Client == nil {
		if err := InitializeS3(); err != nil {
			return "", fmt.Errorf("failed to initialize S3: %w", err)
		}
	}

	presignClient := s3.NewPresignClient(S3Client)
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(S3BucketName),
		Key:    aws.String(s3Key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
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
func GetS3KeyFromURL(url string) string {
	// Extract key from URL like: https://bucket.s3.region.amazonaws.com/key
	parts := strings.Split(url, ".amazonaws.com/")
	if len(parts) > 1 {
		return parts[1]
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
