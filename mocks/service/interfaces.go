// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/interfaces.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	entity "harmonica/internal/entity"
	errs "harmonica/internal/entity/errs"
	multipart "mime/multipart"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	minio "github.com/minio/minio-go/v7"
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

// ClearLike mocks base method.
func (m *MockIService) ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearLike", ctx, pinId, userId)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// ClearLike indicates an expected call of ClearLike.
func (mr *MockIServiceMockRecorder) ClearLike(ctx, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearLike", reflect.TypeOf((*MockIService)(nil).ClearLike), ctx, pinId, userId)
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

// GetImage mocks base method.
func (m *MockIService) GetImage(ctx context.Context, name string) (*minio.Object, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImage", ctx, name)
	ret0, _ := ret[0].(*minio.Object)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImage indicates an expected call of GetImage.
func (mr *MockIServiceMockRecorder) GetImage(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImage", reflect.TypeOf((*MockIService)(nil).GetImage), ctx, name)
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

// GetUsersLiked mocks base method.
func (m *MockIService) GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersLiked", ctx, pinId, limit)
	ret0, _ := ret[0].(entity.UserList)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetUsersLiked indicates an expected call of GetUsersLiked.
func (mr *MockIServiceMockRecorder) GetUsersLiked(ctx, pinId, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersLiked", reflect.TypeOf((*MockIService)(nil).GetUsersLiked), ctx, pinId, limit)
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

// SetLike mocks base method.
func (m *MockIService) SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) errs.ErrorInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLike", ctx, pinId, userId)
	ret0, _ := ret[0].(errs.ErrorInfo)
	return ret0
}

// SetLike indicates an expected call of SetLike.
func (mr *MockIServiceMockRecorder) SetLike(ctx, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLike", reflect.TypeOf((*MockIService)(nil).SetLike), ctx, pinId, userId)
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

// UploadImage mocks base method.
func (m *MockIService) UploadImage(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", ctx, file, fileHeader)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockIServiceMockRecorder) UploadImage(ctx, file, fileHeader interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockIService)(nil).UploadImage), ctx, file, fileHeader)
}