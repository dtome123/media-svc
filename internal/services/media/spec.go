package media

import (
	"context"
	"media-svc/internal/models"
)

type MediaService interface {
	CreateMedia(ctx context.Context, input CreateMediaInput) (string, error)
	UpdateMedia(ctx context.Context, input UpdateMediaInput) error
	GetMedia(ctx context.Context, id string) (*models.Media, error)
	ListMedia(ctx context.Context, input ListMediaInput) ([]*models.Media, error)

	UploadMedia(ctx context.Context, input UploadMediaInput) (string, error)
}
