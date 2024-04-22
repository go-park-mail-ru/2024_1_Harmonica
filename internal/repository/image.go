package repository

import (
	"context"
	"fmt"
	"harmonica/internal/entity"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

const (
	BUCKETNAME       = "images"
	InsertImageQuery = `INSERT INTO public.image ("name") VALUES ($1) RETURNING public.image.image_id, public.image.name`
	GetImageById     = `SELECT public.image.name FROM public.image WHERE public.image.image_id=$1`
)

func (r *DBRepository) GetImageNameById(ctx context.Context, id entity.ImageID) (string, error) {
	imageName := ""
	start := time.Now()
	err := r.db.QueryRowxContext(ctx, GetImageById, id).Scan(&imageName)
	LogDBQuery(r, ctx, GetImageById, time.Since(start))
	return imageName, err
}

func (r *DBRepository) GetImageById(ctx context.Context, id entity.ImageID) (*minio.Object, error) {
	imageName, err := r.GetImageNameById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.GetImage(ctx, imageName)
}

func (r *DBRepository) UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (entity.ImageID, string, error) {
	objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	fileSize := fileHeader.Size
	contentType := fileHeader.Header.Get("Content-Type")
	info, err := r.s3.PutObject(ctx, BUCKETNAME, objectName, file, fileSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return entity.ImageID(0), "", err
	}
	imageId := entity.ImageID(0)
	imageName := ""
	start := time.Now()
	err = r.db.QueryRowxContext(ctx, InsertImageQuery, info.Key).Scan(&imageId, &imageName)
	LogDBQuery(r, ctx, InsertImageQuery, time.Since(start))
	return imageId, imageName, err
}

func (r *DBRepository) GetImage(ctx context.Context, name string) (*minio.Object, error) {
	return r.s3.GetObject(ctx, BUCKETNAME, name, minio.GetObjectOptions{})
}
