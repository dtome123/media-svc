package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"media-svc/config"
	"media-svc/internal/utils"
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

type fileJob struct {
	localPath  string
	objectName string
}

// UploadDir uploads all files under the given directory to the S3 bucket.
// It preserves the directory structure and uploads files concurrently.
// The S3 object keys will be prefixed with the given objectPrefix.
// UploadDir uploads all files under the given directory to the S3 bucket.
// It preserves the directory structure and uses the relative paths as object keys (no prefix).
func (i *impl) UploadDir(ctx context.Context, srcDir, targetDir string) error {
	var files []fileJob

	// Walk qua tất cả file trong thư mục
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		objectName := filepath.ToSlash(relPath)
		if targetDir != "" {
			objectName = filepath.ToSlash(filepath.Join(targetDir, relPath))
		}

		files = append(files, fileJob{
			localPath:  path,
			objectName: objectName,
		})
		return nil
	})
	if err != nil {
		return err
	}

	const concurrency = 5
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var errsMu sync.Mutex
	var errs []error

	for _, f := range files {
		wg.Add(1)
		sem <- struct{}{}

		go func(f fileJob) {
			defer wg.Done()
			defer func() { <-sem }()

			file, err := os.Open(f.localPath)
			if err != nil {
				errsMu.Lock()
				errs = append(errs, fmt.Errorf("failed to open %s: %w", f.localPath, err))
				errsMu.Unlock()
				return
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				errsMu.Lock()
				errs = append(errs, fmt.Errorf("failed to stat %s: %w", f.localPath, err))
				errsMu.Unlock()
				return
			}

			ext := filepath.Ext(f.localPath)
			var contentType string
			switch ext {
			case ".m3u8":
				contentType = "application/vnd.apple.mpegurl"
			case ".ts":
				contentType = "video/MP2T"
			default:
				contentType = mime.TypeByExtension(ext)
				if contentType == "" {
					contentType = "application/octet-stream"
				}
			}

			_, err = i.client.PutObject(ctx, i.cfg.S3.Bucket, f.objectName, file, stat.Size(), minio.PutObjectOptions{
				ContentType: contentType,
			})
			if err != nil {
				errsMu.Lock()
				errs = append(errs, fmt.Errorf("failed to upload %s: %w", f.objectName, err))
				errsMu.Unlock()
				return
			}
		}(f)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (i *impl) PutObject(ctx context.Context, objectName string, reader io.Reader, size int64) error {

	contentType := utils.DetectContentTypeByFileName(objectName)

	_, err := i.client.PutObject(ctx, i.cfg.S3.Bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload %s: %w", objectName, err)
	}
	return nil
}

func (i *impl) GetObject(ctx context.Context, objectName string) (io.ReadCloser, error) {
	obj, err := i.client.GetObject(ctx, i.cfg.S3.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s: %w", objectName, err)
	}
	return obj, nil
}

func (i *impl) PresignPutObject(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := i.client.PresignedPutObject(ctx, i.cfg.S3.Bucket, objectName, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to presign PUT URL for %s: %w", objectName, err)
	}
	return url.String(), nil
}

func (i *impl) PresignGetObject(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := i.client.PresignedGetObject(ctx, i.cfg.S3.Bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to presign GET URL for %s: %w", objectName, err)
	}
	return url.String(), nil
}
