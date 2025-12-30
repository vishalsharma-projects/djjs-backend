package api

import (
	"github.com/followCode/djjs-event-reporting-backend/app/handlers"
	"github.com/followCode/djjs-event-reporting-backend/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupFileRoutes configures file upload/download routes
func SetupFileRoutes(r *gin.RouterGroup) {
	files := r.Group("/files")
	files.Use(middleware.AuthRequired())
	{
		files.POST("/upload", handlers.UploadFileHandler)
		files.POST("/upload-multiple", handlers.UploadMultipleFilesHandler)
		files.POST("/upload-branch", handlers.UploadBranchFilesHandler)
		files.GET("/:media_id/download", handlers.DownloadFileHandler)
		files.DELETE("/:media_id", handlers.DeleteFileHandler)
	}
}

