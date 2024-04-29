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
	"strings"
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

func (r *DBRepository) GetImageBounds(ctx context.Context, url string) (int64, int64, error) {
	if len(url) <= 0 {
		return 0, 0, nil
	}
	splitUrl := strings.Split(url, "/")
	if len(splitUrl) <= 0 {
		return 0, 0, nil
	}
	name := splitUrl[len(splitUrl)-1]

	img, err := r.s3.GetObject(ctx, BUCKETNAME, name, minio.GetObjectOptions{})
	if err != nil {
		return 0, 0, err
	}
	stat, err := img.Stat()
	if err != nil {
		return 0, 0, err
	}
	width := 0
	height := 0

	widthString := stat.Metadata.Get("X-Amz-Meta-Width")
	heightString := stat.Metadata.Get("X-Amz-Meta-Height")

	if widthString != "" {
		width, err = strconv.Atoi(widthString)
		if err != nil {
			return 0, 0, err
		}
	}
	if heightString != "" {
		height, err = strconv.Atoi(heightString)
		if err != nil {
			return 0, 0, err
		}
	}
	return int64(width), int64(height), nil
}
