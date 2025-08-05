package handlers

import (
	"fmt"
	"media-svc/internal/services/media"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *impl) UploadVideo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
	baseName := filename[:len(filename)-len(filepath.Ext(filename))]
	objectPrefix := filepath.Join("videos", baseName) + "/"
	inputPath := filepath.Join("uploads", filename)
	outputDir := filepath.Join("output", baseName)

	if err := c.SaveUploadedFile(file, inputPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	services := s.svc.GetMediaSvc()

	m3u8URL, err := services.UploadMedia(c, media.UploadMediaInput{
		OutputDir:    outputDir,
		InputPath:    inputPath,
		ObjectPrefix: objectPrefix,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Uploaded and transcoded successfully",
		"streamUrl": m3u8URL,
	})
}
