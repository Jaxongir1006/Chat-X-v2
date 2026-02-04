package minio

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/Jaxongir1006/Chat-X-v2/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client *minio.Client
	bucket string
}

func New(cfg config.MinioConfig) (*Storage, error) {
	if cfg.Endpoint == "" || cfg.User == "" || cfg.Password == "" {
		return nil, fmt.Errorf("minio config is incomplete (endpoint/user/password required)")
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("minio bucket is required")
	}

	mc, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio init client: %w", err)
	}

	return &Storage{client: mc, bucket: cfg.Bucket}, nil
}

func (s *Storage) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("minio bucket exists: %w", err)
	}
	if exists {
		return nil
	}

	if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("minio make bucket: %w", err)
	}
	return nil
}

func (s *Storage) Upload(ctx context.Context, objectName string, file multipart.File, size int64, contentType string) (string, error) {
	if objectName == "" {
		return "", fmt.Errorf("objectName is required")
	}
	if file == nil {
		return "", fmt.Errorf("file is nil")
	}
	if size <= 0 {
		return "", fmt.Errorf("size must be > 0")
	}

	info, err := s.client.PutObject(ctx, s.bucket, objectName, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minio put object: %w", err)
	}
	return info.ETag, nil
}

func (s *Storage) Delete(ctx context.Context, objectName string) error {
	if objectName == "" {
		return fmt.Errorf("objectName is required")
	}
	if err := s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("minio remove object: %w", err)
	}
	return nil
}

func (s *Storage) PresignGet(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	if objectName == "" {
		return "", fmt.Errorf("objectName is required")
	}
	if expiry <= 0 {
		expiry = time.Hour
	}

	u, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("minio presign get: %w", err)
	}
	return u.String(), nil
}

func (s *Storage) PresignPut(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	if objectName == "" {
		return "", fmt.Errorf("objectName is required")
	}
	if expiry <= 0 {
		expiry = time.Hour
	}

	u, err := s.client.PresignedPutObject(ctx, s.bucket, objectName, expiry)
	if err != nil {
		return "", fmt.Errorf("minio presign put: %w", err)
	}
	return u.String(), nil
}
