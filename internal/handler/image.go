package handler

import (
	"harmonica/config"
	"harmonica/internal/entity/errs"
	"io"
	"net/http"
)

func (h *APIHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := ctx.Value("request_id").(string)

	name, err := ReadStringSlug(r, "image_name")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrInvalidSlug})
		return
	}
	res, err := h.service.GetImage(r.Context(), name)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrInvalidImg})
		return
	}
	file, err := io.ReadAll(res)
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrInvalidImg})
		return
	}
	w.Write(file)
}

func (h *APIHandler) UploadImage(r *http.Request, imageName string) (string, error) {
	file, header, err := r.FormFile(imageName)
	if err != nil {
		return "", errs.ErrNoImageProvided
	}
	name, err := h.service.UploadImage(r.Context(), file, header)
	if err != nil {
		return "", err
	}
	return name, nil
}

func FormImgURL(name string) string {
	return config.GetEnv("SERVER_URL", "") + name
}
