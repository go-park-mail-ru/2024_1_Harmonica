package server

import (
	"context"
	"fmt"
	"harmonica/internal/entity/errs"
	image "harmonica/internal/microservices/image/proto"
	"harmonica/internal/microservices/image/server/service"
	"io"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type ImageServer struct {
	service service.IService
	logger  *zap.Logger
	image.ImageServer
}

func (s ImageServer) logErr(method string, err error, requestId string) {
	s.logger.Error(fmt.Sprintf("Error occurred in '%s'", method),
		zap.String("request_id", requestId),
		zap.Error(err))
}

func (s ImageServer) logHandle(request string, requestId string) {
	s.logger.Info(fmt.Sprintf("Request handled: '%s'", request),
		zap.String("request_id", requestId),
	)
}

func GetRequestId(ctx context.Context) string {
	if len(metadata.ValueFromIncomingContext(ctx, "request_id")) > 0 {
		return metadata.ValueFromIncomingContext(ctx, "request_id")[0]
	}
	return ""
}

func NewImageServer(service *service.Service, logger *zap.Logger) ImageServer {
	return ImageServer{service: service, logger: logger}
}

func (s ImageServer) GetImage(ctx context.Context, req *image.GetImageRequest) (*image.GetImageResponse, error) {
	requestId := GetRequestId(ctx)
	s.logHandle(fmt.Sprintf("GetImage(%s)", req.Name), requestId)

	obj, err := s.service.GetImage(ctx, req.Name)
	if err != nil {
		s.logErr("GetImage", err, requestId)
		return &image.GetImageResponse{LocalError: 30}, nil
	}
	res, err := io.ReadAll(obj)
	if err != nil {
		s.logErr("GetImage", err, requestId)
		return &image.GetImageResponse{LocalError: 29}, nil
	}
	return &image.GetImageResponse{Image: res}, nil
}

func (s ImageServer) UploadImage(ctx context.Context, req *image.UploadImageRequest) (*image.UploadImageResponse, error) {
	requestId := GetRequestId(ctx)
	s.logHandle(fmt.Sprintf("UploadImage(%s)", req.Filename), requestId)

	name, err := s.service.UploadImage(ctx, req.Image, req.Filename)
	if err != nil {
		if errs.ErrorCodes[err].LocalCode != 0 {
			return &image.UploadImageResponse{LocalError: int64(errs.ErrorCodes[err].LocalCode)}, nil
		}
		s.logErr("UploadImage", err, requestId)
		return &image.UploadImageResponse{LocalError: 18}, nil
	}
	return &image.UploadImageResponse{Name: name}, nil
}

func (s ImageServer) FormUrl(ctx context.Context, req *image.FormUrlRequest) (*image.FormUrlResponse, error) {
	requestId := GetRequestId(ctx)
	s.logHandle(fmt.Sprintf("FormUrl(%s)", req.Name), requestId)

	return &image.FormUrlResponse{Url: os.Getenv("SERVER_URL") + req.Name}, nil
}

func (s ImageServer) GetImageBounds(ctx context.Context, req *image.GetImageBoundsRequest) (*image.GetImageBoundsResponse, error) {
	requestId := GetRequestId(ctx)
	s.logHandle(fmt.Sprintf("FormUrl(%s)", req.Url), requestId)

	dx, dy, err := s.service.GetImageBounds(ctx, req.Url)
	if err != nil {
		s.logErr("GetImageBounds", err, GetRequestId(ctx))
		return &image.GetImageBoundsResponse{LocalError: 29}, nil
	}
	return &image.GetImageBoundsResponse{Dx: dx, Dy: dy}, nil
}

func NewImageServerForTests(service service.IService, logger *zap.Logger) ImageServer {
	return ImageServer{service: service, logger: logger}
}
