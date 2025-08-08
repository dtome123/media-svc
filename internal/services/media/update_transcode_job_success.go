package media

import (
	"context"
	"media-svc/internal/models"
	"media-svc/internal/types"
	"time"
)

type UpdateTranscodeJobSuccessInput struct {
	MediaID    string
	OutputPath string
}

func (i *impl) UpdateTranscodeJobSuccess(ctx context.Context, input UpdateTranscodeJobSuccessInput) error {

	job, err := i.mediaRepo.GetTranscodeJobByMediaID(ctx, input.MediaID)
	if err != nil {
		return err
	}

	if job == nil {
		return nil
	}

	now := time.Now().UTC()
	job.Status = types.TranscodeJobStatusDone.String()
	job.OutputPath = input.OutputPath
	job.FinishedAt = &now

	err = i.mediaRepo.UpdateTranscodeJob(ctx, job)
	if err != nil {
		return err
	}

	media, err := i.mediaRepo.GetMedia(ctx, input.MediaID)
	if err != nil {
		return err
	}

	media.TranscodeSource = &models.TranscodeSource{
		FilePath: input.OutputPath,
	}

	i.mediaRepo.UpdateMedia(ctx, media)

	return nil
}
