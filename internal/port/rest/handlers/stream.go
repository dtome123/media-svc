package handlers

import (
	"io"
	"log"
	"media-svc/internal/services/media"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *impl) Stream(c *gin.Context) {

	filePath := c.Param("file_path")
	filePath = strings.TrimPrefix(filePath, "/")

	services := s.svc.GetMediaSvc()
	presignedURL, err := services.PresignGetStreamObject(c, media.PresignGetObjectInput{
		FilePath: filePath,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	resp, err := http.Get(presignedURL)
	if err != nil {
		c.String(http.StatusBadGateway, "Failed to get object from MinIO: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.String(resp.StatusCode, "MinIO responded with status: %v", resp.Status)
		return
	}

	// Copy header content-type từ MinIO sang client
	for k, vv := range resp.Header {
		for _, v := range vv {
			c.Writer.Header().Add(k, v)
		}
	}

	// Stream dữ liệu về client
	c.Status(http.StatusOK)
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		log.Printf("Error streaming object. Error: %v", err)
	}
}
