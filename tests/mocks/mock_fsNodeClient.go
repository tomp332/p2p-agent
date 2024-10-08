// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/nodes/fsNode/client.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pb "github.com/tomp332/p2p-agent/pkg/pb"
	types "github.com/tomp332/p2p-agent/pkg/utils/types"
)

// MockFileNodeType is a mock of FileNodeType interface.
type MockFileNodeType struct {
	ctrl     *gomock.Controller
	recorder *MockFileNodeTypeMockRecorder
}

// MockFileNodeTypeMockRecorder is the mock recorder for MockFileNodeType.
type MockFileNodeTypeMockRecorder struct {
	mock *MockFileNodeType
}

// NewMockFileNodeType creates a new mock instance.
func NewMockFileNodeType(ctrl *gomock.Controller) *MockFileNodeType {
	mock := &MockFileNodeType{ctrl: ctrl}
	mock.recorder = &MockFileNodeTypeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileNodeType) EXPECT() *MockFileNodeTypeMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockFileNodeType) Authenticate(ctx context.Context, username, password string) (*pb.AuthenticateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", ctx, username, password)
	ret0, _ := ret[0].(*pb.AuthenticateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockFileNodeTypeMockRecorder) Authenticate(ctx, username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockFileNodeType)(nil).Authenticate), ctx, username, password)
}

// Connect mocks base method.
func (m *MockFileNodeType) Connect() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connect indicates an expected call of Connect.
func (mr *MockFileNodeTypeMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockFileNodeType)(nil).Connect))
}

// Disconnect mocks base method.
func (m *MockFileNodeType) Disconnect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockFileNodeTypeMockRecorder) Disconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockFileNodeType)(nil).Disconnect))
}

// DownloadFile mocks base method.
func (m *MockFileNodeType) DownloadFile(searchCtx context.Context, fileName string) (<-chan *types.TransferChunkData, <-chan error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFile", searchCtx, fileName)
	ret0, _ := ret[0].(<-chan *types.TransferChunkData)
	ret1, _ := ret[1].(<-chan error)
	return ret0, ret1
}

// DownloadFile indicates an expected call of DownloadFile.
func (mr *MockFileNodeTypeMockRecorder) DownloadFile(searchCtx, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockFileNodeType)(nil).DownloadFile), searchCtx, fileName)
}

// SearchFile mocks base method.
func (m *MockFileNodeType) SearchFile(searchCtx context.Context, fileName string) (*pb.SearchFileResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchFile", searchCtx, fileName)
	ret0, _ := ret[0].(*pb.SearchFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchFile indicates an expected call of SearchFile.
func (mr *MockFileNodeTypeMockRecorder) SearchFile(searchCtx, fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchFile", reflect.TypeOf((*MockFileNodeType)(nil).SearchFile), searchCtx, fileName)
}

// UploadFile mocks base method.
func (m *MockFileNodeType) UploadFile(ctx context.Context, filePath string) <-chan error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", ctx, filePath)
	ret0, _ := ret[0].(<-chan error)
	return ret0
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockFileNodeTypeMockRecorder) UploadFile(ctx, filePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockFileNodeType)(nil).UploadFile), ctx, filePath)
}
