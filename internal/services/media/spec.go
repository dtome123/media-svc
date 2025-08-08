package media

import (
	"context"
	"media-svc/internal/models"
)

type MediaService interface {
	CreateMedia(ctx context.Context, input CreateMediaInput) (string, error)
	UpdateMedia(ctx context.Context, input UpdateMediaInput) error
	GetMedia(ctx context.Context, id string) (*models.Media, error)
	ListMedia(ctx context.Context, input ListMediaInput) ([]*models.Media, error)

	PresignGetStreamObject(ctx context.Context, input PresignGetObjectInput) (string, error)
	UploadVideo(ctx context.Context, input UploadVideoInput) (*models.Media, error)
	TranscodeVideo(ctx context.Context, input TranscodeVideoInput) (string, error)
	UpdateTranscodeJobError(ctx context.Context, input UpdateTranscodeJobErrorInput) error
	UpdateTranscodeJobSuccess(ctx context.Context, input UpdateTranscodeJobSuccessInput) error
	GetVideoStatus(ctx context.Context, videoId string) (GetVideoStatusResponse, error)
}
