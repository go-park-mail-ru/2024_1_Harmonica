package handler

import (
	"harmonica/internal/entity/errs"
	"io"
	"net/http"
)

func (h *APIHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	name, err := ReadStringSlug(r, "image_name")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrInvalidSlug})
		return
	}
	res, err := h.service.GetImage(r.Context(), name)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrInvalidImg})
		return
	}
	file, err := io.ReadAll(res)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrInvalidImg})
		return
	}
	w.Write(file)
}

func (h *APIHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{GeneralErr: errs.ErrInvalidInputFormat, LocalErr: errs.ErrInvalidInputFormat})
		return
	}
	name, err := h.service.UploadImage(r.Context(), file, header)
	if err != nil {
		WriteErrorResponse(w, h.logger, errs.ErrorInfo{GeneralErr: err, LocalErr: err})
		return
	}
	WriteDefaultResponse(w, h.logger, name)
}
