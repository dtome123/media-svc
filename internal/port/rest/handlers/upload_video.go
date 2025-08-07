package handlers

import (
	"media-svc/internal/services/media"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *impl) UploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	services := s.svc.GetMediaSvc()
	_, err = services.UploadVideo(c, media.UploadVideoInput{
		File: file,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Uploaded and transcoded successfully",
	})
}
