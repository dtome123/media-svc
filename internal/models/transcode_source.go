package models

type TranscodeSource struct {
	VideoID  string `bson:"video_id" json:"video_id"`
	FilePath string `bson:"file_path" json:"file_path"`
}
