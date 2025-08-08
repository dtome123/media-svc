package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetMediaRequest struct {
	MediaID string `uri:"media_id"`
}

func (s *impl) GetMedia(c *gin.Context) {

	var req GetMediaRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services := s.svc.GetMediaSvc()
	media, err := services.GetMedia(c, req.MediaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	c.JSON(http.StatusOK, Media{
		ID:          media.ID.Hex(),
		Name:        media.Name,
		Description: media.Description,
		Path:        media.Path,
		Size:        media.Size,
		ContentType: media.ContentType,
	})
}
