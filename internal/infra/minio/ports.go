package minio

import (
	"context"
	"mime/multipart"
	"time"
)

type ObjectStorage interface {
	EnsureBucket(ctx context.Context) error

	Upload(ctx context.Context, objectName string, file multipart.File, size int64, contentType string) (etag string, err error)
	Delete(ctx context.Context, objectName string) error

	PresignGet(ctx context.Context, objectName string, expiry time.Duration) (url string, err error)
	PresignPut(ctx context.Context, objectName string, expiry time.Duration) (url string, err error)
}
