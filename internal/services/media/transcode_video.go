package media

import (
	"context"
	"fmt"
	"media-svc/internal/utils"
	"os"
	"os/exec"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type TranscodeVideoInput struct {
	FilePath string
}

func (i *impl) TranscodeVideo(ctx context.Context, input TranscodeVideoInput) (string, error) {

	src, err := i.storageAdapter.GetObject(ctx, input.FilePath)
	if err != nil {
		return "", err
	}
	filename := filepath.Base(input.FilePath)

	// save file local
	err = utils.WriteFile("assets", filename, src)
	if err != nil {
		return "", err
	}
	localFilePath := fmt.Sprintf("assets/%s", filename)
	outputDir := fmt.Sprintf("assets/transcode/%s", filename)
	if err := TranscodeAdaptive(localFilePath, outputDir); err != nil {
		return "", err
	}

	targetDir := fmt.Sprintf("videos/%s", filename)

	dirPath, err := i.storageAdapter.UploadStreamDir(context.Background(), outputDir, targetDir)
	if err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s/master.m3u8", dirPath)

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

func TranscodeAdaptive(inputPath, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-filter_complex",
		"[0:v]split=3[v1][v2][v3];"+
			"[v1]scale=w=1920:h=1080:force_original_aspect_ratio=decrease[v1out];"+
			"[v2]scale=w=1280:h=720:force_original_aspect_ratio=decrease[v2out];"+
			"[v3]scale=w=640:h=360:force_original_aspect_ratio=decrease[v3out]",
		// 1080p
		"-map", "[v1out]", "-map", "a:0?", "-c:v:0", "libx264", "-b:v:0", "5000k",
		"-c:a:0", "aac", "-b:a:0", "192k",
		"-f", "hls", "-hls_time", "4", "-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputDir, "1080p_%03d.ts"),
		filepath.Join(outputDir, "1080p.m3u8"),

		// 720p
		"-map", "[v2out]", "-map", "a:0?", "-c:v:1", "libx264", "-b:v:1", "3000k",
		"-c:a:1", "aac", "-b:a:1", "128k",
		"-f", "hls", "-hls_time", "4", "-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputDir, "720p_%03d.ts"),
		filepath.Join(outputDir, "720p.m3u8"),

		// 360p
		"-map", "[v3out]", "-map", "a:0?", "-c:v:2", "libx264", "-b:v:2", "1000k",
		"-c:a:2", "aac", "-b:a:2", "96k",
		"-f", "hls", "-hls_time", "4", "-hls_playlist_type", "vod",
		"-hls_segment_filename", filepath.Join(outputDir, "360p_%03d.ts"),
		filepath.Join(outputDir, "360p.m3u8"),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %v\nOutput: %s", err, string(out))
	}

	master := `#EXTM3U
#EXT-X-STREAM-INF:BANDWIDTH=5000000,RESOLUTION=1920x1080
1080p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=3000000,RESOLUTION=1280x720
720p.m3u8
#EXT-X-STREAM-INF:BANDWIDTH=1000000,RESOLUTION=640x360
360p.m3u8
`

	err = os.WriteFile(filepath.Join(outputDir, "master.m3u8"), []byte(master), 0644)
	if err != nil {
		return fmt.Errorf("failed to write master playlist: %v", err)
	}

	return nil
}
