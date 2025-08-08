package transcoder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// Rendition defines one output quality profile for transcoding.
type Rendition struct {
	Name         string // Rendition name, e.g., "1080p"
	Width        int    // Target width
	Height       int    // Target height
	VideoBitrate string // Video bitrate string, e.g., "5000k"
	AudioBitrate string // Audio bitrate string, e.g., "192k"
}

// DefaultRenditions contains common adaptive streaming resolutions and bitrates.
var DefaultRenditions = []Rendition{
	{"1080p", 1920, 1080, "5000k", "192k"},
	{"720p", 1280, 720, "3000k", "128k"},
	{"360p", 640, 360, "1000k", "96k"},
}

// Transcoder provides methods to perform video transcoding.
type Transcoder struct{}

// New returns a new Transcoder instance.
func New() *Transcoder {
	return &Transcoder{}
}

// TranscodeAdaptiveCMAF performs adaptive bitrate transcoding using CMAF segments,
// automatically selecting output renditions based on the input video's resolution.
//
// inputPath: path to source video file.
// outputDir: directory where transcoded files will be stored.
//
// Returns the list of renditions created, or an error if transcoding fails.
func (t *Transcoder) TranscodeAdaptiveCMAF(inputPath, outputDir string) ([]Rendition, error) {
	// Get input video resolution using ffprobe.
	srcW, srcH, err := getVideoResolution(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get source resolution: %w", err)
	}

	// Select renditions that are smaller or equal to source resolution.
	selected := filterRenditions(srcW, srcH, DefaultRenditions)

	// Create output directory and variant subdirectories.
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}
	for i := range selected {
		variantDir := filepath.Join(outputDir, fmt.Sprintf("%d", i))
		if err := os.MkdirAll(variantDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create variant directory: %w", err)
		}
	}

	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return nil, err
	}

	// Build ffmpeg filter_complex argument for splitting and scaling.
	filterComplex := buildFilterComplex(selected)

	// Relative segment and playlist patterns for HLS.
	hlsSegmentPattern := filepath.Join("%v", "seg_%03d.m4s")
	hlsPlaylistPattern := filepath.Join("%v", "stream.m3u8")
	masterPlaylist := "master.m3u8"
	dashManifest := "manifest.mpd"

	// Construct the var_stream_map argument for ffmpeg.
	varStreamMap := buildVarStreamMap(selected)

	// Build the full ffmpeg command-line arguments.
	args := buildFFmpegArgs(absInputPath, filterComplex, selected, hlsSegmentPattern, hlsPlaylistPattern, masterPlaylist, varStreamMap, dashManifest)

	// Create and execute ffmpeg command.
	cmd := exec.Command("ffmpeg", args...)
	cmd.Dir = outputDir // Set working directory for output files

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ffmpeg error: %v\nOutput:\n%s", err, string(output))
		return nil, fmt.Errorf("ffmpeg failed: %v", err)
	}

	return selected, nil
}
