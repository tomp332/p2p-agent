// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: protos/files_node.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Message for file upload request
type UploadFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// File metadata
	FileName  string `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	ChunkData []byte `protobuf:"bytes,2,opt,name=chunk_data,json=chunkData,proto3" json:"chunk_data,omitempty"`
}

func (x *UploadFileRequest) Reset() {
	*x = UploadFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_node_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileRequest) ProtoMessage() {}

func (x *UploadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_node_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileRequest.ProtoReflect.Descriptor instead.
func (*UploadFileRequest) Descriptor() ([]byte, []int) {
	return file_protos_files_node_proto_rawDescGZIP(), []int{0}
}

func (x *UploadFileRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *UploadFileRequest) GetChunkData() []byte {
	if x != nil {
		return x.ChunkData
	}
	return nil
}

// Message for file upload response
type UploadFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId   string  `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	FileSize float64 `protobuf:"fixed64,2,opt,name=file_size,json=fileSize,proto3" json:"file_size,omitempty"`
	Success  bool    `protobuf:"varint,3,opt,name=success,proto3" json:"success,omitempty"`
	Message  string  `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *UploadFileResponse) Reset() {
	*x = UploadFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_node_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadFileResponse) ProtoMessage() {}

func (x *UploadFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_node_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadFileResponse.ProtoReflect.Descriptor instead.
func (*UploadFileResponse) Descriptor() ([]byte, []int) {
	return file_protos_files_node_proto_rawDescGZIP(), []int{1}
}

func (x *UploadFileResponse) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *UploadFileResponse) GetFileSize() float64 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

func (x *UploadFileResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UploadFileResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type DownloadFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId string `protobuf:"bytes,1,opt,name=fileId,proto3" json:"fileId,omitempty"`
}

func (x *DownloadFileRequest) Reset() {
	*x = DownloadFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_node_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadFileRequest) ProtoMessage() {}

func (x *DownloadFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_node_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadFileRequest.ProtoReflect.Descriptor instead.
func (*DownloadFileRequest) Descriptor() ([]byte, []int) {
	return file_protos_files_node_proto_rawDescGZIP(), []int{2}
}

func (x *DownloadFileRequest) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

type DownloadFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Chunk  []byte `protobuf:"bytes,1,opt,name=chunk,proto3" json:"chunk,omitempty"`
	FileId string `protobuf:"bytes,2,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	Exists bool   `protobuf:"varint,3,opt,name=exists,proto3" json:"exists,omitempty"`
}

func (x *DownloadFileResponse) Reset() {
	*x = DownloadFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_node_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadFileResponse) ProtoMessage() {}

func (x *DownloadFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_node_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadFileResponse.ProtoReflect.Descriptor instead.
func (*DownloadFileResponse) Descriptor() ([]byte, []int) {
	return file_protos_files_node_proto_rawDescGZIP(), []int{3}
}

func (x *DownloadFileResponse) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

func (x *DownloadFileResponse) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *DownloadFileResponse) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}

type DeleteFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId string `protobuf:"bytes,1,opt,name=fileId,proto3" json:"fileId,omitempty"`
}

func (x *DeleteFileRequest) Reset() {
	*x = DeleteFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_node_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileRequest) ProtoMessage() {}

func (x *DeleteFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_node_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileRequest.ProtoReflect.Descriptor instead.
func (*DeleteFileRequest) Descriptor() ([]byte, []int) {
	return file_protos_files_node_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteFileRequest) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

type DeleteFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
}

func (x *DeleteFileResponse) Reset() {
	*x = DeleteFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_node_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFileResponse) ProtoMessage() {}

func (x *DeleteFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_node_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFileResponse.ProtoReflect.Descriptor instead.
func (*DeleteFileResponse) Descriptor() ([]byte, []int) {
	return file_protos_files_node_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteFileResponse) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

type SearchFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
}

func (x *SearchFileRequest) Reset() {
	*x = SearchFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_node_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchFileRequest) ProtoMessage() {}

