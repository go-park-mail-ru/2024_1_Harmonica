package service

import (
	"context"
	"harmonica/internal/entity/errs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/minio/minio-go/v7"
)

var (
	MAX_IMG_SIZE = 10000000 // 10 МБ
	MIN_IMG_SIZE = 0
	ALLOWED_EXTS = []string{".png", ".jpg", ".jpeg"}
)

func (s *RepositoryService) GetImage(ctx context.Context, name string) (*minio.Object, error) {
	return s.repo.GetImage(ctx, name)
}

func (s *RepositoryService) UploadImage(ctx context.Context, image []byte, filename string) (string, error) {
	if len(image) < MIN_IMG_SIZE || len(image) > MAX_IMG_SIZE {
		return "", errs.ErrInvalidImg
	}
	if slices.Index(ALLOWED_EXTS, strings.ToLower(filepath.Ext(filename))) == -1 {
		return "", errs.ErrNotAllowedExtension
	}
	return s.repo.UploadImage(ctx, image, filename)
}

func (s *RepositoryService) GetImageBounds(ctx context.Context, url string) (int64, int64, error) {
	if len(url) <= 0 {
		return 0, 0, nil
	}
	splitUrl := strings.Split(url, "/")
	name := splitUrl[len(splitUrl)-1]
	return s.repo.GetImageBounds(ctx, name)
}
