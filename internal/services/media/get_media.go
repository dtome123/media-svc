package media

import (
	"context"
	"media-svc/internal/models"
)

func (i *impl) GetMedia(ctx context.Context, id string) (*models.Media, error) {

	media, err := i.mediaRepo.GetMedia(ctx, id)
	if err != nil {
		return nil, err
	}

	return media, nil
}
