package media

import (
	"context"
	"media-svc/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateMediaInput struct {
	ID          string
	Name        string
	Description string
	PlaybackUrl string
}

func (i *impl) UpdateMedia(ctx context.Context, input UpdateMediaInput) error {

	objectID, _ := primitive.ObjectIDFromHex(input.ID)

	media := models.Media{
		ID:          objectID,
		Name:        input.Name,
		Description: input.Description,
		PlaybackUrl: input.PlaybackUrl,
		UpdatedAt:   time.Now().UTC(),
	}

	err := i.mediaRepo.UpdateMedia(ctx, media)
	if err != nil {
		return err
	}

	return nil
}
