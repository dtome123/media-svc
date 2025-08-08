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
	rabbitClient *rabbitmq.Publisher
}

func NewService(
	cfg *config.Config,
	db *mongodb.Database,
	rabbitClient *rabbitmq.Publisher,
) *Service {

	mediaRepo := media.NewMediaRepository(db)
	mediaStorage, _ := minio.New(cfg, cfg.S3.Bucket)
	streamStorage, _ := minio.New(cfg, cfg.S3.StreamBucket)

	return &Service{
		cfg:          cfg,
		rabbitClient: rabbitClient,
		mediaSvc:     mediaSvc.NewService(cfg, mediaRepo, mediaStorage, streamStorage, rabbitClient),
	}
}

func (i *Service) GetMediaSvc() mediaSvc.MediaService {
	return i.mediaSvc
}
