// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: protos/files_node.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	FilesNodeService_UploadFile_FullMethodName          = "/p2p_agent.FilesNodeService/UploadFile"
	FilesNodeService_DownloadFile_FullMethodName        = "/p2p_agent.FilesNodeService/DownloadFile"
	FilesNodeService_DeleteFile_FullMethodName          = "/p2p_agent.FilesNodeService/DeleteFile"
	FilesNodeService_BroadcastSearchFile_FullMethodName = "/p2p_agent.FilesNodeService/BroadcastSearchFile"
)

// FilesNodeServiceClient is the client API for FilesNodeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FilesNodeServiceClient interface {
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (FilesNodeService_UploadFileClient, error)
	DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (FilesNodeService_DownloadFileClient, error)
	DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error)
	BroadcastSearchFile(ctx context.Context, in *BroadcastRequest, opts ...grpc.CallOption) (*BroadcastResponse, error)
}

type filesNodeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFilesNodeServiceClient(cc grpc.ClientConnInterface) FilesNodeServiceClient {
	return &filesNodeServiceClient{cc}
}

func (c *filesNodeServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (FilesNodeService_UploadFileClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FilesNodeService_ServiceDesc.Streams[0], FilesNodeService_UploadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &filesNodeServiceUploadFileClient{ClientStream: stream}
	return x, nil
}

type FilesNodeService_UploadFileClient interface {
	Send(*UploadFileRequest) error
	CloseAndRecv() (*UploadFileResponse, error)
	grpc.ClientStream
}

type filesNodeServiceUploadFileClient struct {
	grpc.ClientStream
}

func (x *filesNodeServiceUploadFileClient) Send(m *UploadFileRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *filesNodeServiceUploadFileClient) CloseAndRecv() (*UploadFileResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filesNodeServiceClient) DownloadFile(ctx context.Context, in *DownloadFileRequest, opts ...grpc.CallOption) (FilesNodeService_DownloadFileClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FilesNodeService_ServiceDesc.Streams[1], FilesNodeService_DownloadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &filesNodeServiceDownloadFileClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FilesNodeService_DownloadFileClient interface {
	Recv() (*DownloadFileResponse, error)
	grpc.ClientStream
}

type filesNodeServiceDownloadFileClient struct {
	grpc.ClientStream
}

func (x *filesNodeServiceDownloadFileClient) Recv() (*DownloadFileResponse, error) {
	m := new(DownloadFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *filesNodeServiceClient) DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteFileResponse)
	err := c.cc.Invoke(ctx, FilesNodeService_DeleteFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *filesNodeServiceClient) BroadcastSearchFile(ctx context.Context, in *BroadcastRequest, opts ...grpc.CallOption) (*BroadcastResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BroadcastResponse)
	err := c.cc.Invoke(ctx, FilesNodeService_BroadcastSearchFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilesNodeServiceServer is the server API for FilesNodeService service.
// All implementations must embed UnimplementedFilesNodeServiceServer
// for forward compatibility
type FilesNodeServiceServer interface {
	UploadFile(FilesNodeService_UploadFileServer) error
	DownloadFile(*DownloadFileRequest, FilesNodeService_DownloadFileServer) error
	DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error)
	BroadcastSearchFile(context.Context, *BroadcastRequest) (*BroadcastResponse, error)
	mustEmbedUnimplementedFilesNodeServiceServer()
}

// UnimplementedFilesNodeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFilesNodeServiceServer struct {
}

func (UnimplementedFilesNodeServiceServer) UploadFile(FilesNodeService_UploadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedFilesNodeServiceServer) DownloadFile(*DownloadFileRequest, FilesNodeService_DownloadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadFile not implemented")
}
func (UnimplementedFilesNodeServiceServer) DeleteFile(context.Context, *DeleteFileRequest) (*DeleteFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFile not implemented")
}
func (UnimplementedFilesNodeServiceServer) BroadcastSearchFile(context.Context, *BroadcastRequest) (*BroadcastResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastSearchFile not implemented")
}
func (UnimplementedFilesNodeServiceServer) mustEmbedUnimplementedFilesNodeServiceServer() {}

// UnsafeFilesNodeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FilesNodeServiceServer will
// result in compilation errors.
type UnsafeFilesNodeServiceServer interface {
	mustEmbedUnimplementedFilesNodeServiceServer()
}

func RegisterFilesNodeServiceServer(s grpc.ServiceRegistrar, srv FilesNodeServiceServer) {
	s.RegisterService(&FilesNodeService_ServiceDesc, srv)
}

func _FilesNodeService_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FilesNodeServiceServer).UploadFile(&filesNodeServiceUploadFileServer{ServerStream: stream})
}

type FilesNodeService_UploadFileServer interface {
	SendAndClose(*UploadFileResponse) error
	Recv() (*UploadFileRequest, error)
	grpc.ServerStream
}

type filesNodeServiceUploadFileServer struct {
	grpc.ServerStream
}

func (x *filesNodeServiceUploadFileServer) SendAndClose(m *UploadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *filesNodeServiceUploadFileServer) Recv() (*UploadFileRequest, error) {
	m := new(UploadFileRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FilesNodeService_DownloadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FilesNodeServiceServer).DownloadFile(m, &filesNodeServiceDownloadFileServer{ServerStream: stream})
}

type FilesNodeService_DownloadFileServer interface {
	Send(*DownloadFileResponse) error
	grpc.ServerStream
}

type filesNodeServiceDownloadFileServer struct {
	grpc.ServerStream
}

func (x *filesNodeServiceDownloadFileServer) Send(m *DownloadFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _FilesNodeService_DeleteFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilesNodeServiceServer).DeleteFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FilesNodeService_DeleteFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilesNodeServiceServer).DeleteFile(ctx, req.(*DeleteFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _FilesNodeService_BroadcastSearchFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BroadcastRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilesNodeServiceServer).BroadcastSearchFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FilesNodeService_BroadcastSearchFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilesNodeServiceServer).BroadcastSearchFile(ctx, req.(*BroadcastRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FilesNodeService_ServiceDesc is the grpc.ServiceDesc for FilesNodeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FilesNodeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "p2p_agent.FilesNodeService",
	HandlerType: (*FilesNodeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteFile",
			Handler:    _FilesNodeService_DeleteFile_Handler,
		},
		{
			MethodName: "BroadcastSearchFile",
			Handler:    _FilesNodeService_BroadcastSearchFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _FilesNodeService_UploadFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadFile",
			Handler:       _FilesNodeService_DownloadFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protos/files_node.proto",
}
