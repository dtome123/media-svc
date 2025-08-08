package media

import (
	"context"
	"media-svc/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *MediaRepository) UpdateMedia(ctx context.Context, media *models.Media) error {
	
	media.BeforeUpdate()
	
	err := repo.mediaCol.UpdateOne(ctx, bson.M{
		"_id": media.ID,
	}, bson.M{"$set": media}, options.Update().SetHint("_id_"))
	return err
}
