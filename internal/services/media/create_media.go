package media

import (
	"context"
	"media-svc/internal/models"
	"time"
)

type CreateMediaInput struct {
	Name         string
	Description  string
	Path         string
	ContentType  string
	Size         int64
	Duration     float64
	Width        int
	Height       int
	IsStreamable bool
	Tags         []string
}

func (i *impl) CreateMedia(ctx context.Context, input CreateMediaInput) (string, error) {
	media := &models.Media{
		Name:         input.Name,
		Description:  input.Description,
		Path:         input.Path,
		ContentType:  input.ContentType,
		Size:         input.Size,
		Duration:     input.Duration,
		Width:        input.Width,
		Height:       input.Height,
		IsStreamable: input.IsStreamable,
		Tags:         input.Tags,
		CreatedAt:    time.Now().UTC(),
	}

	err := i.mediaRepo.CreateMedia(ctx, media)
	if err != nil {
		return "", err
	}

	return media.ID.Hex(), nil
}
