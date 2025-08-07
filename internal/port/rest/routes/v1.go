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
	group := r.Group("videos")
	group.POST("/upload", handler.UploadVideo)
	group.GET("/stream/*file_path", handler.Stream)
}
