package transcoder

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// videoStream is used to parse JSON output from ffprobe for video streams.
type videoStream struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// ffprobeOutput holds the ffprobe JSON output structure for streams.
type ffprobeOutput struct {
	Streams []videoStream `json:"streams"`
}

// getVideoResolution runs ffprobe on the input video file and extracts
// the width and height of the first video stream.
//
// Returns width, height or an error if probing fails or no video stream found.
func getVideoResolution(inputPath string) (int, int, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_streams", "-select_streams", "v", inputPath)
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe error: %w", err)
	}

	var probeData ffprobeOutput
	if err := json.Unmarshal(out, &probeData); err != nil {
		return 0, 0, fmt.Errorf("ffprobe json unmarshal error: %w", err)
	}

	if len(probeData.Streams) == 0 {
		return 0, 0, fmt.Errorf("no video stream found")
	}

	return probeData.Streams[0].Width, probeData.Streams[0].Height, nil
}
