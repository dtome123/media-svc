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
	Renditions []types.Rendition
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

	var renditions []models.Rendition
	for _, rendition := range input.Renditions {
		renditions = append(renditions, models.Rendition{
			Width:        rendition.Width,
			Height:       rendition.Height,
			Name:         rendition.Name,
			VideoBitrate: rendition.VideoBitrate,
			AudioBitrate: rendition.AudioBitrate,
		})
	}

	media.TranscodeSource = &models.TranscodeSource{
		FilePath:   input.OutputPath,
		Renditions: renditions,
	}

	err = i.mediaRepo.UpdateMedia(ctx, media)
	if err != nil {
		return err
	}

	return nil
}
