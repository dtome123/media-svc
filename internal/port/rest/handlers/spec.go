package handlers

import "github.com/gin-gonic/gin"

type Handler interface {
	UploadVideo(c *gin.Context)
	Stream(c *gin.Context)
	GetVideoStatus(c *gin.Context)
}
