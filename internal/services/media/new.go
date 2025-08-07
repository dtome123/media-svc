package media

import (
	"media-svc/config"
	"media-svc/internal/adapters/minio"
	"media-svc/internal/adapters/mongodb/media"
	"media-svc/pkgs/rabbitmq"
)

type impl struct {
	cfg            *config.Config
	mediaRepo      *media.MediaRepository
	storageAdapter minio.StorageAdapter
	rabbitClient   *rabbitmq.Client
}

func NewService(cfg *config.Config, mediaRepo *media.MediaRepository, storageAdapter minio.StorageAdapter, rabbitClient *rabbitmq.Client) MediaService {
	return &impl{
		cfg:            cfg,
		mediaRepo:      mediaRepo,
		storageAdapter: storageAdapter,
		rabbitClient:   rabbitClient,
	}
}
