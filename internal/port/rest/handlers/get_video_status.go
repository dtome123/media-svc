package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetVideoStatusRequest struct {
	VideoID string `uri:"video_id"`
}

type GetVideoStatusResponse struct {
	Status          string `json:"status"`
	TranscodeSource string `json:"transcode_source"`
}

func (s *impl) GetVideoStatus(c *gin.Context) {

	var req GetVideoStatusRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services := s.svc.GetMediaSvc()
	res, err := services.GetVideoStatus(c, req.VideoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	c.JSON(http.StatusOK, GetVideoStatusResponse{
		Status:          res.Status,
		TranscodeSource: res.TranscodeSource,
	})
}
