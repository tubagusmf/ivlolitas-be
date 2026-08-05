package storage

import (
	"context"
	"mime/multipart"
)

type UploadResult struct {
	URL      string
	PublicID string
}

type Storage interface {
	Upload(ctx context.Context, file multipart.File, folder string) (*UploadResult, error)
	Delete(ctx context.Context, publicID string) error
}
