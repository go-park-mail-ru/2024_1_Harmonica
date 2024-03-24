package handler

import (
	"harmonica/internal/entity/errs"
	"net/http"
	"strconv"
	"strings"
)

func ReadInt64Slug(r *http.Request) (int64, error) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) <= 1 {
		return 0, errs.ErrInvalidSlug
	}
	slug := parts[len(parts)-1]
	int64Slug, err := strconv.ParseInt(slug, 10, 64)
	if err != nil {
		return 0, errs.ErrInvalidSlug
	}
	return int64Slug, nil
}
