package minio

import (
	"context"
	"fmt"
	"media-svc/config"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type impl struct {
	cfg    *config.Config
	client *minio.Client
}

func New(cfg *config.Config) StorageAdapter {

	client, err := minio.New(cfg.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
		Region: cfg.S3.Region,
		Secure: cfg.S3.Insecure,
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.S3.Bucket)
	if err != nil {
		panic(err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.S3.Bucket, minio.MakeBucketOptions{}); err != nil {
			panic(err)
		}
	}

	return &impl{
		cfg:    cfg,
		client: client,
	}
}

func (i *impl) UploadDir(outputDir, objectPrefix string) error {
	ctx := context.Background()

	return filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		objectName := objectPrefix + info.Name()
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		stat, _ := file.Stat()

		_, err = i.client.PutObject(ctx, i.cfg.S3.Bucket, objectName, file, stat.Size(), minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		})
		if err != nil {
			return fmt.Errorf("failed uploading %s: %w", objectName, err)
		}

		return nil
	})
}
