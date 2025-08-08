package media

import (
	"context"
	"media-svc/internal/types"
)

type UpdateTranscodeJobErrorInput struct {
	MediaID string
	Err     string
}

func (i *impl) UpdateTranscodeJobError(ctx context.Context, input UpdateTranscodeJobErrorInput) error {

	job, err := i.mediaRepo.GetTranscodeJobByMediaID(ctx, input.MediaID)
	if err != nil {
		return err
	}

	if job == nil {
		return nil
	}

	job.Status = types.TranscodeJobStatusError.String()
	job.Error = input.Err

	err = i.mediaRepo.UpdateTranscodeJob(ctx, job)
	if err != nil {
		return err
	}

	return nil
}
