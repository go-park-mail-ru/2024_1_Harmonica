package handler

import (
	"net/http"
	"strconv"
)

func ReadUint64Slug(r *http.Request, name string) (uint64, error) {
	value := r.PathValue(name)
	res, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func ReadStringSlug(r *http.Request, name string) (string, error) {
	res := r.PathValue(name)
	if len(res) == 0 {
		return "", errs.ErrInvalidSlug
	}
	return res, nil
}
