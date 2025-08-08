package routes

import (
	"media-svc/internal/port/rest/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(r *gin.RouterGroup, handler handlers.Handler) {
	v1 := r.Group("v1")
	v1VideoRoutes(v1, handler)
}

func v1VideoRoutes(r *gin.RouterGroup, handler handlers.Handler) {
	videoRoutes := r.Group("videos")
	videoRoutes.POST("/upload", handler.UploadVideo)
	videoRoutes.GET("/:video_id/status", handler.GetVideoStatus)
	videoRoutes.GET("/stream/*file_path", handler.Stream)
}
