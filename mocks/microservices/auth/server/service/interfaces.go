// Code generated by MockGen. DO NOT EDIT.
// Source: internal/microservices/auth/server/service/interfaces.go

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
