// Code generated by MockGen. DO NOT EDIT.
// Source: internal/microservices/image/proto/image_grpc.pb.go

// Package mock_proto is a generated GoMock package.
package mock_proto

import (
	context "context"
	proto "harmonica/internal/microservices/image/proto"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockImageClient is a mock of ImageClient interface.
type MockImageClient struct {
	ctrl     *gomock.Controller
	recorder *MockImageClientMockRecorder
}

// MockImageClientMockRecorder is the mock recorder for MockImageClient.
type MockImageClientMockRecorder struct {
	mock *MockImageClient
}

// NewMockImageClient creates a new mock instance.
func NewMockImageClient(ctrl *gomock.Controller) *MockImageClient {
	mock := &MockImageClient{ctrl: ctrl}
	mock.recorder = &MockImageClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageClient) EXPECT() *MockImageClientMockRecorder {
	return m.recorder
}

// FormUrl mocks base method.
func (m *MockImageClient) FormUrl(ctx context.Context, in *proto.FormUrlRequest, opts ...grpc.CallOption) (*proto.FormUrlResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FormUrl", varargs...)
	ret0, _ := ret[0].(*proto.FormUrlResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FormUrl indicates an expected call of FormUrl.
func (mr *MockImageClientMockRecorder) FormUrl(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FormUrl", reflect.TypeOf((*MockImageClient)(nil).FormUrl), varargs...)
}

// GetImage mocks base method.
func (m *MockImageClient) GetImage(ctx context.Context, in *proto.GetImageRequest, opts ...grpc.CallOption) (*proto.GetImageResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetImage", varargs...)
	ret0, _ := ret[0].(*proto.GetImageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImage indicates an expected call of GetImage.
func (mr *MockImageClientMockRecorder) GetImage(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImage", reflect.TypeOf((*MockImageClient)(nil).GetImage), varargs...)
}

// GetImageBounds mocks base method.
func (m *MockImageClient) GetImageBounds(ctx context.Context, in *proto.GetImageBoundsRequest, opts ...grpc.CallOption) (*proto.GetImageBoundsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetImageBounds", varargs...)
	ret0, _ := ret[0].(*proto.GetImageBoundsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImageBounds indicates an expected call of GetImageBounds.
func (mr *MockImageClientMockRecorder) GetImageBounds(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageBounds", reflect.TypeOf((*MockImageClient)(nil).GetImageBounds), varargs...)
}

// UploadImage mocks base method.
func (m *MockImageClient) UploadImage(ctx context.Context, in *proto.UploadImageRequest, opts ...grpc.CallOption) (*proto.UploadImageResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UploadImage", varargs...)
	ret0, _ := ret[0].(*proto.UploadImageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockImageClientMockRecorder) UploadImage(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockImageClient)(nil).UploadImage), varargs...)
}

// MockImageServer is a mock of ImageServer interface.
type MockImageServer struct {
	ctrl     *gomock.Controller
	recorder *MockImageServerMockRecorder
}

// MockImageServerMockRecorder is the mock recorder for MockImageServer.
type MockImageServerMockRecorder struct {
	mock *MockImageServer
}

// NewMockImageServer creates a new mock instance.
func NewMockImageServer(ctrl *gomock.Controller) *MockImageServer {
	mock := &MockImageServer{ctrl: ctrl}
	mock.recorder = &MockImageServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageServer) EXPECT() *MockImageServerMockRecorder {
	return m.recorder
}

// FormUrl mocks base method.
func (m *MockImageServer) FormUrl(arg0 context.Context, arg1 *proto.FormUrlRequest) (*proto.FormUrlResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FormUrl", arg0, arg1)
	ret0, _ := ret[0].(*proto.FormUrlResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FormUrl indicates an expected call of FormUrl.
func (mr *MockImageServerMockRecorder) FormUrl(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FormUrl", reflect.TypeOf((*MockImageServer)(nil).FormUrl), arg0, arg1)
}

// GetImage mocks base method.
func (m *MockImageServer) GetImage(arg0 context.Context, arg1 *proto.GetImageRequest) (*proto.GetImageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImage", arg0, arg1)
	ret0, _ := ret[0].(*proto.GetImageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImage indicates an expected call of GetImage.
func (mr *MockImageServerMockRecorder) GetImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImage", reflect.TypeOf((*MockImageServer)(nil).GetImage), arg0, arg1)
}

// GetImageBounds mocks base method.
func (m *MockImageServer) GetImageBounds(arg0 context.Context, arg1 *proto.GetImageBoundsRequest) (*proto.GetImageBoundsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageBounds", arg0, arg1)
	ret0, _ := ret[0].(*proto.GetImageBoundsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImageBounds indicates an expected call of GetImageBounds.
func (mr *MockImageServerMockRecorder) GetImageBounds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageBounds", reflect.TypeOf((*MockImageServer)(nil).GetImageBounds), arg0, arg1)
}

// UploadImage mocks base method.
func (m *MockImageServer) UploadImage(arg0 context.Context, arg1 *proto.UploadImageRequest) (*proto.UploadImageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadImage", arg0, arg1)
	ret0, _ := ret[0].(*proto.UploadImageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadImage indicates an expected call of UploadImage.
func (mr *MockImageServerMockRecorder) UploadImage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadImage", reflect.TypeOf((*MockImageServer)(nil).UploadImage), arg0, arg1)
}

// mustEmbedUnimplementedImageServer mocks base method.
func (m *MockImageServer) mustEmbedUnimplementedImageServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedImageServer")
}

// mustEmbedUnimplementedImageServer indicates an expected call of mustEmbedUnimplementedImageServer.
func (mr *MockImageServerMockRecorder) mustEmbedUnimplementedImageServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedImageServer", reflect.TypeOf((*MockImageServer)(nil).mustEmbedUnimplementedImageServer))
}

// MockUnsafeImageServer is a mock of UnsafeImageServer interface.
type MockUnsafeImageServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeImageServerMockRecorder
}

// MockUnsafeImageServerMockRecorder is the mock recorder for MockUnsafeImageServer.
type MockUnsafeImageServerMockRecorder struct {
	mock *MockUnsafeImageServer
}

// NewMockUnsafeImageServer creates a new mock instance.
func NewMockUnsafeImageServer(ctrl *gomock.Controller) *MockUnsafeImageServer {
	mock := &MockUnsafeImageServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeImageServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeImageServer) EXPECT() *MockUnsafeImageServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedImageServer mocks base method.
func (m *MockUnsafeImageServer) mustEmbedUnimplementedImageServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedImageServer")
}

// mustEmbedUnimplementedImageServer indicates an expected call of mustEmbedUnimplementedImageServer.
func (mr *MockUnsafeImageServerMockRecorder) mustEmbedUnimplementedImageServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedImageServer", reflect.TypeOf((*MockUnsafeImageServer)(nil).mustEmbedUnimplementedImageServer))
}
