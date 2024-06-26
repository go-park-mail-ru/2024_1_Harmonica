// Code generated by MockGen. DO NOT EDIT.
// Source: internal/microservices/like/server/repository/interfaces.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	entity "harmonica/internal/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIRepository is a mock of IRepository interface.
type MockIRepository struct {
	ctrl     *gomock.Controller
	recorder *MockIRepositoryMockRecorder
}

// MockIRepositoryMockRecorder is the mock recorder for MockIRepository.
type MockIRepositoryMockRecorder struct {
	mock *MockIRepository
}

// NewMockIRepository creates a new mock instance.
func NewMockIRepository(ctrl *gomock.Controller) *MockIRepository {
	mock := &MockIRepository{ctrl: ctrl}
	mock.recorder = &MockIRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRepository) EXPECT() *MockIRepositoryMockRecorder {
	return m.recorder
}

// CheckIsLiked mocks base method.
func (m *MockIRepository) CheckIsLiked(ctx context.Context, pinId entity.PinID, userId entity.UserID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIsLiked", ctx, pinId, userId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIsLiked indicates an expected call of CheckIsLiked.
func (mr *MockIRepositoryMockRecorder) CheckIsLiked(ctx, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIsLiked", reflect.TypeOf((*MockIRepository)(nil).CheckIsLiked), ctx, pinId, userId)
}

// ClearLike mocks base method.
func (m *MockIRepository) ClearLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearLike", ctx, pinId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearLike indicates an expected call of ClearLike.
func (mr *MockIRepositoryMockRecorder) ClearLike(ctx, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearLike", reflect.TypeOf((*MockIRepository)(nil).ClearLike), ctx, pinId, userId)
}

// GetFavorites mocks base method.
func (m *MockIRepository) GetFavorites(ctx context.Context, userId entity.UserID, limit, offset int) (entity.FeedPins, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavorites", ctx, userId, limit, offset)
	ret0, _ := ret[0].(entity.FeedPins)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavorites indicates an expected call of GetFavorites.
func (mr *MockIRepositoryMockRecorder) GetFavorites(ctx, userId, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavorites", reflect.TypeOf((*MockIRepository)(nil).GetFavorites), ctx, userId, limit, offset)
}

// GetUsersLiked mocks base method.
func (m *MockIRepository) GetUsersLiked(ctx context.Context, pinId entity.PinID, limit int) (entity.UserList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersLiked", ctx, pinId, limit)
	ret0, _ := ret[0].(entity.UserList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersLiked indicates an expected call of GetUsersLiked.
func (mr *MockIRepositoryMockRecorder) GetUsersLiked(ctx, pinId, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersLiked", reflect.TypeOf((*MockIRepository)(nil).GetUsersLiked), ctx, pinId, limit)
}

// SetLike mocks base method.
func (m *MockIRepository) SetLike(ctx context.Context, pinId entity.PinID, userId entity.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLike", ctx, pinId, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLike indicates an expected call of SetLike.
func (mr *MockIRepositoryMockRecorder) SetLike(ctx, pinId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLike", reflect.TypeOf((*MockIRepository)(nil).SetLike), ctx, pinId, userId)
}
