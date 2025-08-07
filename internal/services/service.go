package services

import (
	"media-svc/config"
	"media-svc/internal/adapters/minio"
	"media-svc/internal/adapters/mongodb/media"
	mediaSvc "media-svc/internal/services/media"
	"media-svc/pkgs/rabbitmq"

	mongodb "github.com/dtome123/go-mongo-generic"
)

type Service struct {
	cfg          *config.Config
	mediaSvc     mediaSvc.MediaService
	rabbitClient *rabbitmq.Client
}

func NewService(
	cfg *config.Config,
	db *mongodb.Database,
	rabbitClient *rabbitmq.Client,
) *Service {

	mediaRepo := media.NewMediaRepository(db)
	storageAdapter := minio.New(cfg)

	return &Service{
		cfg:          cfg,
		rabbitClient: rabbitClient,
		mediaSvc:     mediaSvc.NewService(cfg, mediaRepo, storageAdapter, rabbitClient),
	}
}

func (i *Service) GetMediaSvc() mediaSvc.MediaService {
	return i.mediaSvc
}
