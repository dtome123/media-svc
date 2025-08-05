package media

import (
	"media-svc/config"
	"media-svc/internal/adapters/minio"
	"media-svc/internal/adapters/mongodb/media"
)

type impl struct {
	cfg            *config.Config
	mediaRepo      *media.MediaRepository
	storageAdapter minio.StorageAdapter
}

func NewService(cfg *config.Config, mediaRepo *media.MediaRepository, storageAdapter minio.StorageAdapter) MediaService {
	return &impl{
		cfg:            cfg,
		mediaRepo:      mediaRepo,
		storageAdapter: storageAdapter,
	}
}
