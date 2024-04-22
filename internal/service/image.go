package service

import (
	"context"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"mime/multipart"
	"strings"

	"github.com/minio/minio-go/v7"
)

const (
	MAX_IMG_SIZE = 10000000 // 10 МБ
	MIN_IMG_SIZE = 100000
)

func (s *RepositoryService) GetImageNameById(ctx context.Context, id entity.ImageID) (string, error) {
	return s.repo.GetImageNameById(ctx, id)
}

func (s *RepositoryService) GetImage(ctx context.Context, name string) (*minio.Object, error) {
	return s.repo.GetImage(ctx, name)
}

func (s *RepositoryService) GetImageById(ctx context.Context, id entity.ImageID) (*minio.Object, error) {
	return s.repo.GetImageById(ctx, id)
}

func (s *RepositoryService) UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (entity.ImageID, string, error) {
	contentType := fileHeader.Header.Get("Content-Type")
	if len(contentType) == 0 || strings.Split(contentType, "/")[0] != "image" {
		return entity.ImageID(0), "", errs.ErrInvalidContentType
	}
	if fileHeader.Size == 0 || fileHeader.Size > MAX_IMG_SIZE || fileHeader.Size < MIN_IMG_SIZE {
		return entity.ImageID(0), "", errs.ErrInvalidImg
	}
	return s.repo.UploadImage(ctx, file, fileHeader)
}
