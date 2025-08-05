package media

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type UploadMediaInput struct {
	OutputDir    string
	InputPath    string
	ObjectPrefix string
}

func (i *impl) UploadMedia(ctx context.Context, input UploadMediaInput) (string, error) {

	// Transcode
	if err := TranscodeToHLS(input.InputPath, input.OutputDir); err != nil {
		return "", err
	}

	err := i.storageAdapter.UploadDir(input.OutputDir, input.ObjectPrefix)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://localhost:9000/videos/%soutput.m3u8", input.ObjectPrefix), nil
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
