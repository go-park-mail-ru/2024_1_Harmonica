package handler

import (
	"context"
	"harmonica/config"
	"harmonica/internal/entity/errs"
	image "harmonica/internal/microservices/image/proto"
	"io"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

func (h *APIHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	requestId := r.Context().Value("request_id").(string)
	name, err := ReadStringSlug(r, "image_name")
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrInvalidSlug})
		return
	}
	ctx := metadata.NewOutgoingContext(r.Context(), metadata.Pairs("request_id", requestId))
	res, err := h.ImageService.GetImage(ctx, &image.GetImageRequest{Name: name})
	if err != nil {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{GeneralErr: err, LocalErr: errs.ErrGRPCWentWrong})
		return
	}
	if res.LocalError != 0 {
		WriteErrorResponse(w, h.logger, requestId, errs.ErrorInfo{LocalErr: errs.GetLocalErrorByCode[res.LocalError]})
		return
	}
	w.Write(res.Image)
}

func (h *APIHandler) UploadImage(r *http.Request, imageName string) (string, error) {
	file, header, err := r.FormFile(imageName)
	if err != nil {
		return "", errs.ErrNoImageProvided
	}
	contentType := header.Header.Get("Content-Type")
	if len(contentType) == 0 || strings.Split(contentType, "/")[0] != "image" {
		return "", errs.ErrInvalidContentType
	}
	f, err := io.ReadAll(file)
	if err != nil {
		return "", errs.ErrInvalidImg
	}
	res, err := h.ImageService.UploadImage(r.Context(), &image.UploadImageRequest{Image: f, Filename: header.Filename})
	if err != nil {
		return "", errs.ErrGRPCWentWrong
	}
	if res.LocalError != 0 {
		return "", errs.GetLocalErrorByCode[res.LocalError]
	}
	return res.Name, nil
}

func (h *APIHandler) FormImgURL(name string) string {
	res, err := h.ImageService.FormUrl(context.Background(), &image.FormUrlRequest{Name: name})
	if err != nil {
		return config.GetEnv("SERVER_URL", "") + name
	}
	return res.Url
}
