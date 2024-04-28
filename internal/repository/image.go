package repository

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
)

const BUCKETNAME = "images"

func (r *DBRepository) UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
	fileSize := fileHeader.Size
	var buf bytes.Buffer
	tee := io.TeeReader(file, &buf)
	img, _, err := image.Decode(tee)
	if err != nil {
		return "", err
	}
	width := strconv.Itoa(img.Bounds().Dx())
	height := strconv.Itoa(img.Bounds().Dy())
	contentType := fileHeader.Header.Get("Content-Type")
	info, err := r.s3.PutObject(ctx, BUCKETNAME, objectName, &buf, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"width":  width,
			"height": height,
		},
	})
	return info.Key, err
}

func (r *DBRepository) GetImage(ctx context.Context, name string) (*minio.Object, error) {
	return r.s3.GetObject(ctx, BUCKETNAME, name, minio.GetObjectOptions{})
}
