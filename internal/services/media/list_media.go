package media

import (
	"context"
	"media-svc/internal/adapters/mongodb/media"
	"media-svc/internal/models"
)

type ListMediaInput struct {
	Keyword string
}

func (i *impl) ListMedia(ctx context.Context, input ListMediaInput) ([]*models.Media, error) {

	medias, err := i.mediaRepo.ListMedia(ctx, media.ListMediaInput{
		Keyword: input.Keyword,
	})
	if err != nil {
		return nil, err
	}

	return medias, nil
}
