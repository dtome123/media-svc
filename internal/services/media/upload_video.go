package media

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"media-svc/internal/models"
	"media-svc/internal/types"
	"mime/multipart"
	"path/filepath"
	"time"
)

type UploadVideoInput struct {
	File *multipart.FileHeader
}

func (i *impl) UploadVideo(ctx context.Context, input UploadVideoInput) (string, error) {
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(input.File.Filename))

	filePath := filepath.Join("videos", filename)

	src, err := input.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	filePath, err = i.storageAdapter.PutObject(ctx, filePath, src, input.File.Size)
	if err != nil {
		return "", err
	}

	media := models.Media{
		Name:        input.File.Filename,
		Description: input.File.Filename,
		Path:        filePath,
		ContentType: input.File.Header.Get("Content-Type"),
	}
	err = i.mediaRepo.CreateMedia(ctx, media)
	if err != nil {
		return "", err
	}

	job := types.TranscodeJob{
		InputPath: filePath,
	}

	data, err := json.Marshal(job)
	if err != nil {
		return "", err
	}

	err = i.rabbitClient.Publish(i.cfg.RabbitMQ.Queue, data)
	if err != nil {
		log.Println("Failed to publish job:", err)
		return "", err
	}

	return filePath, nil
}
