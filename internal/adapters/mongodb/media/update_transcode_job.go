package media

import (
	"context"
	"media-svc/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *MediaRepository) UpdateTranscodeJob(ctx context.Context, transcode *models.TranscodeJob) error {

	transcode.BeforeUpdate()

	err := repo.transcodeJobCol.UpdateOne(ctx, bson.M{
		"media_id": transcode.MediaID,
	}, bson.M{"$set": transcode}, options.Update())
	return err
}
