package media

import (
	"context"
	"media-svc/internal/types"
)

type GetVideoStatusResponse struct {
	Status          string
	TranscodeSource string
}

func (i *impl) GetVideoStatus(ctx context.Context, videoId string) (GetVideoStatusResponse, error) {

	media, err := i.mediaRepo.GetMedia(ctx, videoId)
	if err != nil {
		return GetVideoStatusResponse{}, err
	}

	var source string

	if media.TranscodeSource != nil {
		source = media.TranscodeSource.FilePath
	}

	transcodeJob, err := i.mediaRepo.GetTranscodeJobByMediaID(ctx, videoId)
	if err != nil {
		return GetVideoStatusResponse{}, err
	}

	if transcodeJob == nil {
		return GetVideoStatusResponse{
			Status: types.TranscodeJobStatusPending.String(),
		}, nil
	}

	return GetVideoStatusResponse{
		Status:          transcodeJob.Status,
		TranscodeSource: source,
	}, nil
}
