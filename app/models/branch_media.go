package models

import (
	"time"
)

// BranchMedia represents media files for a branch
type BranchMedia struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	BranchID        uint      `gorm:"not null" json:"branch_id"`
	IsChildBranch   bool      `gorm:"default:false" json:"is_child_branch"`
	FileURL         string    `json:"-" gorm:"column:file_url"` // Internal: NEVER serialize to JSON - stores presigned URL temporarily
	S3Key           string    `json:"s3_key,omitempty" gorm:"column:s3_key"`   // Opaque S3 object key (UUID-based)
	OriginalFilename string   `json:"original_filename,omitempty" gorm:"column:original_filename"` // Original filename from upload
	FileType        string    `json:"file_type,omitempty" gorm:"column:file_type"` // image, video, audio, file
	Name            string    `json:"name,omitempty"`
	URL             string    `json:"url,omitempty" gorm:"-"` // Computed: presigned URL (populated by ConvertBranchMediaToPresignedURLs)
	Category    string    `json:"category,omitempty"` // Branch Photos, Video Coverage, Documents, Other
	CreatedOn   time.Time `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn   time.Time `gorm:"autoUpdateTime" json:"updated_on"`
	CreatedBy   string    `json:"created_by,omitempty" gorm:"<-:create"`
	UpdatedBy   string    `json:"updated_by,omitempty"`
	Branch      Branch    `gorm:"foreignKey:BranchID;references:ID" json:"branch,omitempty"`
}

func (BranchMedia) TableName() string {
	return "branch_media"
}


