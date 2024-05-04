package repository

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type IRepository interface {
	UploadImage(ctx context.Context, image []byte, filename string) (string, error)
	GetImageBounds(ctx context.Context, url string) (int64, int64, error)
	GetImage(ctx context.Context, name string) (*minio.Object, error)
}
