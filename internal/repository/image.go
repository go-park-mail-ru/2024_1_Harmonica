package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

const BUCKETNAME = "images"

func (r *DBRepository) UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	fileSize := fileHeader.Size
	contentType := fileHeader.Header.Get("Content-Type")
	info, err := r.s3.PutObject(ctx, BUCKETNAME, objectName, file, fileSize, minio.PutObjectOptions{ContentType: contentType})
	return info.Key, err
}

func (r *DBRepository) GetImage(ctx context.Context, name string) (*minio.Object, error) {
	return r.s3.GetObject(ctx, BUCKETNAME, name, minio.GetObjectOptions{})
}
