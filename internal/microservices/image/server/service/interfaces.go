package service

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type IService interface {
	UploadImage(ctx context.Context, file []byte, fileHeader string) (string, error)
	GetImageBounds(ctx context.Context, url string) (int64, int64, error)
	GetImage(ctx context.Context, name string) (*minio.Object, error)
}
