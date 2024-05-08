package server

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"harmonica/internal/entity/errs"
	image "harmonica/internal/microservices/image/proto"
	mock_service "harmonica/mocks/microservices/image/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestImageServer_UploadImage(t *testing.T) {
	type ExpectedReturn struct {
		Response *image.UploadImageResponse
		Error    error
	}
	type ExpectedMockReturn struct {
		Name  string
		Error error
	}
	mockBehaviour := func(service *mock_service.MockIService, ctx context.Context,
		image []byte, filename string, mockReturn ExpectedMockReturn) {
		service.EXPECT().UploadImage(ctx, image, filename).Return(
			mockReturn.Name, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               *image.UploadImageRequest
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: &image.UploadImageRequest{
				Image:    []byte{0x89, 0x50, 0x4e, 0x47},
				Filename: "test_image.jpg",
			},
			expectedReturn: ExpectedReturn{
				Response: &image.UploadImageResponse{Name: "test_image.jpg"},
			},
			expectedMockReturn: ExpectedMockReturn{
				Name: "test_image.jpg",
			},
		},
		{
			name: "Error test case 1",
			args: &image.UploadImageRequest{
				Image:    []byte{0x89, 0x50, 0x4e, 0x47},
				Filename: "test_image.jpg",
			},
			expectedReturn: ExpectedReturn{
				Response: &image.UploadImageResponse{LocalError: 18},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errors.New("some error"),
			},
		},
		{
			name: "Error test case 2",
			args: &image.UploadImageRequest{
				Image:    []byte{0x89, 0x50, 0x4e, 0x47},
				Filename: "test_image.jpg",
			},
			expectedReturn: ExpectedReturn{
				Response: &image.UploadImageResponse{LocalError: int64(errs.ErrorCodes[errs.ErrDBInternal].LocalCode)},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := mock_service.NewMockIService(ctrl)
			server := NewImageServerForTests(service, zap.NewNop())
			mockBehaviour(service, context.Background(), testCase.args.Image, testCase.args.Filename, testCase.expectedMockReturn)
			response, err := server.UploadImage(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.Response, response)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}
func TestImageServer_FormUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service := mock_service.NewMockIService(ctrl)
	server := NewImageServerForTests(service, zap.NewNop())
	resp, err := server.FormUrl(context.Background(), &image.FormUrlRequest{Name: "image.jpg"})
	assert.Equal(t, "image.jpg", resp.Url)
	assert.Equal(t, nil, err)
}

func TestImageServer_GetImageBounds(t *testing.T) {
	type ExpectedReturn struct {
		Response *image.GetImageBoundsResponse
		Error    error
	}
	type ExpectedMockReturn struct {
		Dx    int64
		Dy    int64
		Error error
	}
	mockBehaviour := func(service *mock_service.MockIService, ctx context.Context,
		url string, mockReturn ExpectedMockReturn) {
		service.EXPECT().GetImageBounds(ctx, url).Return(
			mockReturn.Dx, mockReturn.Dy, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               *image.GetImageBoundsRequest
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: &image.GetImageBoundsRequest{
				Url: "http://example.com/image.jpg",
			},
			expectedReturn: ExpectedReturn{
				Response: &image.GetImageBoundsResponse{Dx: 100, Dy: 200},
			},
			expectedMockReturn: ExpectedMockReturn{
				Dx: 100,
				Dy: 200,
			},
		},
		{
			name: "Error test case 1",
			args: &image.GetImageBoundsRequest{
				Url: "http://example.com/image.jpg",
			},
			expectedReturn: ExpectedReturn{
				Response: &image.GetImageBoundsResponse{LocalError: 29},
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := mock_service.NewMockIService(ctrl)
			server := NewImageServerForTests(service, zap.NewNop())
			mockBehaviour(service, context.Background(), testCase.args.Url, testCase.expectedMockReturn)
			response, err := server.GetImageBounds(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.Response, response)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}

/* // не получается работа с объектом минио
func TestImageServer_GetImage(t *testing.T) {
	type ExpectedReturn struct {
		Response *image.GetImageResponse
		Error    error
	}
	type ExpectedMockReturn struct {
		Object *minio.Object
		Error  error
	}
	mockBehaviour := func(service *mock_service.MockIService, ctx context.Context,
		name string, mockReturn ExpectedMockReturn) {
		service.EXPECT().GetImage(ctx, name).Return(
			mockReturn.Object, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               string
		expectedReturn     ExpectedReturn
		expectedMockArgs   string
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: "test_image.jpg",
			expectedReturn: ExpectedReturn{
				Response: &image.GetImageResponse{Image: []byte{0x89, 0x50, 0x4e, 0x47}},
			},
			expectedMockArgs: "test_image.jpg",
			expectedMockReturn: ExpectedMockReturn{
				Object: &minio.Object{},
			},
		},
		{
			name: "Error test case 1",
			args: "test_image.jpg",
			expectedReturn: ExpectedReturn{
				Error: errs.ErrDBInternal,
			},
			expectedMockArgs: "test_image.jpg",
			expectedMockReturn: ExpectedMockReturn{
				Error: errs.ErrDBInternal,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			service := mock_service.NewMockIService(ctrl)
			server := NewImageServerForTests(service, zap.NewNop())
			mockBehaviour(service, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			response, err := server.GetImage(context.Background(), &image.GetImageRequest{Name: testCase.args})
			assert.Equal(t, testCase.expectedReturn.Response, response)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}
*/
