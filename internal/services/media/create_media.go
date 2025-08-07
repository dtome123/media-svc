package media

import (
	"context"
	"media-svc/internal/models"
	"time"
)

type CreateMediaInput struct {
	Name        string
	Description string
	PlaybackUrl string
}

func (i *impl) CreateMedia(ctx context.Context, input CreateMediaInput) (string, error) {
	media := models.Media{
		Name:        input.Name,
		Description: input.Description,
		Path:        input.PlaybackUrl,
		CreatedAt:   time.Now().UTC(),
	}

	err := i.mediaRepo.CreateMedia(ctx, media)
	if err != nil {
		return "", err
	}

	return media.ID.Hex(), nil
}
