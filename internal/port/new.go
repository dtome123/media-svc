package port

import (
	"log"
	"media-svc/config"
	"media-svc/internal/port/rest"
	"media-svc/internal/services"
	"runtime/debug"

	mongodb "github.com/dtome123/go-mongo-generic"
)

type Server struct {
	cfg *config.Config
	svc *services.Service
}

func NewServer(cfg *config.Config, db *mongodb.Database) *Server {

	return &Server{
		cfg: cfg,
		svc: services.NewService(cfg, db),
	}
}

func (s *Server) Run() {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("❗ Recovered from panic: %v\n%s", r, debug.Stack())
		}
	}()

	restSvr := rest.NewRestServer(s.cfg, s.svc)

	// Run HTTP in parallel
	go restSvr.Run()

	// Prevent main from exiting
	select {}
}
