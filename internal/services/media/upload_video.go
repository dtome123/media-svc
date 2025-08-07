package media

import (
	"context"
	"fmt"
	"media-svc/internal/models"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type UploadVideoInput struct {
	File *multipart.FileHeader
}

func (i *impl) UploadVideo(ctx context.Context, input UploadVideoInput) (string, error) {

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(input.File.Filename))
	// baseName := filename[:len(filename)-len(filepath.Ext(filename))]
	// objectDir := filepath.Join("videos", baseName) + "/"

	filePath := filepath.Join("videos", filename)

	src, err := input.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	// inputPath := filepath.Join("uploads", filename)

	// Transcode
	// if err := TranscodeToHLS(input.InputPath, input.OutputDir); err != nil {
	// 	return "", err
	// }

	// err := i.storageAdapter.UploadDir(context.Background(), input.OutputDir, input.ObjectDir)
	// if err != nil {
	// 	return "", err
	// }

	// return fmt.Sprintf("http://localhost:9000/videos/%soutput.m3u8", input.ObjectDir), nil

	err = i.storageAdapter.PutObject(ctx, filePath, src, input.File.Size)
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

	return filePath, nil
}

func TranscodeToHLS(inputPath, outputDir string) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	outputPattern := filepath.Join(outputDir, "output.m3u8")

	return ffmpeg.Input(inputPath).
		Output(outputPattern,
			ffmpeg.KwArgs{
				"codec":         "copy",
				"start_number":  "0",
				"hls_time":      "10",
				"hls_list_size": "0",
				"f":             "hls",
			},
		).
		OverWriteOutput().
		Run()
}
