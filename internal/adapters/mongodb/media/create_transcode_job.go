package media

import (
	"context"
	"media-svc/internal/models"
)

func (repo *MediaRepository) CreateTranscodeJob(ctx context.Context, transcode models.TranscodeJob) error {
	err := repo.transcodeJobCol.InsertOne(ctx, transcode)
	return err
}
