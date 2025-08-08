package types

type TranscodeJob struct {
	MediaID string `json:"media_id"`
}

type TranscodeJobStatus string

const (
	TranscodeJobStatusPending    TranscodeJobStatus = "pending"
	TranscodeJobStatusProcessing TranscodeJobStatus = "processing"
	TranscodeJobStatusDone       TranscodeJobStatus = "done"
	TranscodeJobStatusError      TranscodeJobStatus = "error"
)

func (t TranscodeJobStatus) String() string {
	return string(t)
}
