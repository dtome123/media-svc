package media

import (
	"context"
	"media-svc/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (svc *MediaRepository) GetTranscodeJobByMediaID(ctx context.Context, mediaId string) (*models.TranscodeJob, error) {

	oid, err := primitive.ObjectIDFromHex(mediaId)
	if err != nil {
		return nil, err
	}

	model, err := svc.transcodeJobCol.FindOne(ctx, bson.M{
		"media_id": oid,
	}, options.FindOne().SetHint(IndexTranscodeJobMediaID))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return model, nil
}
