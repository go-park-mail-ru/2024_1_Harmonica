package repository

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const BUCKETNAME = "images"

func (r *MinioRepository) UploadImage(ctx context.Context, file []byte, filename string) (string, error) {
	objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(filename))
	fileSize := int64(len(file))
	obj := bytes.NewBuffer(file)
	var buf bytes.Buffer
	tee := io.TeeReader(obj, &buf)
	img, _, err := image.Decode(tee)
	if err != nil {
		return "", err
	}
	width := strconv.Itoa(img.Bounds().Dx())
	height := strconv.Itoa(img.Bounds().Dy())
	contentType := fmt.Sprintf("%s/%s", "image", filepath.Ext(filename))

	start := time.Now()
	info, err := r.s3.PutObject(ctx, BUCKETNAME, objectName, &buf, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"Width":  width,
			"Height": height,
		},
	})
	r.LogMinioQuery(ctx, fmt.Sprintf("PutObject(%s, %d)", objectName, fileSize), time.Since(start))
	return info.Key, err
}

func (r *MinioRepository) GetImage(ctx context.Context, name string) (*minio.Object, error) {
	start := time.Now()
	res, err := r.s3.GetObject(ctx, BUCKETNAME, name, minio.GetObjectOptions{})
	r.LogMinioQuery(ctx, fmt.Sprintf("GetObject(%s)", name), time.Since(start))
	return res, err
}

func (r *MinioRepository) GetImageBounds(ctx context.Context, name string) (int64, int64, error) {
	start := time.Now()
	img, err := r.s3.GetObject(ctx, BUCKETNAME, name, minio.GetObjectOptions{})
	r.LogMinioQuery(ctx, fmt.Sprintf("GetObject(%s)", name), time.Since(start))
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

func GetRequestId(ctx context.Context) string {
	if len(metadata.ValueFromIncomingContext(ctx, "request_id")) > 0 {
		return metadata.ValueFromIncomingContext(ctx, "request_id")[0]
	}
	return ""
}

func (r *MinioRepository) LogMinioQuery(ctx context.Context, query string, duration time.Duration) {
	requestId := GetRequestId(ctx)
	r.logger.Info("Minio query handled",
		zap.String("request_id", requestId),
		zap.String("query", query),
		zap.String("duration", duration.String()))
}
