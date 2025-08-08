package transcoder

import (
	"fmt"
	"strings"
)

// filterRenditions filters the renditions list, returning only those
// with width and height less than or equal to the source resolution.
// If no renditions match, returns the smallest one by default.
func filterRenditions(srcW, srcH int, renditions []Rendition) []Rendition {
	selected := []Rendition{}
	for _, r := range renditions {
		if r.Width <= srcW && r.Height <= srcH {
			selected = append(selected, r)
		}
	}
	if len(selected) == 0 {
		selected = append(selected, renditions[len(renditions)-1])
	}
	return selected
}

// buildFilterComplex generates the filter_complex argument for ffmpeg
// that splits the input video into multiple streams and scales them
// according to the selected renditions.
func buildFilterComplex(selected []Rendition) string {
	splitOutputs := []string{}
	filterScales := []string{}
	for i, r := range selected {
		splitOutputs = append(splitOutputs, fmt.Sprintf("[v%d]", i))
		filterScales = append(filterScales, fmt.Sprintf("[v%d]scale=w=%d:h=%d:force_original_aspect_ratio=decrease[v%dout]", i, r.Width, r.Height, i))
	}
	return fmt.Sprintf("[0:v]split=%d%s;%s",
		len(selected),
		strings.Join(splitOutputs, ""),
		strings.Join(filterScales, ";"),
	)
}

// buildVarStreamMap constructs the var_stream_map parameter used by ffmpeg to map
// video and audio streams for each rendition in adaptive streaming.
func buildVarStreamMap(selected []Rendition) string {
	parts := []string{}
	for i := range selected {
		parts = append(parts, fmt.Sprintf("v:%d,a:%d", i, i))
	}
	return strings.Join(parts, " ")
}

// buildFFmpegArgs assembles the complete list of ffmpeg command-line arguments
// required to transcode into multiple renditions with HLS CMAF segments and DASH manifest.
func buildFFmpegArgs(
	inputPath, filterComplex string,
	selected []Rendition,
	hlsSegmentPattern, hlsPlaylistPattern, masterPlaylist, varStreamMap, dashManifest string,
) []string {
	args := []string{
		"-y",
		"-i", inputPath,
		"-filter_complex", filterComplex,
	}

	// Map video and audio streams for each rendition.
	for i := range selected {
		args = append(args, "-map", fmt.Sprintf("[v%dout]", i), "-map", "a:0?")
	}

	// Encoding settings for each rendition.
	for i, r := range selected {
		args = append(args,
			"-c:v:"+fmt.Sprint(i), "libx264",
			"-b:v:"+fmt.Sprint(i), r.VideoBitrate,
			"-preset", "veryfast",
			"-profile:v:"+fmt.Sprint(i), "main",
			"-g", "48", "-keyint_min", "48",
			"-sc_threshold", "0",
			"-c:a:"+fmt.Sprint(i), "aac",
			"-b:a:"+fmt.Sprint(i), r.AudioBitrate,
		)
	}

	// HLS and DASH output options.
	args = append(args,
		"-f", "hls",
		"-hls_time", "4",
		"-hls_playlist_type", "vod",
		"-hls_segment_type", "fmp4",
		"-hls_segment_filename", hlsSegmentPattern,
		"-master_pl_name", masterPlaylist,
		"-var_stream_map", varStreamMap,
		hlsPlaylistPattern,
		"-f", "dash",
		"-use_template", "1",
		"-use_timeline", "1",
		"-adaptation_sets", "id=0,streams=v id=1,streams=a",
		dashManifest,
	)

	return args
}
