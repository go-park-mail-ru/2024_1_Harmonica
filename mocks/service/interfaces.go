// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/interfaces.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	entity "harmonica/internal/entity"
	errs "harmonica/internal/entity/errs"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIService is a mock of IService interface.
type MockIService struct {
	ctrl     *gomock.Controller
	recorder *MockIServiceMockRecorder
}

// MockIServiceMockRecorder is the mock recorder for MockIService.
type MockIServiceMockRecorder struct {
	mock *MockIService
}

// NewMockIService creates a new mock instance.
func NewMockIService(ctrl *gomock.Controller) *MockIService {
	mock := &MockIService{ctrl: ctrl}
	mock.recorder = &MockIServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIService) EXPECT() *MockIServiceMockRecorder {
	return m.recorder
}

// AddComment mocks base method.
func (m *MockIService) AddComment(ctx context.Context, comment string, pinId entity.PinID, userId entity.UserID) (entity.PinPageResponse, entity.CommentID, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddComment", ctx, comment, pinId, userId)
	ret0, _ := ret[0].(entity.PinPageResponse)
	ret1, _ := ret[1].(entity.CommentID)
	ret2, _ := ret[2].(errs.ErrorInfo)
	return ret0, ret1, ret2
}

// AddComment indicates an expected call of AddComment.
func (mr *MockIServiceMockRecorder) AddComment(ctx, comment, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddComment", reflect.TypeOf((*MockIService)(nil).AddComment), ctx, comment, pinId, userId)
}

// AddPinToBoard mocks base method.
func (m *MockIService) AddPinToBoard(ctx context.Context, boardId entity.BoardID, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPinToBoard", ctx, boardId, pinId, userId)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// AddPinToBoard indicates an expected call of AddPinToBoard.
func (mr *MockIServiceMockRecorder) AddPinToBoard(ctx, boardId, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPinToBoard", reflect.TypeOf((*MockIService)(nil).AddPinToBoard), ctx, boardId, pinId, userId)
}

// AddSubscriptionToUser mocks base method.
func (m *MockIService) AddSubscriptionToUser(ctx context.Context, userId, subscribeUserId entity.UserID) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSubscriptionToUser", ctx, userId, subscribeUserId)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// AddSubscriptionToUser indicates an expected call of AddSubscriptionToUser.
func (mr *MockIServiceMockRecorder) AddSubscriptionToUser(ctx, userId, subscribeUserId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubscriptionToUser", reflect.TypeOf((*MockIService)(nil).AddSubscriptionToUser), ctx, userId, subscribeUserId)
}

// CreateBoard mocks base method.
func (m *MockIService) CreateBoard(ctx context.Context, board entity.Board, userId entity.UserID) (entity.FullBoard, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBoard", ctx, board, userId)
	ret0, _ := ret[0].(entity.FullBoard)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// CreateBoard indicates an expected call of CreateBoard.
func (mr *MockIServiceMockRecorder) CreateBoard(ctx, board, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBoard", reflect.TypeOf((*MockIService)(nil).CreateBoard), ctx, board, userId)
}

// CreateMessage mocks base method.
func (m *MockIService) CreateMessage(ctx context.Context, message entity.Message) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessage", ctx, message)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// CreateMessage indicates an expected call of CreateMessage.
func (mr *MockIServiceMockRecorder) CreateMessage(ctx, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessage", reflect.TypeOf((*MockIService)(nil).CreateMessage), ctx, message)
}

// CreateNotification mocks base method.
func (m *MockIService) CreateNotification(ctx context.Context, n entity.Notification) (entity.NotificationID, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNotification", ctx, n)
	ret0, _ := ret[0].(entity.NotificationID)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// CreateNotification indicates an expected call of CreateNotification.
func (mr *MockIServiceMockRecorder) CreateNotification(ctx, n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNotification", reflect.TypeOf((*MockIService)(nil).CreateNotification), ctx, n)
}

// CreatePin mocks base method.
func (m *MockIService) CreatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePin", ctx, pin)
	ret0, _ := ret[0].(entity.PinPageResponse)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// CreatePin indicates an expected call of CreatePin.
func (mr *MockIServiceMockRecorder) CreatePin(ctx, pin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePin", reflect.TypeOf((*MockIService)(nil).CreatePin), ctx, pin)
}

// DeleteBoard mocks base method.
func (m *MockIService) DeleteBoard(ctx context.Context, boardId entity.BoardID, userId entity.UserID) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBoard", ctx, boardId, userId)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// DeleteBoard indicates an expected call of DeleteBoard.
func (mr *MockIServiceMockRecorder) DeleteBoard(ctx, boardId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBoard", reflect.TypeOf((*MockIService)(nil).DeleteBoard), ctx, boardId, userId)
}

// DeletePin mocks base method.
func (m *MockIService) DeletePin(ctx context.Context, pin entity.Pin) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePin", ctx, pin)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// DeletePin indicates an expected call of DeletePin.
func (mr *MockIServiceMockRecorder) DeletePin(ctx, pin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePin", reflect.TypeOf((*MockIService)(nil).DeletePin), ctx, pin)
}

// DeletePinFromBoard mocks base method.
func (m *MockIService) DeletePinFromBoard(ctx context.Context, boardId entity.BoardID, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePinFromBoard", ctx, boardId, pinId, userId)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// DeletePinFromBoard indicates an expected call of DeletePinFromBoard.
func (mr *MockIServiceMockRecorder) DeletePinFromBoard(ctx, boardId, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePinFromBoard", reflect.TypeOf((*MockIService)(nil).DeletePinFromBoard), ctx, boardId, pinId, userId)
}

// DeleteSubscriptionToUser mocks base method.
func (m *MockIService) DeleteSubscriptionToUser(ctx context.Context, userId, unsubscribeUserId entity.UserID) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSubscriptionToUser", ctx, userId, unsubscribeUserId)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// DeleteSubscriptionToUser indicates an expected call of DeleteSubscriptionToUser.
func (mr *MockIServiceMockRecorder) DeleteSubscriptionToUser(ctx, userId, unsubscribeUserId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSubscriptionToUser", reflect.TypeOf((*MockIService)(nil).DeleteSubscriptionToUser), ctx, userId, unsubscribeUserId)
}

// GetAllUsers mocks base method.
func (m *MockIService) GetAllUsers(ctx context.Context) (entity.SearchUsers, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", ctx)
	ret0, _ := ret[0].(entity.SearchUsers)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockIServiceMockRecorder) GetAllUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockIService)(nil).GetAllUsers), ctx)
}

