package media

import (
	"context"
	"time"
)

type PresignGetObjectInput struct {
	FilePath string
}

func (i *impl) PresignGetStreamObject(ctx context.Context, input PresignGetObjectInput) (string, error) {

	presignUrl, err := i.streamStorage.PresignGetObject(ctx, input.FilePath, 5*time.Minute)
	if err != nil {
		return "", err
	}

	return presignUrl, nil
}
