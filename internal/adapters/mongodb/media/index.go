package media

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	IndexTranscodeJobMediaID = "transcode_job_media_id"
)

func GetMediaIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{}
}

func GetTranscodeJobIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.M{
				"media_id": 1,
			},
			Options: options.Index().SetName(IndexTranscodeJobMediaID),
		},
	}
}
