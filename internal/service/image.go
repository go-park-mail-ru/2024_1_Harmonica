package service

import (
	"context"
	"harmonica/internal/entity/errs"
	"mime/multipart"
	"strings"

	"github.com/minio/minio-go/v7"
)

const MAX_IMG_SIZE = 10000000 // 10 МБ

func (s *RepositoryService) GetImage(ctx context.Context, name string) (*minio.Object, error) {
	return s.repo.GetImage(ctx, name)
}

func (s *RepositoryService) UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	contentType := fileHeader.Header.Get("Content-Type")
	if len(contentType) == 0 || strings.Split(contentType, "/")[0] != "image" {
		return "", errs.ErrInvalidContentType
	}
	if fileHeader.Size == 0 || fileHeader.Size > MAX_IMG_SIZE {
		return "", errs.ErrInvalidImg
	}
	return s.repo.UploadImage(ctx, file, fileHeader)
}
