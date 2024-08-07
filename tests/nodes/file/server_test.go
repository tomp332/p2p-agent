package file

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tomp332/p2p-agent/pkg/nodes"
	"github.com/tomp332/p2p-agent/pkg/nodes/file_node"
	"github.com/tomp332/p2p-agent/pkg/pb"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"github.com/tomp332/p2p-agent/tests/mocks"
	"io"
	"testing"
	"time"
)

func setupFileNode(_ *testing.T, ctrl *gomock.Controller) (*file_node.FileNode, *mocks.MockStorage) {
	// Mock dependencies
	mockStorage := mocks.NewMockStorage(ctrl)

	// Create a basic NodeConfig for testing
	nodeConfig := &configs.NodeConfig{
		ID:   "test-nodes-id",
		Type: configs.FilesNodeType,
	}

	// Use the NewP2PFilesNode constructor to initialize the FileNode
	fileNode := file_node.NewP2PFilesNode(nodes.NewBaseNode(nodeConfig), mockStorage)

	return fileNode, mockStorage
}

func Test_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileNode, mockStorage := setupFileNode(t, ctrl)
	mockStream := mocks.NewMockFilesNodeService_UploadFileServer(ctrl)

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

	fileNode, mockStorage := setupFileNode(t, ctrl)
	mockStream := mocks.NewMockFilesNodeService_DownloadFileServer(ctrl)

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

	fileNode, mockStorage := setupFileNode(t, ctrl)

	mockStorage.EXPECT().Delete(gomock.Any(), "fileId").Return(nil).Times(1)

	response, err := fileNode.DeleteFile(context.Background(), &pb.DeleteFileRequest{FileId: "fileId"})
	assert.NoError(t, err)
	assert.Equal(t, "fileId", response.FileId)
}

func Test_SearchFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileNode, mockStorage := setupFileNode(t, ctrl)

	mockStorage.EXPECT().Search("fileId").Return(true).Times(1)

	response, err := fileNode.SearchFile(context.Background(), &pb.SearchFileRequest{FileId: "fileId"})
	assert.NoError(t, err)
	assert.True(t, response.Exists)
}

func Test_DirectDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileNode, mockStorage := setupFileNode(t, ctrl)
	mockStream := mocks.NewMockFilesNodeService_DirectDownloadFileServer(ctrl)

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
	err := fileNode.DirectDownloadFile(&pb.DirectDownloadFileRequest{FileId: "fileId", NodeId: "test-nodes-id"}, mockStream)
	assert.NoError(t, err)
}

func Test_SearchFileInNetwork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock the gRPC client
	mockGRPCClient := mocks.NewMockFilesNodeServiceClient(ctrl)
	mockStorage := mocks.NewMockStorage(ctrl)

	// Create a FileNodeClient using the mock gRPC client
	fileNodeClient := file_node.NewFileNodeClient(mockGRPCClient, 5*time.Second)

	// Create a FileNode and attach the FileNodeClient to a peer
	nodeConfig := &configs.NodeConfig{
		ID:                   "test-nodes-id",
		Type:                 configs.FilesNodeType,
		BootstrapPeerAddrs:   []string{"localhost:50051"},
		BootstrapNodeTimeout: 5,
	}
	fileNode := file_node.NewP2PFilesNode(nodes.NewBaseNode(nodeConfig), mockStorage)
	fileNode.ConnectedPeers = append(fileNode.ConnectedPeers, file_node.FileNodeConnection{
		NodeConnection: nodes.NodeConnection{
			Address:        "peer-1",
			GrpcConnection: nil,
		},
		NodeClient: *fileNodeClient, // Use the actual FileNodeClient instance
	})
	t.Run("TestFileFoundInNetwork", func(t *testing.T) {
		// Test case where the file is found on a peer nodes
		mockGRPCClient.EXPECT().SearchFile(gomock.Any(), &pb.SearchFileRequest{FileId: "fileId"}).Return(&pb.SearchFileResponse{Exists: true}, nil).Times(1)
		foundFileNode, err := fileNode.SearchFileInNetwork("fileId")
		assert.NoError(t, err)
		assert.Equal(t, "peer-1", foundFileNode.Address)
	})
	t.Run("TestFileNotFoundInNetwork", func(t *testing.T) {
		// Test case where the file is not found on any peer nodes
		mockGRPCClient.EXPECT().SearchFile(gomock.Any(), &pb.SearchFileRequest{FileId: "fileId"}).Return(&pb.SearchFileResponse{Exists: false}, nil).Times(1)
		_, errNotFound := fileNode.SearchFileInNetwork("fileId")
		if errNotFound == nil {
			t.Fatalf("File found when not supposed to.")
		}
		assert.Equal(t, "file was not found in network", errNotFound.Error())
	})
	t.Run("TestNoPeersConnected", func(t *testing.T) {
		// Test case where no peers are connected
		fileNodeNoPeers := file_node.NewP2PFilesNode(nodes.NewBaseNode(nodeConfig), mockStorage)
		_, errNoPeers := fileNodeNoPeers.SearchFileInNetwork("fileId")
		assert.Error(t, errNoPeers)
		assert.Equal(t, "no connected peers", errNoPeers.Error())
	})
	t.Run("TestPeerErrorInNetwork", func(t *testing.T) {
		// Test case where a peer returns an error during search
		mockGRPCClient.EXPECT().SearchFile(gomock.Any(), &pb.SearchFileRequest{FileId: "fileId"}).Return(nil, errors.New("search error")).Times(1)
		_, errPeerError := fileNode.SearchFileInNetwork("fileId")
		assert.Error(t, errPeerError)
		assert.Equal(t, "file was not found in network", errPeerError.Error())
	})
}
