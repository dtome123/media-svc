package media

import (
	"context"
	"media-svc/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (svc *MediaRepository) GetMedia(ctx context.Context, id string) (*models.Media, error) {

	media, err := svc.mediaCol.FindOne(ctx, bson.M{
		"_id": id,
	}, options.FindOne())

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return media, nil
}
