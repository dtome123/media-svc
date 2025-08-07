package media

import (
	"context"
	"media-svc/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateMediaInput struct {
	ID           string
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

func (i *impl) UpdateMedia(ctx context.Context, input UpdateMediaInput) error {

	objectID, _ := primitive.ObjectIDFromHex(input.ID)

	media := models.Media{
		ID:           objectID,
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
		UpdatedAt:    time.Now().UTC(),
	}

	err := i.mediaRepo.UpdateMedia(ctx, media)
	if err != nil {
		return err
	}

	return nil
}
