package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name"`                                             // Display name of the media file
	Description     string             `bson:"description,omitempty" json:"description,omitempty"`           // Optional description or caption
	Path            string             `bson:"path" json:"path"`                                             // Physical or remote path (e.g., local path or S3 key)
	ContentType     string             `bson:"content_type" json:"content_type"`                             // MIME type (e.g., video/mp4, image/png)
	Size            int64              `bson:"size" json:"size"`                                             // File size in bytes
	Duration        float64            `bson:"duration,omitempty" json:"duration,omitempty"`                 // Duration in seconds (for video or audio)
	Width           int                `bson:"width,omitempty" json:"width,omitempty"`                       // Media width in pixels (if applicable)
	Height          int                `bson:"height,omitempty" json:"height,omitempty"`                     // Media height in pixels (if applicable)
	IsStreamable    bool               `bson:"is_streamable" json:"is_streamable"`                           // Indicates whether the media supports streaming
	Tags            []string           `bson:"tags,omitempty" json:"tags,omitempty"`                         // Tags or keywords for filtering/searching
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`                                 // Timestamp when the media was created
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`                                 // Timestamp when the media was last updated
	TranscodeSource *TranscodeSource   `bson:"transcode_source,omitempty" json:"transcode_source,omitempty"` // Optional transcode source only for video
}

func (coll Media) CollectionName() string {
	return "medias"
}

func (coll *Media) BeforeCreate() {

	if coll.ID.IsZero() {
		coll.ID = primitive.NewObjectID()
	}

	coll.CreatedAt = time.Now().UTC()
	coll.UpdatedAt = time.Now().UTC()
}

func (coll *Media) BeforeUpdate() {
	coll.UpdatedAt = time.Now().UTC()
}

type TranscodeSource struct {
	FilePath string `bson:"file_path" json:"file_path"`
}
