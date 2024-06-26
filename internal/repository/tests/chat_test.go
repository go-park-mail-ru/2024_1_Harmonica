package test_repository

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"harmonica/internal/entity"
	"harmonica/internal/entity/errs"
	"harmonica/internal/microservices/image/proto"
	"harmonica/internal/repository"
	"testing"
	"time"
)

func TestRepository_CreateMessage(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	message := entity.Message{
		SenderId:   entity.UserID(1),
		ReceiverId: entity.UserID(2),
		Text:       "Hello!",
	}

	tests := []struct {
		name        string
		setupMocks  func()
		expectedErr error
	}{
		{
			name: "OK test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryCreateMessage).
					WithArgs(message.SenderId, message.ReceiverId, message.Text).
					WillReturnResult(driver.ResultNoRows)
			},
			expectedErr: nil,
		},
		{
			name: "Error test",
			setupMocks: func() {
				mock.ExpectExec(repository.QueryCreateMessage).
					WithArgs(message.SenderId, message.ReceiverId, message.Text).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedErr: errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			err := repo.CreateMessage(CtxWithRequestId, message)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetMessages(t *testing.T) {
	db, mock, ctrl, imageClient, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()

	firstUserID := entity.UserID(1)
	secondUserID := entity.UserID(2)
	dialogUser := entity.User{
		UserID:    firstUserID,
		Nickname:  "testuser",
		AvatarURL: "http://example.com/avatar.png",
	}

	tests := []struct {
		name           string
		setupMocks     func()
		expectedResult entity.Messages
		expectedErr    error
	}{
		{
			name: "OK test",
			setupMocks: func() {
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: dialogUser.AvatarURL}).
					Return(&proto.GetImageBoundsResponse{}, nil)
				mock.ExpectQuery(repository.QueryGetUserById).
					WithArgs(firstUserID).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "nickname", "avatar_url"}).
						AddRow(dialogUser.UserID, dialogUser.Nickname, dialogUser.AvatarURL))
				mock.ExpectQuery(repository.QueryGetMessages).
					WithArgs(firstUserID, secondUserID).
					WillReturnRows(sqlmock.NewRows([]string{"sender_id", "receiver_id", "text"}).
						AddRow(1, 2, "Hello!"))
			},
			expectedResult: entity.Messages{
				User: entity.UserFromChat{
					UserID:    dialogUser.UserID,
					Nickname:  dialogUser.Nickname,
					AvatarURL: dialogUser.AvatarURL,
				},
				Messages: []entity.MessageResponse{
					{
						SenderId: 1,
						Text:     "Hello!",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error test 1",
			setupMocks: func() {
				imageClient.EXPECT().GetImageBounds(CtxWithRequestId, &proto.GetImageBoundsRequest{Url: dialogUser.AvatarURL}).
					Return(&proto.GetImageBoundsResponse{}, nil)
				mock.ExpectQuery(repository.QueryGetUserById).
					WithArgs(firstUserID).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "nickname", "avatar_url"}).
						AddRow(dialogUser.UserID, dialogUser.Nickname, dialogUser.AvatarURL))
				mock.ExpectQuery(repository.QueryGetMessages).
					WithArgs(firstUserID, secondUserID).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedResult: entity.Messages{},
			expectedErr:    errs.ErrDBInternal,
		},
		{
			name: "Error test 2",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserById).
					WithArgs(firstUserID).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedResult: entity.Messages{},
			expectedErr:    errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			result, err := repo.GetMessages(CtxWithRequestId, firstUserID, secondUserID)
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestRepository_GetUserChats(t *testing.T) {
	db, mock, ctrl, _, repo := SetupDBMock(t)
	defer ctrl.Finish()
	defer db.Close()
	testTime := time.Now()

	tests := []struct {
		name           string
		setupMocks     func()
		expectedResult entity.UserChats
		expectedErr    error
	}{
		{
			name: "OK test",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserChats).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"user_user_id", "user_nickname", "user_avatar_url",
						"chat_last_message_sender_id", "chat_last_message_text", "chat_last_message_message_read",
						"chat_last_message_sent_at"}).
						AddRow(1, "user", "", 0, "", false, testTime))
			},
			expectedResult: entity.UserChats{
				Chats: []entity.UserChat{
					{
						User:        entity.UserFromChat{UserID: 1, Nickname: "user"},
						LastMessage: entity.MessageResponse{SentAt: testTime},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Error test",
			setupMocks: func() {
				mock.ExpectQuery(repository.QueryGetUserChats).
					WillReturnError(errs.ErrDBInternal)
			},
			expectedResult: entity.UserChats{},
			expectedErr:    errs.ErrDBInternal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			result, err := repo.GetUserChats(CtxWithRequestId, entity.UserID(1))
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
