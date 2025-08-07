package minio

import (
	"context"
	"io"
	"time"
)

type StorageAdapter interface {
	UploadDir(ctx context.Context, srcDir, targetDir string) error
	PutObject(ctx context.Context, objectName string, reader io.Reader, size int64) error
	GetObject(ctx context.Context, objectName string) (io.ReadCloser, error)
	PresignPutObject(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	PresignGetObject(ctx context.Context, objectName string, expiry time.Duration) (string, error)
}
