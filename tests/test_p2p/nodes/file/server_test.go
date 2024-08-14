package file

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/node/p2p/file_node"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/utils/configs"
	mocks2 "github.com/tomp332/p2p-agent/tests/mocks"
	"io"
	"testing"
)

func setupFileNode(_ *testing.T, ctrl *gomock.Controller) (*file_node.FileNode, *mocks2.MockStorage, *mocks2.MockAgentGRPCServer) {
	// Mock dependencies
	mockStorage := mocks2.NewMockStorage(ctrl)
	mockServer := mocks2.NewMockAgentGRPCServer(ctrl)

	// Create a basic NodeConfig for testing
	nodeConfig := &configs.NodeConfig{
		ID:   "test-node-id",
		Type: configs.FilesNodeType,
	}

	// Use the NewP2PFilesNode constructor to initialize the FileNode
	fileNode := file_node.NewP2PFilesNode(p2p.NewBaseNode(nodeConfig), mockStorage)

	return fileNode, mockStorage, mockServer
}

func Test_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileNode, mockStorage, _ := setupFileNode(t, ctrl)
	mockStream := mocks2.NewMockFilesNodeService_UploadFileServer(ctrl)

	// Create a buffer with some data to simulate the file being uploaded
	buffer := []byte("filedata")

	// Set up the stream to return chunks of data
	mockStream.EXPECT().Recv().Return(&pb.UploadFileRequest{ChunkData: buffer}, nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	mockStream.EXPECT().Context().Return(context.Background()).Times(1)
	mockStorage.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(float64(len(buffer)), nil).Times(1)
	mockStream.EXPECT().SendAndClose(gomock.Any()).Return(nil).Times(1)

	err := fileNode.UploadFile(mockStream)
	assert.NoError(t, err)

}

func Test_DownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileNode, mockStorage, _ := setupFileNode(t, ctrl)
	mockStream := mocks2.NewMockFilesNodeService_DownloadFileServer(ctrl)

	// Set up the storage to return a data stream
	dataChan := make(chan []byte, 1)
	dataChan <- []byte("chunk1")
	close(dataChan)

	mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dataChan, nil).Times(1)
	mockStream.EXPECT().Context().Return(context.Background()).Times(1)
	mockStream.EXPECT().Send(gomock.Any()).Return(nil).Times(1)

	err := fileNode.DownloadFile(&pb.DownloadFileRequest{FileId: "fileId"}, mockStream)
	assert.NoError(t, err)
}

func Test_DeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileNode, mockStorage, _ := setupFileNode(t, ctrl)

	mockStorage.EXPECT().Delete(gomock.Any(), "fileId").Return(nil).Times(1)

	response, err := fileNode.DeleteFile(context.Background(), &pb.DeleteFileRequest{FileId: "fileId"})
	assert.NoError(t, err)
	assert.Equal(t, "fileId", response.FileId)
}

func Test_SearchFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileNode, mockStorage, _ := setupFileNode(t, ctrl)

	mockStorage.EXPECT().Search("fileId").Return(true).Times(1)

	response, err := fileNode.SearchFile(context.Background(), &pb.SearchFileRequest{FileId: "fileId"})
	assert.NoError(t, err)
	assert.True(t, response.Exists)
}

func Test_DirectDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileNode, mockStorage, _ := setupFileNode(t, ctrl)
	mockStream := mocks2.NewMockFilesNodeService_DirectDownloadFileServer(ctrl)

	// Set up the storage to return a data stream
	dataChan := make(chan []byte, 2)
	dataChan <- []byte("chunk1")
	dataChan <- []byte("chunk2")
	close(dataChan)

	// Mock the storage Get method to return the data channel
	mockStorage.EXPECT().Get(gomock.Any(), "fileId").Return(dataChan, nil).Times(1)
	mockStream.EXPECT().Context().Return(context.Background()).Times(1)

	// Expect the server to send the chunks back via the stream
	mockStream.EXPECT().Send(&pb.DirectDownloadFileResponse{
		FileId: "fileId",
		Exists: true,
		Chunk:  []byte("chunk1"),
	}).Return(nil).Times(1)
	mockStream.EXPECT().Send(&pb.DirectDownloadFileResponse{
		FileId: "fileId",
		Exists: true,
		Chunk:  []byte("chunk2"),
	}).Return(nil).Times(1)

	// Run the DirectDownloadFile method
	err := fileNode.DirectDownloadFile(&pb.DirectDownloadFileRequest{FileId: "fileId", NodeId: "test-node-id"}, mockStream)
	assert.NoError(t, err)
}