// GetBoardById mocks base method.
func (m *MockIService) GetBoardById(ctx context.Context, boardId entity.BoardID, userId entity.UserID, limit, offset int) (entity.FullBoard, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardById", ctx, boardId, userId, limit, offset)
	ret0, _ := ret[0].(entity.FullBoard)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetBoardById indicates an expected call of GetBoardById.
func (mr *MockIServiceMockRecorder) GetBoardById(ctx, boardId, userId, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardById", reflect.TypeOf((*MockIService)(nil).GetBoardById), ctx, boardId, userId, limit, offset)
}

// GetComments mocks base method.
func (m *MockIService) GetComments(ctx context.Context, pinId entity.PinID) (entity.GetCommentsResponse, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComments", ctx, pinId)
	ret0, _ := ret[0].(entity.GetCommentsResponse)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetComments indicates an expected call of GetComments.
func (mr *MockIServiceMockRecorder) GetComments(ctx, pinId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComments", reflect.TypeOf((*MockIService)(nil).GetComments), ctx, pinId)
}

// GetFeedPins mocks base method.
func (m *MockIService) GetFeedPins(ctx context.Context, limit, offset int) (entity.FeedPins, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeedPins", ctx, limit, offset)
	ret0, _ := ret[0].(entity.FeedPins)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetFeedPins indicates an expected call of GetFeedPins.
func (mr *MockIServiceMockRecorder) GetFeedPins(ctx, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeedPins", reflect.TypeOf((*MockIService)(nil).GetFeedPins), ctx, limit, offset)
}

// GetMessages mocks base method.
func (m *MockIService) GetMessages(ctx context.Context, dialogUserId, authUserId entity.UserID) (entity.Messages, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessages", ctx, dialogUserId, authUserId)
	ret0, _ := ret[0].(entity.Messages)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetMessages indicates an expected call of GetMessages.
func (mr *MockIServiceMockRecorder) GetMessages(ctx, dialogUserId, authUserId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessages", reflect.TypeOf((*MockIService)(nil).GetMessages), ctx, dialogUserId, authUserId)
}

// GetPinById mocks base method.
func (m *MockIService) GetPinById(ctx context.Context, PinId entity.PinID, UserId entity.UserID) (entity.PinPageResponse, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPinById", ctx, PinId, UserId)
	ret0, _ := ret[0].(entity.PinPageResponse)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetPinById indicates an expected call of GetPinById.
func (mr *MockIServiceMockRecorder) GetPinById(ctx, PinId, UserId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPinById", reflect.TypeOf((*MockIService)(nil).GetPinById), ctx, PinId, UserId)
}

// GetSubscriptionsFeedPins mocks base method.
func (m *MockIService) GetSubscriptionsFeedPins(ctx context.Context, userId entity.UserID, limit, offset int) (entity.FeedPins, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubscriptionsFeedPins", ctx, userId, limit, offset)
	ret0, _ := ret[0].(entity.FeedPins)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetSubscriptionsFeedPins indicates an expected call of GetSubscriptionsFeedPins.
func (mr *MockIServiceMockRecorder) GetSubscriptionsFeedPins(ctx, userId, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionsFeedPins", reflect.TypeOf((*MockIService)(nil).GetSubscriptionsFeedPins), ctx, userId, limit, offset)
}

// GetUnreadNotifications mocks base method.
func (m *MockIService) GetUnreadNotifications(ctx context.Context, userId entity.UserID) (entity.Notifications, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnreadNotifications", ctx, userId)
	ret0, _ := ret[0].(entity.Notifications)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUnreadNotifications indicates an expected call of GetUnreadNotifications.
func (mr *MockIServiceMockRecorder) GetUnreadNotifications(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnreadNotifications", reflect.TypeOf((*MockIService)(nil).GetUnreadNotifications), ctx, userId)
}

// GetUserBoards mocks base method.
func (m *MockIService) GetUserBoards(ctx context.Context, authorNickname string, userId entity.UserID, limit, offset int) (entity.UserBoards, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBoards", ctx, authorNickname, userId, limit, offset)
	ret0, _ := ret[0].(entity.UserBoards)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserBoards indicates an expected call of GetUserBoards.
func (mr *MockIServiceMockRecorder) GetUserBoards(ctx, authorNickname, userId, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBoards", reflect.TypeOf((*MockIService)(nil).GetUserBoards), ctx, authorNickname, userId, limit, offset)
}

// GetUserBoardsWithoutPin mocks base method.
func (m *MockIService) GetUserBoardsWithoutPin(ctx context.Context, pinId entity.PinID, userId entity.UserID) (entity.UserBoardsWithoutPin, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBoardsWithoutPin", ctx, pinId, userId)
	ret0, _ := ret[0].(entity.UserBoardsWithoutPin)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserBoardsWithoutPin indicates an expected call of GetUserBoardsWithoutPin.
func (mr *MockIServiceMockRecorder) GetUserBoardsWithoutPin(ctx, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBoardsWithoutPin", reflect.TypeOf((*MockIService)(nil).GetUserBoardsWithoutPin), ctx, pinId, userId)
}

// GetUserByEmail mocks base method.
func (m *MockIService) GetUserByEmail(ctx context.Context, email string) (entity.User, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockIServiceMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockIService)(nil).GetUserByEmail), ctx, email)
}

// GetUserById mocks base method.
func (m *MockIService) GetUserById(ctx context.Context, id entity.UserID) (entity.User, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", ctx, id)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockIServiceMockRecorder) GetUserById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockIService)(nil).GetUserById), ctx, id)
}

// GetUserByNickname mocks base method.
func (m *MockIService) GetUserByNickname(ctx context.Context, nickname string) (entity.User, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByNickname", ctx, nickname)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserByNickname indicates an expected call of GetUserByNickname.
func (mr *MockIServiceMockRecorder) GetUserByNickname(ctx, nickname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByNickname", reflect.TypeOf((*MockIService)(nil).GetUserByNickname), ctx, nickname)
}

// GetUserChats mocks base method.
func (m *MockIService) GetUserChats(ctx context.Context, userId entity.UserID) (entity.UserChats, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserChats", ctx, userId)
	ret0, _ := ret[0].(entity.UserChats)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserChats indicates an expected call of GetUserChats.
func (mr *MockIServiceMockRecorder) GetUserChats(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserChats", reflect.TypeOf((*MockIService)(nil).GetUserChats), ctx, userId)
}

// GetUserPins mocks base method.
func (m *MockIService) GetUserPins(ctx context.Context, authorNickname string, limit, offset int) (entity.UserPins, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserPins", ctx, authorNickname, limit, offset)
	ret0, _ := ret[0].(entity.UserPins)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserPins indicates an expected call of GetUserPins.
func (mr *MockIServiceMockRecorder) GetUserPins(ctx, authorNickname, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserPins", reflect.TypeOf((*MockIService)(nil).GetUserPins), ctx, authorNickname, limit, offset)
}

// GetUserProfileByNickname mocks base method.
func (m *MockIService) GetUserProfileByNickname(ctx context.Context, nickname string, userId entity.UserID) (entity.UserProfileResponse, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfileByNickname", ctx, nickname, userId)
	ret0, _ := ret[0].(entity.UserProfileResponse)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserProfileByNickname indicates an expected call of GetUserProfileByNickname.
func (mr *MockIServiceMockRecorder) GetUserProfileByNickname(ctx, nickname, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfileByNickname", reflect.TypeOf((*MockIService)(nil).GetUserProfileByNickname), ctx, nickname, userId)
}

// GetUserSubscribers mocks base method.
func (m *MockIService) GetUserSubscribers(ctx context.Context, userId entity.UserID) (entity.UserSubscribers, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSubscribers", ctx, userId)
	ret0, _ := ret[0].(entity.UserSubscribers)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserSubscribers indicates an expected call of GetUserSubscribers.
func (mr *MockIServiceMockRecorder) GetUserSubscribers(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSubscribers", reflect.TypeOf((*MockIService)(nil).GetUserSubscribers), ctx, userId)
}

// GetUserSubscriptions mocks base method.
func (m *MockIService) GetUserSubscriptions(ctx context.Context, userId entity.UserID) (entity.UserSubscriptions, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSubscriptions", ctx, userId)
	ret0, _ := ret[0].(entity.UserSubscriptions)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUserSubscriptions indicates an expected call of GetUserSubscriptions.
func (mr *MockIServiceMockRecorder) GetUserSubscriptions(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSubscriptions", reflect.TypeOf((*MockIService)(nil).GetUserSubscriptions), ctx, userId)
}

// ReadAllNotifications mocks base method.
func (m *MockIService) ReadAllNotifications(ctx context.Context, userId entity.UserID) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAllNotifications", ctx, userId)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// ReadAllNotifications indicates an expected call of ReadAllNotifications.
func (mr *MockIServiceMockRecorder) ReadAllNotifications(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAllNotifications", reflect.TypeOf((*MockIService)(nil).ReadAllNotifications), ctx, userId)
}

// ReadNotification mocks base method.
func (m *MockIService) ReadNotification(ctx context.Context, notificationId entity.NotificationID, userId entity.UserID) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadNotification", ctx, notificationId, userId)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// ReadNotification indicates an expected call of ReadNotification.
func (mr *MockIServiceMockRecorder) ReadNotification(ctx, notificationId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadNotification", reflect.TypeOf((*MockIService)(nil).ReadNotification), ctx, notificationId, userId)
}

// RegisterUser mocks base method.
func (m *MockIService) RegisterUser(ctx context.Context, user entity.User) []errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, user)
	ret0, _ := ret[0].([]errs.ErrorInfo)
	return ret0
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockIServiceMockRecorder) RegisterUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockIService)(nil).RegisterUser), ctx, user)
}

// Search mocks base method.
func (m *MockIService) Search(ctx context.Context, query string) (entity.SearchResult, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, query)
	ret0, _ := ret[0].(entity.SearchResult)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockIServiceMockRecorder) Search(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockIService)(nil).Search), ctx, query)
}

// UpdateBoard mocks base method.
func (m *MockIService) UpdateBoard(ctx context.Context, board entity.Board, userId entity.UserID) (entity.FullBoard, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBoard", ctx, board, userId)
	ret0, _ := ret[0].(entity.FullBoard)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// UpdateBoard indicates an expected call of UpdateBoard.
func (mr *MockIServiceMockRecorder) UpdateBoard(ctx, board, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBoard", reflect.TypeOf((*MockIService)(nil).UpdateBoard), ctx, board, userId)
}

// UpdateDraft mocks base method.
func (m *MockIService) UpdateDraft(ctx context.Context, draft entity.Draft) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDraft", ctx, draft)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// UpdateDraft indicates an expected call of UpdateDraft.
func (mr *MockIServiceMockRecorder) UpdateDraft(ctx, draft interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDraft", reflect.TypeOf((*MockIService)(nil).UpdateDraft), ctx, draft)
}

// UpdatePin mocks base method.
func (m *MockIService) UpdatePin(ctx context.Context, pin entity.Pin) (entity.PinPageResponse, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePin", ctx, pin)
	ret0, _ := ret[0].(entity.PinPageResponse)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// UpdatePin indicates an expected call of UpdatePin.
func (mr *MockIServiceMockRecorder) UpdatePin(ctx, pin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePin", reflect.TypeOf((*MockIService)(nil).UpdatePin), ctx, pin)
}

// UpdateUser mocks base method.
func (m *MockIService) UpdateUser(ctx context.Context, user entity.User) (entity.User, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, user)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockIServiceMockRecorder) UpdateUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockIService)(nil).UpdateUser), ctx, user)
}
