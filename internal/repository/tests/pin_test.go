package test_repository

import (
	"database/sql/driver"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/microservices/image/proto"
	"harmonica/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	Limit  = 10
	Offset = 10
)

func TestRepository_GetFeedPins(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	tests := []struct {
		name         string
		setupMocks   func()
		expectedPins entity.FeedPins
		expectedErr  error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetPinsFeed).
					WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, nil).MaxTimes(2)
			},
			expectedPins: entity.FeedPins{Pins: []entity.FeedPinResponse{{PinId: 1}}},
			expectedErr:  nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetPinsFeed).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedPins: entity.FeedPins{},
			expectedErr:  errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetPinsFeed).
					WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal)
			},
			expectedPins: entity.FeedPins{},
			expectedErr:  errs.ErrDBInternal,
		},
		{
			name: "Error test 3",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetPinsFeed).
					WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, nil)
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal)
			},
			expectedPins: entity.FeedPins{},
			expectedErr:  errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetFeedPins(CtxWithRequestId, Limit, Offset)
			assert.Equal(t, tc.expectedPins, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetSubscriptionsFeedPins(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	userId := entity.UserID(1)

	tests := []struct {
		name         string
		setupMocks   func()
		expectedPins entity.FeedPins
		expectedErr  error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetSubscriptionsFeedPins).
					WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, nil)
			},
			expectedPins: entity.FeedPins{Pins: []entity.FeedPinResponse{{PinId: 1}}},
			expectedErr:  nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetSubscriptionsFeedPins).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedPins: entity.FeedPins{},
			expectedErr:  errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetSubscriptionsFeedPins).
					WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow(1))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{}, errs.ErrDBInternal)
			},
			expectedPins: entity.FeedPins{},
			expectedErr:  errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetSubscriptionsFeedPins(CtxWithRequestId, userId, Limit, Offset)
			assert.Equal(t, tc.expectedPins, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetUserPins(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	authorId := entity.UserID(1)

	tests := []struct {
		name         string
		setupMocks   func()
		expectedPins entity.UserPins
		expectedErr  error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserPins).
					WithArgs(authorId, Limit, Offset).
					WillReturnRows(sqlmock.NewRows([]string{"content_url"}).AddRow("url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{Dx: 1, Dy: 2}, nil).MaxTimes(1)
			},
			expectedPins: entity.UserPins{
				Pins: []entity.UserPinResponse{
					{ContentUrl: "url", ContentDX: 1, ContentDY: 2},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserPins).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedPins: entity.UserPins{},
			expectedErr:  errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserPins).
					WithArgs(authorId, Limit, Offset).
					WillReturnRows(sqlmock.NewRows([]string{"content_url"}).AddRow("url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(nil, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedPins: entity.UserPins{},
			expectedErr:  errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetUserPins(CtxWithRequestId, authorId, Limit, Offset)
			assert.Equal(t, tc.expectedPins, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetPinById(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	pinId := entity.PinID(1)

	tests := []struct {
		name        string
		setupMocks  func()
		expectedPin entity.PinPageResponse
		expectedErr error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetPinById).
					WillReturnRows(sqlmock.NewRows([]string{"content_url"}).AddRow("url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{Dx: 1, Dy: 2}, nil)
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(&proto.GetImageBoundsResponse{Dx: 1, Dy: 2}, nil)
			},
			expectedPin: entity.PinPageResponse{
				ContentUrl: "url",
				ContentDX:  1,
				ContentDY:  2,
				PinAuthor: entity.PinAuthor{
					AvatarDX: 1,
					AvatarDY: 2,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetPinById).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedPin: entity.PinPageResponse{},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetPinById).
					WillReturnRows(sqlmock.NewRows([]string{"content_url"}).AddRow("url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(nil, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedPin: entity.PinPageResponse{},
			expectedErr: errs.ErrDBInternal,
		},
		{
			name: "Error test 3",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetPinById).
					WillReturnRows(sqlmock.NewRows([]string{"content_url"}).AddRow("url"))
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: "url"}).
					Return(&proto.GetImageBoundsResponse{}, nil).MaxTimes(1)
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{}).
					Return(nil, errs.ErrDBInternal).MaxTimes(1)
			},
			expectedPin: entity.PinPageResponse{},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.GetPinById(CtxWithRequestId, pinId)
			assert.Equal(t, tc.expectedPin, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_CreatePin(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	pin := entity.Pin{Title: "Title"}
	expectedID := entity.PinID(10)

	tests := []struct {
		name        string
		setupMocks  func()
		expectedID  entity.PinID
		expectedErr error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryCreatePin).
					WillReturnRows(sqlmock.NewRows([]string{"pin_id"}).AddRow("10"))
			},
			expectedID:  expectedID,
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryCreatePin).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedID:  entity.PinID(0),
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.CreatePin(CtxWithRequestId, pin)
			assert.Equal(t, tc.expectedID, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_UpdatePin(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	pin := entity.Pin{Title: "Title"}

	tests := []struct {
		name        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdatePin).WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryUpdatePin).WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.UpdatePin(CtxWithRequestId, pin)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_DeletePin(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	pinId := entity.PinID(1)

	tests := []struct {
		name        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryDeletePin).WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryDeletePin).WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.DeletePin(CtxWithRequestId, pinId)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_CheckPinExistence(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	pinId := entity.PinID(1)

	tests := []struct {
		name        string
		setupMocks  func()
		expectedRes bool
		expectedErr error
	}{
		{
			name: "OK test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryCheckPinExistence).
					WithArgs(entity.PinID(1)).
					WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("true"))
			},
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryCheckPinExistence).
					WithArgs(entity.PinID(1)).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedRes: false,
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			res, err := repo.CheckPinExistence(CtxWithRequestId, pinId)
			assert.Equal(t, tc.expectedRes, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
