// Code generated by MockGen. DO NOT EDIT.
// Source: internal/microservices/like/server/service/interfaces.go

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

// CheckIsLiked mocks base method.
func (m *MockIService) CheckIsLiked(ctx context.Context, pinId entity.PinID, userId entity.UserID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIsLiked", ctx, pinId, userId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIsLiked indicates an expected call of CheckIsLiked.
func (mr *MockIServiceMockRecorder) CheckIsLiked(ctx, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIsLiked", reflect.TypeOf((*MockIService)(nil).CheckIsLiked), ctx, pinId, userId)
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

// GetFavorites mocks base method.
func (m *MockIService) GetFavorites(ctx context.Context, userId entity.UserID, limit, offset int) (entity.FeedPins, errs.ErrorInfo) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavorites", ctx, userId, limit, offset)
	ret0, _ := ret[0].(entity.FeedPins)
	ret1, _ := ret[1].(errs.ErrorInfo)
	return ret0, ret1
}

// GetFavorites indicates an expected call of GetFavorites.
func (mr *MockIServiceMockRecorder) GetFavorites(ctx, userId, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavorites", reflect.TypeOf((*MockIService)(nil).GetFavorites), ctx, userId, limit, offset)
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
