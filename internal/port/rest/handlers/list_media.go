package handlers

import (
	"media-svc/internal/services/media"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListMediaRequest struct {
	MediaID string `uri:"media_id"`
}

func (s *impl) ListMedia(c *gin.Context) {

	var req ListMediaRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services := s.svc.GetMediaSvc()
	medias, err := services.ListMedia(c, media.ListMediaInput{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	var mediasResponse []Media

	for _, media := range medias {
		mediasResponse = append(mediasResponse, Media{
			ID:          media.ID.Hex(),
			Name:        media.Name,
			Description: media.Description,
			Path:        media.Path,
			Size:        media.Size,
			ContentType: media.ContentType,
		})
	}

	c.JSON(http.StatusOK, mediasResponse)
}
