package minio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path"
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
	bucket string
}

// New tạo MinIO client mới, đảm bảo bucket tồn tại
func New(cfg *config.Config, bucketName string) (StorageAdapter, error) {
	client, err := minio.New(cfg.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
		Region: cfg.S3.Region,
		Secure: cfg.S3.Secure,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	if err := ensureBucket(client, bucketName); err != nil {
		return nil, err
	}

	return &impl{
		cfg:    cfg,
		client: client,
		bucket: bucketName,
	}, nil
}

func ensureBucket(client *minio.Client, bucket string) error {
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket exists: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return nil
}

type fileJob struct {
	localPath  string
	objectName string
}

func getFilesFromDir(srcDir, targetDir string) ([]fileJob, error) {
	var files []fileJob

	err := filepath.Walk(srcDir, func(fullPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(srcDir, fullPath)
		if err != nil {
			return err
		}
		objectName := filepath.ToSlash(relPath)
		if targetDir != "" {
			objectName = filepath.ToSlash(path.Join(targetDir, relPath))
		}
		files = append(files, fileJob{
			localPath:  fullPath,
			objectName: objectName,
		})
		return nil
	})

	return files, err
}

func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".m3u8":
		return "application/vnd.apple.mpegurl"
	case ".ts":
		return "video/MP2T"
	default:
		if ct := mime.TypeByExtension(ext); ct != "" {
			return ct
		}
		return "application/octet-stream"
	}
}

func (i *impl) uploadFilesConcurrently(ctx context.Context, files []fileJob, concurrency int) error {
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

			contentType := getContentType(f.localPath)

			_, err = i.client.PutObject(ctx, i.bucket, f.objectName, file, stat.Size(), minio.PutObjectOptions{
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

// UploadDir upload toàn bộ folder srcDir vào bucket, targetDir là prefix trên bucket
func (i *impl) UploadDir(ctx context.Context, srcDir, targetDir string) (string, error) {
	files, err := getFilesFromDir(srcDir, targetDir)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", i.bucket, targetDir), i.uploadFilesConcurrently(ctx, files, 5)
}

func (i *impl) PutObject(ctx context.Context, objectName string, reader io.Reader, size int64) (string, error) {
	contentType := utils.DetectContentTypeByFileName(objectName)

	_, err := i.client.PutObject(ctx, i.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload %s: %w", objectName, err)
	}
	return objectName, nil
}

func (i *impl) GetObject(ctx context.Context, objectName string) ([]byte, error) {
	obj, err := i.client.GetObject(ctx, i.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s: %w", objectName, err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to read object %s: %w", objectName, err)
	}

	return data, nil
}

func (i *impl) PresignPutObject(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := i.client.PresignedPutObject(ctx, i.bucket, objectName, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to presign PUT URL for %s: %w", objectName, err)
	}
	return url.String(), nil
}

func (i *impl) PresignGetObject(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := i.client.PresignedGetObject(ctx, i.bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to presign GET URL for %s: %w", objectName, err)
	}
	return url.String(), nil
}
