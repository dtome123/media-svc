package minio

type StorageAdapter interface {
	UploadDir(outputDir, objectPrefix string) error
}
