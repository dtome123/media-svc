package media

import (
	"context"
	"media-svc/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ListMediaInput struct {
	Keyword string
}

func (svc *MediaRepository) ListMedia(ctx context.Context, input ListMediaInput) ([]*models.Media, error) {

	filter := bson.M{}
	if input.Keyword != "" {
		filter["name"] = input.Keyword
	}

	medias, err := svc.mediaCol.Find(
		ctx,
		filter,
		options.Find().SetSort(bson.M{"_id": -1}),
		nil)

	if err != nil {
		return nil, err
	}

	return medias, nil
}
