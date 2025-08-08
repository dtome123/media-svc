package media

import (
	"media-svc/config"
	"media-svc/internal/adapters/minio"
	"media-svc/internal/adapters/mongodb/media"
	"media-svc/pkgs/rabbitmq"
)

type impl struct {
	cfg           *config.Config
	mediaRepo     *media.MediaRepository
	mediaStorage  minio.StorageAdapter
	streamStorage minio.StorageAdapter
	rabbitClient  *rabbitmq.Publisher
}

func NewService(
	cfg *config.Config,
	mediaRepo *media.MediaRepository,
	mediaStorage minio.StorageAdapter,
	streamStorage minio.StorageAdapter,
	rabbitClient *rabbitmq.Publisher,
) MediaService {
	return &impl{
		cfg:           cfg,
		mediaRepo:     mediaRepo,
		mediaStorage:  mediaStorage,
		streamStorage: streamStorage,
		rabbitClient:  rabbitClient,
	}
}
