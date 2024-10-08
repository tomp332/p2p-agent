// Code generated by MockGen. DO NOT EDIT.
// Source: ./tests/interfaces/fsNodeClient.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pb "github.com/tomp332/p2p-agent/pkg/pb"
	metadata "google.golang.org/grpc/metadata"
)

// MockUploadFileServerStream is a mock of UploadFileServerStream interface.
type MockUploadFileServerStream struct {
	ctrl     *gomock.Controller
	recorder *MockUploadFileServerStreamMockRecorder
}

// MockUploadFileServerStreamMockRecorder is the mock recorder for MockUploadFileServerStream.
type MockUploadFileServerStreamMockRecorder struct {
	mock *MockUploadFileServerStream
}

// NewMockUploadFileServerStream creates a new mock instance.
func NewMockUploadFileServerStream(ctrl *gomock.Controller) *MockUploadFileServerStream {
	mock := &MockUploadFileServerStream{ctrl: ctrl}
	mock.recorder = &MockUploadFileServerStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUploadFileServerStream) EXPECT() *MockUploadFileServerStreamMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockUploadFileServerStream) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockUploadFileServerStreamMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockUploadFileServerStream)(nil).Context))
}

// Recv mocks base method.
func (m *MockUploadFileServerStream) Recv() (*pb.UploadFileRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*pb.UploadFileRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockUploadFileServerStreamMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockUploadFileServerStream)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockUploadFileServerStream) RecvMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockUploadFileServerStreamMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockUploadFileServerStream)(nil).RecvMsg), m)
}

// SendAndClose mocks base method.
func (m *MockUploadFileServerStream) SendAndClose(arg0 *pb.UploadFileResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAndClose", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAndClose indicates an expected call of SendAndClose.
func (mr *MockUploadFileServerStreamMockRecorder) SendAndClose(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAndClose", reflect.TypeOf((*MockUploadFileServerStream)(nil).SendAndClose), arg0)
}

// SendHeader mocks base method.
func (m *MockUploadFileServerStream) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockUploadFileServerStreamMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockUploadFileServerStream)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockUploadFileServerStream) SendMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockUploadFileServerStreamMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockUploadFileServerStream)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockUploadFileServerStream) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockUploadFileServerStreamMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockUploadFileServerStream)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockUploadFileServerStream) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockUploadFileServerStreamMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockUploadFileServerStream)(nil).SetTrailer), arg0)
}

// MockDownloadFileClientStream is a mock of DownloadFileClientStream interface.
type MockDownloadFileClientStream struct {
	ctrl     *gomock.Controller
	recorder *MockDownloadFileClientStreamMockRecorder
}

// MockDownloadFileClientStreamMockRecorder is the mock recorder for MockDownloadFileClientStream.
type MockDownloadFileClientStreamMockRecorder struct {
	mock *MockDownloadFileClientStream
}

// NewMockDownloadFileClientStream creates a new mock instance.
func NewMockDownloadFileClientStream(ctrl *gomock.Controller) *MockDownloadFileClientStream {
	mock := &MockDownloadFileClientStream{ctrl: ctrl}
	mock.recorder = &MockDownloadFileClientStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDownloadFileClientStream) EXPECT() *MockDownloadFileClientStreamMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockDownloadFileClientStream) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockDownloadFileClientStreamMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockDownloadFileClientStream)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockDownloadFileClientStream) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDownloadFileClientStreamMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDownloadFileClientStream)(nil).Context))
}

// Header mocks base method.
func (m *MockDownloadFileClientStream) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockDownloadFileClientStreamMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockDownloadFileClientStream)(nil).Header))
}

// Recv mocks base method.
func (m *MockDownloadFileClientStream) Recv() (*pb.DownloadFileResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*pb.DownloadFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockDownloadFileClientStreamMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockDownloadFileClientStream)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockDownloadFileClientStream) RecvMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDownloadFileClientStreamMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDownloadFileClientStream)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockDownloadFileClientStream) Send(arg0 *pb.DownloadFileResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockDownloadFileClientStreamMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDownloadFileClientStream)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockDownloadFileClientStream) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockDownloadFileClientStreamMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockDownloadFileClientStream)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockDownloadFileClientStream) SendMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDownloadFileClientStreamMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDownloadFileClientStream)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockDownloadFileClientStream) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockDownloadFileClientStreamMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockDownloadFileClientStream)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockDownloadFileClientStream) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockDownloadFileClientStreamMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockDownloadFileClientStream)(nil).SetTrailer), arg0)
}

// Trailer mocks base method.
func (m *MockDownloadFileClientStream) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockDownloadFileClientStreamMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockDownloadFileClientStream)(nil).Trailer))
}

// MockDownloadFileServerStream is a mock of DownloadFileServerStream interface.
type MockDownloadFileServerStream struct {
	ctrl     *gomock.Controller
	recorder *MockDownloadFileServerStreamMockRecorder
}

// MockDownloadFileServerStreamMockRecorder is the mock recorder for MockDownloadFileServerStream.
type MockDownloadFileServerStreamMockRecorder struct {
	mock *MockDownloadFileServerStream
}

// NewMockDownloadFileServerStream creates a new mock instance.
func NewMockDownloadFileServerStream(ctrl *gomock.Controller) *MockDownloadFileServerStream {
	mock := &MockDownloadFileServerStream{ctrl: ctrl}
	mock.recorder = &MockDownloadFileServerStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDownloadFileServerStream) EXPECT() *MockDownloadFileServerStreamMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockDownloadFileServerStream) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDownloadFileServerStreamMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDownloadFileServerStream)(nil).Context))
}

// Recv mocks base method.
func (m *MockDownloadFileServerStream) Recv() (*pb.DownloadFileRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*pb.DownloadFileRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockDownloadFileServerStreamMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockDownloadFileServerStream)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockDownloadFileServerStream) RecvMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDownloadFileServerStreamMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDownloadFileServerStream)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockDownloadFileServerStream) Send(arg0 *pb.DownloadFileResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockDownloadFileServerStreamMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDownloadFileServerStream)(nil).Send), arg0)
}

// SendAndClose mocks base method.
func (m *MockDownloadFileServerStream) SendAndClose(response *pb.DownloadFileResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAndClose", response)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAndClose indicates an expected call of SendAndClose.
func (mr *MockDownloadFileServerStreamMockRecorder) SendAndClose(response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAndClose", reflect.TypeOf((*MockDownloadFileServerStream)(nil).SendAndClose), response)
}

// SendHeader mocks base method.
func (m *MockDownloadFileServerStream) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockDownloadFileServerStreamMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockDownloadFileServerStream)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockDownloadFileServerStream) SendMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDownloadFileServerStreamMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDownloadFileServerStream)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockDownloadFileServerStream) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockDownloadFileServerStreamMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockDownloadFileServerStream)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockDownloadFileServerStream) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockDownloadFileServerStreamMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockDownloadFileServerStream)(nil).SetTrailer), arg0)
}
