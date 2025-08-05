package services

import (
	"media-svc/config"
	"media-svc/internal/adapters/minio"
	"media-svc/internal/adapters/mongodb/media"
	mediaSvc "media-svc/internal/services/media"

	mongodb "github.com/dtome123/go-mongo-generic"
)

type Service struct {
	cfg      *config.Config
	mediaSvc mediaSvc.MediaService
}

func NewService(
	cfg *config.Config,
	db *mongodb.Database,
) *Service {

	mediaRepo := media.NewMediaRepository(db)
	storageAdapter := minio.New(cfg)

	return &Service{
		cfg:      cfg,
		mediaSvc: mediaSvc.NewService(cfg, mediaRepo, storageAdapter),
	}
}

func (i *Service) GetMediaSvc() mediaSvc.MediaService {
	return i.mediaSvc
}
