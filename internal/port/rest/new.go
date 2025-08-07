package rest

import (
	"log"
	"media-svc/config"
	"media-svc/internal/port/rest/handlers"
	"media-svc/internal/port/rest/routes"
	"media-svc/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RestServer struct {
	cfg *config.Config
	svc *services.Service
}

func NewRestServer(cfg *config.Config, svc *services.Service) *RestServer {
	return &RestServer{
		cfg: cfg,
		svc: svc,
	}
}

func (s *RestServer) Run() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	handler := handlers.NewHandler(s.svc)

	routes.RegisterV1Routes(r.Group("/"), handler)

	log.Println("🚀 REST server running at :", s.cfg.Server.HttpPort)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
