package handler

import (
	"harmonica/internal/entity/errs"
	"net/http"
	"strconv"
)

func ReadInt64Slug(r *http.Request, name string) (uint64, error) {
	stringId := r.PathValue(name)
	if len(stringId) == 0 {
		return 0, errs.ErrInvalidSlug
	}
	res, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		return 0, errs.ErrInvalidSlug
	}
	return res, nil
}