func (x *SearchFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_node_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchFileRequest.ProtoReflect.Descriptor instead.
func (*SearchFileRequest) Descriptor() ([]byte, []int) {
	return file_protos_files_node_proto_rawDescGZIP(), []int{6}
}

func (x *SearchFileRequest) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

type SearchFileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileId string `protobuf:"bytes,1,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	Exists bool   `protobuf:"varint,2,opt,name=exists,proto3" json:"exists,omitempty"`
}

func (x *SearchFileResponse) Reset() {
	*x = SearchFileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_files_node_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchFileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchFileResponse) ProtoMessage() {}

func (x *SearchFileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_files_node_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchFileResponse.ProtoReflect.Descriptor instead.
func (*SearchFileResponse) Descriptor() ([]byte, []int) {
	return file_protos_files_node_proto_rawDescGZIP(), []int{7}
}

func (x *SearchFileResponse) GetFileId() string {
	if x != nil {
		return x.FileId
	}
	return ""
}

func (x *SearchFileResponse) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}

var File_protos_files_node_proto protoreflect.FileDescriptor

var file_protos_files_node_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x5f, 0x6e,
	0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x70, 0x32, 0x70, 0x5f, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x22, 0x4f, 0x0a, 0x11, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x5f,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x44, 0x61, 0x74, 0x61, 0x22, 0x7e, 0x0a, 0x12, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46,
	0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x66,
	0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x6c, 0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x2d, 0x0a, 0x13, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69,
	0x6c, 0x65, 0x49, 0x64, 0x22, 0x5d, 0x0a, 0x14, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64,
	0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x63, 0x68, 0x75,
	0x6e, 0x6b, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x65,
	0x78, 0x69, 0x73, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x78, 0x69,
	0x73, 0x74, 0x73, 0x22, 0x2b, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x65,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64,
	0x22, 0x2d, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22,
	0x2c, 0x0a, 0x11, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x45, 0x0a,
	0x12, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06,
	0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x65, 0x78,
	0x69, 0x73, 0x74, 0x73, 0x32, 0xc8, 0x02, 0x0a, 0x10, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x4e, 0x6f,
	0x64, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4b, 0x0a, 0x0a, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1c, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x12, 0x51, 0x0a, 0x0c, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f,
	0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1e, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x61, 0x67, 0x65,
	0x6e, 0x74, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x49, 0x0a, 0x0a, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1c, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x49, 0x0a, 0x0a, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x1c, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x53,
	0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x70, 0x32, 0x70, 0x5f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x73, 0x72, 0x63, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_protos_files_node_proto_rawDescOnce sync.Once
	file_protos_files_node_proto_rawDescData = file_protos_files_node_proto_rawDesc
)

func file_protos_files_node_proto_rawDescGZIP() []byte {
	file_protos_files_node_proto_rawDescOnce.Do(func() {
		file_protos_files_node_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_files_node_proto_rawDescData)
	})
	return file_protos_files_node_proto_rawDescData
}

var file_protos_files_node_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_protos_files_node_proto_goTypes = []any{
	(*UploadFileRequest)(nil),    // 0: p2p_agent.UploadFileRequest
	(*UploadFileResponse)(nil),   // 1: p2p_agent.UploadFileResponse
	(*DownloadFileRequest)(nil),  // 2: p2p_agent.DownloadFileRequest
	(*DownloadFileResponse)(nil), // 3: p2p_agent.DownloadFileResponse
	(*DeleteFileRequest)(nil),    // 4: p2p_agent.DeleteFileRequest
	(*DeleteFileResponse)(nil),   // 5: p2p_agent.DeleteFileResponse
	(*SearchFileRequest)(nil),    // 6: p2p_agent.SearchFileRequest
	(*SearchFileResponse)(nil),   // 7: p2p_agent.SearchFileResponse
}
var file_protos_files_node_proto_depIdxs = []int32{
	0, // 0: p2p_agent.FilesNodeService.UploadFile:input_type -> p2p_agent.UploadFileRequest
	2, // 1: p2p_agent.FilesNodeService.DownloadFile:input_type -> p2p_agent.DownloadFileRequest
	4, // 2: p2p_agent.FilesNodeService.DeleteFile:input_type -> p2p_agent.DeleteFileRequest
	6, // 3: p2p_agent.FilesNodeService.SearchFile:input_type -> p2p_agent.SearchFileRequest
	1, // 4: p2p_agent.FilesNodeService.UploadFile:output_type -> p2p_agent.UploadFileResponse
	3, // 5: p2p_agent.FilesNodeService.DownloadFile:output_type -> p2p_agent.DownloadFileResponse
	5, // 6: p2p_agent.FilesNodeService.DeleteFile:output_type -> p2p_agent.DeleteFileResponse
	7, // 7: p2p_agent.FilesNodeService.SearchFile:output_type -> p2p_agent.SearchFileResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_files_node_proto_init() }
func file_protos_files_node_proto_init() {
	if File_protos_files_node_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_files_node_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*UploadFileRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_node_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*UploadFileResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_node_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*DownloadFileRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_node_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*DownloadFileResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_node_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteFileRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_node_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteFileResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_node_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*SearchFileRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_files_node_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*SearchFileResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_files_node_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_files_node_proto_goTypes,
		DependencyIndexes: file_protos_files_node_proto_depIdxs,
		MessageInfos:      file_protos_files_node_proto_msgTypes,
	}.Build()
	File_protos_files_node_proto = out.File
	file_protos_files_node_proto_rawDesc = nil
	file_protos_files_node_proto_goTypes = nil
	file_protos_files_node_proto_depIdxs = nil
}
