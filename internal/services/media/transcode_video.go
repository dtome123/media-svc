package media

import (
	"context"
	"fmt"
	"media-svc/internal/models"
	"media-svc/internal/types"
	"media-svc/internal/utils"
	"media-svc/pkgs/transcoder"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TranscodeVideoInput struct {
	MediaID string
}

type TranscodeVideoOutput struct {
	Path       string
	Renditions []types.Rendition
}

// TranscodeVideo downloads a video file, transcodes it into adaptive streams,
// uploads the transcoded files back to storage, and returns the master playlist path.
func (i *impl) TranscodeVideo(ctx context.Context, input TranscodeVideoInput) (TranscodeVideoOutput, error) {

	media, err := i.mediaRepo.GetMedia(ctx, input.MediaID)
	if err != nil {
		return TranscodeVideoOutput{}, fmt.Errorf("get media: %w", err)
	}

	if media == nil {
		return TranscodeVideoOutput{}, fmt.Errorf("media not found")
	}

	filePath := media.Path

	// Download video file from storage
	src, err := i.mediaStorage.GetObject(ctx, filePath)
	if err != nil {
		return TranscodeVideoOutput{}, fmt.Errorf("get object from storage: %w", err)
	}

	filename := filepath.Base(filePath)

	// Save downloaded file locally
	if err := utils.WriteFile("assets", filename, src); err != nil {
		return TranscodeVideoOutput{}, fmt.Errorf("write file local: %w", err)
	}

	// create transcode job db
	mediaObjectId, _ := primitive.ObjectIDFromHex(input.MediaID)
	now := time.Now().UTC()
	job := &models.TranscodeJob{
		MediaID:   mediaObjectId,
		Status:    types.TranscodeJobStatusProcessing.String(),
		StartedAt: &now,
	}
	err = i.mediaRepo.CreateTranscodeJob(ctx, job)
	if err != nil {
		return TranscodeVideoOutput{}, fmt.Errorf("create transcode job: %w", err)
	}

	localFilePath := filepath.Join("assets", filename)
	outputDir := filepath.Join("assets", "transcode", filename)

	// Transcode the video into adaptive bitrate streams using ffmpeg
	transcoder := transcoder.New()
	renditions, err := transcoder.TranscodeAdaptiveCMAF(localFilePath, outputDir)
	if err != nil {
		return TranscodeVideoOutput{}, fmt.Errorf("transcode adaptive: %w", err)
	}

	// Upload the transcoded directory back to storage
	targetDir := filepath.Join(filename)
	dirPath, err := i.streamStorage.UploadDir(context.Background(), outputDir, targetDir)
	if err != nil {
		return TranscodeVideoOutput{}, fmt.Errorf("upload transcode dir: %w", err)
	}

	// Construct and return the path to the master playlist file
	filePath = filepath.Join(dirPath, "master.m3u8")

	// remove local files
	if err := utils.RemoveFile(localFilePath); err != nil {
		return TranscodeVideoOutput{}, fmt.Errorf("remove local file: %w", err)
	}
	if err := utils.RemoveDir(outputDir); err != nil {
		return TranscodeVideoOutput{}, fmt.Errorf("remove output dir: %w", err)
	}

	var outRenditions []types.Rendition
	for _, r := range renditions {
		outRenditions = append(outRenditions, types.Rendition{
			Width:        r.Width,
			Height:       r.Height,
			Name:         r.Name,
			VideoBitrate: r.VideoBitrate,
			AudioBitrate: r.AudioBitrate,
		})
	}

	return TranscodeVideoOutput{
		Path:       filePath,
		Renditions: outRenditions,
	}, nil
}
