package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TranscodeStatus string

const (
	TranscodeStatusPending TranscodeStatus = "pending"
	TranscodeStatusSuccess TranscodeStatus = "success"
	TranscodeStatusFailed  TranscodeStatus = "failed"
)

type TranscodeJob struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MediaID    primitive.ObjectID `bson:"media_id" json:"media_id"`                           // Reference to the original media
	Status     string             `bson:"status" json:"status"`                               // pending, processing, success, failed
	OutputPath string             `bson:"output_path,omitempty" json:"output_path,omitempty"` // Folder or key where HLS/DASH is stored
	Quality    string             `bson:"quality,omitempty" json:"quality,omitempty"`         // e.g. "720p", "1080p"
	Error      string             `bson:"error,omitempty" json:"error,omitempty"`             // Error message if failed
	StartedAt  *time.Time         `bson:"started_at,omitempty" json:"started_at,omitempty"`   // When processing started
	FinishedAt *time.Time         `bson:"finished_at,omitempty" json:"finished_at,omitempty"` // When processing finished
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`                       // Job creation time
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`                       // Last update time
}

func (coll TranscodeJob) CollectionName() string {
	return "transcode_jobs"
}
