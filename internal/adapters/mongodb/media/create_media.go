package media

import (
	"context"
	"media-svc/internal/models"
)

func (repo *MediaRepository) CreateMedia(ctx context.Context, media models.Media) error {
	err := repo.mediaCol.InsertOne(ctx, media)
	return err
}
