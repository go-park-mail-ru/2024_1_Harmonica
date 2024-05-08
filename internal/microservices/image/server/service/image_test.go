package service

import (
	"context"
	"github.com/minio/minio-go/v7"
	"harmonica/internal/entity/errs"
	mock_repository "harmonica/mocks/microservices/image/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryService_GetImage(t *testing.T) {
	type ExpectedReturn struct {
		Object *minio.Object
		Error  error
	}
	type ExpectedMockReturn struct {
		Object *minio.Object
		Error  error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		mockArgs string, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetImage(ctx, mockArgs).Return(
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
				Object: &minio.Object{},
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
				Error: errs.ErrInvalidImg,
			},
			expectedMockArgs: "test_image.jpg",
			expectedMockReturn: ExpectedMockReturn{
				Error: errs.ErrInvalidImg,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			mockBehaviour(repo, context.Background(), testCase.expectedMockArgs, testCase.expectedMockReturn)
			s := NewRepositoryService(repo)
			object, err := s.GetImage(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.Object, object)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}

func TestRepositoryService_UploadImage(t *testing.T) {
	type ExpectedReturn struct {
		Filename string
		Error    error
	}
	type ExpectedMockReturn struct {
		Error error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		image []byte, filename string, mockReturn ExpectedMockReturn) {
		repo.EXPECT().UploadImage(ctx, image, filename).Return(
			filename, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               []byte
		filename           string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name:     "OK test case 1",
			args:     []byte{0x89, 0x50, 0x4e, 0x47},
			filename: "test_image.jpg",
			expectedReturn: ExpectedReturn{
				Filename: "test_image.jpg",
			},
			expectedMockReturn: ExpectedMockReturn{},
		},
		{
			name:     "Error test case 1",
			args:     []byte{},
			filename: "test_image.jpg",
			expectedReturn: ExpectedReturn{
				Filename: "test_image.jpg",
				Error:    errs.ErrInvalidImg,
			},
			expectedMockReturn: ExpectedMockReturn{
				Error: errs.ErrInvalidImg,
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockIRepository(ctrl)
			mockBehaviour(repo, context.Background(), testCase.args, testCase.filename, testCase.expectedMockReturn)
			s := NewRepositoryService(repo)
			filename, err := s.UploadImage(context.Background(), testCase.args, testCase.filename)
			assert.Equal(t, testCase.expectedReturn.Filename, filename)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}

func TestRepositoryService_GetImageBounds(t *testing.T) {
	type ExpectedReturn struct {
		Width  int64
		Height int64
		Error  error
	}
	type ExpectedMockReturn struct {
		Width  int64
		Height int64
		Error  error
	}
	mockBehaviour := func(repo *mock_repository.MockIRepository, ctx context.Context,
		url string, mockReturn ExpectedMockReturn) {
		repo.EXPECT().GetImageBounds(ctx, url).Return(
			mockReturn.Width, mockReturn.Height, mockReturn.Error)
	}
	testTable := []struct {
		name               string
		args               string
		expectedReturn     ExpectedReturn
		expectedMockReturn ExpectedMockReturn
	}{
		{
			name: "OK test case 1",
			args: "test_image.jpg",
			expectedReturn: ExpectedReturn{
				Width:  800,
				Height: 600,
			},
			expectedMockReturn: ExpectedMockReturn{
				Width:  800,
				Height: 600,
			},
		},
		{
			name: "Error test case 1",
			args: "test_image.jpg",
			expectedReturn: ExpectedReturn{
				Error: errs.ErrDBInternal,
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
			repo := mock_repository.NewMockIRepository(ctrl)
			mockBehaviour(repo, context.Background(), testCase.args, testCase.expectedMockReturn)
			s := NewRepositoryService(repo)
			width, height, err := s.GetImageBounds(context.Background(), testCase.args)
			assert.Equal(t, testCase.expectedReturn.Width, width)
			assert.Equal(t, testCase.expectedReturn.Height, height)
			assert.Equal(t, testCase.expectedReturn.Error, err)
		})
	}
}
