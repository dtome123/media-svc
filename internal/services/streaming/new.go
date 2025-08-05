package streaming

import (
	"media-svc/config"
	"media-svc/internal/adapters/mongodb/media"
)

type impl struct {
	cfg       *config.Config
	mediaRepo *media.MediaRepository
}

func NewService(cfg *config.Config, mediaRepo *media.MediaRepository) StreamingService {
	return &impl{
		cfg:       cfg,
		mediaRepo: mediaRepo,
	}
}
