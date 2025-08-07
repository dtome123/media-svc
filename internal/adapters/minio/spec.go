package minio

import (
	"context"
	"io"
	"time"
)

type StorageAdapter interface {
	UploadDir(ctx context.Context, srcDir, targetDir string) error
	UploadStreamDir(ctx context.Context, srcDir, targetDir string) (string, error)
	PutObject(ctx context.Context, objectName string, reader io.Reader, size int64) (string, error)
	GetObject(ctx context.Context, objectName string) ([]byte, error)
	PresignPutObject(ctx context.Context, objectName string, expiry time.Duration) (string, error)
	PresignGetObject(ctx context.Context, objectName string, expiry time.Duration) (string, error)
}
