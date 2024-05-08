// Code generated by MockGen. DO NOT EDIT.
// Source: internal/microservices/like/proto/like_grpc.pb.go

// Package mock_proto is a generated GoMock package.
package mock_proto

import (
	context "context"
	proto "harmonica/internal/microservices/like/proto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockLikeClient is a mock of LikeClient interface.
type MockLikeClient struct {
	ctrl     *gomock.Controller
	recorder *MockLikeClientMockRecorder
}

// MockLikeClientMockRecorder is the mock recorder for MockLikeClient.
type MockLikeClientMockRecorder struct {
	mock *MockLikeClient
}

// NewMockLikeClient creates a new mock instance.
func NewMockLikeClient(ctrl *gomock.Controller) *MockLikeClient {
	mock := &MockLikeClient{ctrl: ctrl}
	mock.recorder = &MockLikeClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLikeClient) EXPECT() *MockLikeClientMockRecorder {
	return m.recorder
}

// CheckIsLiked mocks base method.
func (m *MockLikeClient) CheckIsLiked(ctx context.Context, in *proto.CheckIsLikedRequest, opts ...grpc.CallOption) (*proto.CheckIsLikedResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckIsLiked", varargs...)
	ret0, _ := ret[0].(*proto.CheckIsLikedResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIsLiked indicates an expected call of CheckIsLiked.
func (mr *MockLikeClientMockRecorder) CheckIsLiked(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIsLiked", reflect.TypeOf((*MockLikeClient)(nil).CheckIsLiked), varargs...)
}

// ClearLike mocks base method.
func (m *MockLikeClient) ClearLike(ctx context.Context, in *proto.MakeLikeRequest, opts ...grpc.CallOption) (*proto.MakeLikeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ClearLike", varargs...)
	ret0, _ := ret[0].(*proto.MakeLikeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClearLike indicates an expected call of ClearLike.
func (mr *MockLikeClientMockRecorder) ClearLike(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearLike", reflect.TypeOf((*MockLikeClient)(nil).ClearLike), varargs...)
}

// GetFavorites mocks base method.
func (m *MockLikeClient) GetFavorites(ctx context.Context, in *proto.GetFavoritesRequest, opts ...grpc.CallOption) (*proto.GetFavoritesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFavorites", varargs...)
	ret0, _ := ret[0].(*proto.GetFavoritesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavorites indicates an expected call of GetFavorites.
func (mr *MockLikeClientMockRecorder) GetFavorites(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavorites", reflect.TypeOf((*MockLikeClient)(nil).GetFavorites), varargs...)
}

// GetUsersLiked mocks base method.
func (m *MockLikeClient) GetUsersLiked(ctx context.Context, in *proto.GetUsersLikedRequest, opts ...grpc.CallOption) (*proto.GetUsersLikedResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUsersLiked", varargs...)
	ret0, _ := ret[0].(*proto.GetUsersLikedResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersLiked indicates an expected call of GetUsersLiked.
func (mr *MockLikeClientMockRecorder) GetUsersLiked(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersLiked", reflect.TypeOf((*MockLikeClient)(nil).GetUsersLiked), varargs...)
}

// SetLike mocks base method.
func (m *MockLikeClient) SetLike(ctx context.Context, in *proto.MakeLikeRequest, opts ...grpc.CallOption) (*proto.MakeLikeResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetLike", varargs...)
	ret0, _ := ret[0].(*proto.MakeLikeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetLike indicates an expected call of SetLike.
func (mr *MockLikeClientMockRecorder) SetLike(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLike", reflect.TypeOf((*MockLikeClient)(nil).SetLike), varargs...)
}

// MockLikeServer is a mock of LikeServer interface.
type MockLikeServer struct {
	ctrl     *gomock.Controller
	recorder *MockLikeServerMockRecorder
}

// MockLikeServerMockRecorder is the mock recorder for MockLikeServer.
type MockLikeServerMockRecorder struct {
	mock *MockLikeServer
}

// NewMockLikeServer creates a new mock instance.
func NewMockLikeServer(ctrl *gomock.Controller) *MockLikeServer {
	mock := &MockLikeServer{ctrl: ctrl}
	mock.recorder = &MockLikeServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLikeServer) EXPECT() *MockLikeServerMockRecorder {
	return m.recorder
}

// CheckIsLiked mocks base method.
func (m *MockLikeServer) CheckIsLiked(arg0 context.Context, arg1 *proto.CheckIsLikedRequest) (*proto.CheckIsLikedResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIsLiked", arg0, arg1)
	ret0, _ := ret[0].(*proto.CheckIsLikedResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIsLiked indicates an expected call of CheckIsLiked.
func (mr *MockLikeServerMockRecorder) CheckIsLiked(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIsLiked", reflect.TypeOf((*MockLikeServer)(nil).CheckIsLiked), arg0, arg1)
}

// ClearLike mocks base method.
func (m *MockLikeServer) ClearLike(arg0 context.Context, arg1 *proto.MakeLikeRequest) (*proto.MakeLikeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearLike", arg0, arg1)
	ret0, _ := ret[0].(*proto.MakeLikeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClearLike indicates an expected call of ClearLike.
func (mr *MockLikeServerMockRecorder) ClearLike(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearLike", reflect.TypeOf((*MockLikeServer)(nil).ClearLike), arg0, arg1)
}

// GetFavorites mocks base method.
func (m *MockLikeServer) GetFavorites(arg0 context.Context, arg1 *proto.GetFavoritesRequest) (*proto.GetFavoritesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavorites", arg0, arg1)
	ret0, _ := ret[0].(*proto.GetFavoritesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavorites indicates an expected call of GetFavorites.
func (mr *MockLikeServerMockRecorder) GetFavorites(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavorites", reflect.TypeOf((*MockLikeServer)(nil).GetFavorites), arg0, arg1)
}

// GetUsersLiked mocks base method.
func (m *MockLikeServer) GetUsersLiked(arg0 context.Context, arg1 *proto.GetUsersLikedRequest) (*proto.GetUsersLikedResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersLiked", arg0, arg1)
	ret0, _ := ret[0].(*proto.GetUsersLikedResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersLiked indicates an expected call of GetUsersLiked.
func (mr *MockLikeServerMockRecorder) GetUsersLiked(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersLiked", reflect.TypeOf((*MockLikeServer)(nil).GetUsersLiked), arg0, arg1)
}

// SetLike mocks base method.
func (m *MockLikeServer) SetLike(arg0 context.Context, arg1 *proto.MakeLikeRequest) (*proto.MakeLikeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLike", arg0, arg1)
	ret0, _ := ret[0].(*proto.MakeLikeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetLike indicates an expected call of SetLike.
func (mr *MockLikeServerMockRecorder) SetLike(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLike", reflect.TypeOf((*MockLikeServer)(nil).SetLike), arg0, arg1)
}

// mustEmbedUnimplementedLikeServer mocks base method.
func (m *MockLikeServer) mustEmbedUnimplementedLikeServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedLikeServer")
}

// mustEmbedUnimplementedLikeServer indicates an expected call of mustEmbedUnimplementedLikeServer.
func (mr *MockLikeServerMockRecorder) mustEmbedUnimplementedLikeServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedLikeServer", reflect.TypeOf((*MockLikeServer)(nil).mustEmbedUnimplementedLikeServer))
}

// MockUnsafeLikeServer is a mock of UnsafeLikeServer interface.
type MockUnsafeLikeServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeLikeServerMockRecorder
}

// MockUnsafeLikeServerMockRecorder is the mock recorder for MockUnsafeLikeServer.
type MockUnsafeLikeServerMockRecorder struct {
	mock *MockUnsafeLikeServer
}

// NewMockUnsafeLikeServer creates a new mock instance.
func NewMockUnsafeLikeServer(ctrl *gomock.Controller) *MockUnsafeLikeServer {
	mock := &MockUnsafeLikeServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeLikeServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeLikeServer) EXPECT() *MockUnsafeLikeServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedLikeServer mocks base method.
func (m *MockUnsafeLikeServer) mustEmbedUnimplementedLikeServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedLikeServer")
}

// mustEmbedUnimplementedLikeServer indicates an expected call of mustEmbedUnimplementedLikeServer.
func (mr *MockUnsafeLikeServerMockRecorder) mustEmbedUnimplementedLikeServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedLikeServer", reflect.TypeOf((*MockUnsafeLikeServer)(nil).mustEmbedUnimplementedLikeServer))
}
