package file

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tomp332/p2p-agent/pkg/nodes"
	"github.com/tomp332/p2p-agent/pkg/nodes/fsNode"
	"github.com/tomp332/p2p-agent/pkg/pb"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"github.com/tomp332/p2p-agent/tests/mocks"
	"golang.org/x/crypto/bcrypt"
	"io"
	"testing"
	"time"
)

func setupFileNode(t *testing.T, ctrl *gomock.Controller) (*fsNode.FileNode, *mocks.MockStorage, *mocks.MockFilesNodeServiceClient) {
	// Mock dependencies
	mockStorage := mocks.NewMockStorage(ctrl)
	mockGRPCClient := mocks.NewMockFilesNodeServiceClient(ctrl)

	// Create a basic NodeConfig for testing
	nodeConfig := &configs.NodeConfig{
		Type: configs.FilesNodeType,
	}

	// Use the NewP2PFilesNode constructor to initialize the FileNode
	fileNode := fsNode.NewP2PFilesNode(nodes.NewBaseNode(nodeConfig), mockStorage)
	return fileNode, mockStorage, mockGRPCClient
}
func Test_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)

	fileNode, mockStorage, _ := setupFileNode(t, ctrl)
	mockStream := mocks.NewMockFilesNodeService_UploadFileServer(ctrl)

	// Create a buffer with some data to simulate the file being uploaded
	buffer := []byte("filedata")

	// Set up the stream to return chunks of data
	mockStream.EXPECT().Recv().Return(&pb.UploadFileRequest{ChunkData: buffer}, nil).Times(1)

	mockStream.EXPECT().Context().Return(context.Background()).Times(1)
	mockStorage.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(int64(len(buffer)), nil).Times(1)
	mockStream.EXPECT().SendAndClose(gomock.Any()).Return(nil).Times(1)

	err := fileNode.UploadFile(mockStream)
	assert.NoError(t, err)
}

func TestFileNode_DownloadFile_RemotePeer(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Create mocks
	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
	mockStorage := mocks.NewMockStorage(ctrl)
	mockStream := mocks.NewMockFilesNodeService_DownloadFileClient(ctrl) // Mock stream

	// Set up mock response for Storage.Get method
	mockStorage.EXPECT().
		Get(gomock.Any(), "testfile").
		Return(nil, fmt.Errorf("no file found")) // Simulate file not found locally

	// Mock the client to return the mock stream
	mockClient.EXPECT().
		DownloadFile(gomock.Any(), gomock.Any()).
		Return(mockStream, nil)

	// Setup expectations for the mock stream
	mockStream.EXPECT().Recv().Return(&pb.DownloadFileResponse{Chunk: []byte("chunk1")}, nil).Times(1)
	mockStream.EXPECT().Recv().Return(&pb.DownloadFileResponse{Chunk: []byte("chunk2")}, nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	// Create FileNode instance with mocks
	fileNode := &fsNode.FileNode{
		BaseNode: &nodes.BaseNode{
			ConnectedPeers: []nodes.PeerConnection{
				{
					ConnectionInfo: &configs.BootStrapNodeConnection{
						Host: "test",
						Port: 1111,
					},
					NodeClient: fsNode.NewFileNodeClient(mockClient, 5*time.Minute),
				},
			},
		},
		Storage: mockStorage,
	}

	// Create a mock DownloadFileServer
	mockDownloadFileServer := mocks.NewMockFilesNodeService_DownloadFileServer(ctrl)

	// Set up the mock DownloadFileServer to expect Context and Send calls
	mockDownloadFileServer.EXPECT().
		Context().
		Return(context.Background()).
		Times(1)

	// Setup mock to expect Send calls for each chunk
	mockDownloadFileServer.EXPECT().
		Send(&pb.DownloadFileResponse{
			FileId: "testfile",
			Exists: true,
			Chunk:  []byte("chunk1"),
		}).Return(nil).Times(1)

	mockDownloadFileServer.EXPECT().
		Send(&pb.DownloadFileResponse{
			FileId: "testfile",
			Exists: true,
			Chunk:  []byte("chunk2"),
		}).Return(nil).Times(1)

	// Perform the test
	err := fileNode.DownloadFile(&pb.DownloadFileRequest{FileId: "testfile"}, mockDownloadFileServer)
	assert.NoError(t, err)
}

func Test_DeleteFile(t *testing.T) {
	ctrl := gomock.NewController(t)

	fileNode, mockStorage, _ := setupFileNode(t, ctrl)

	mockStorage.EXPECT().Delete(gomock.Any(), "fileId").Return(nil).Times(1)

	response, err := fileNode.DeleteFile(context.Background(), &pb.DeleteFileRequest{FileId: "fileId"})
	assert.NoError(t, err)
	assert.Equal(t, "fileId", response.FileId)
}

func TestFileNode_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockAuthManager := mocks.NewMockAuthenticationManager(ctrl)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("valid"), bcrypt.DefaultCost)
	fileNode := &fsNode.FileNode{
		BaseNode: &nodes.BaseNode{
			NodeConfig: configs.NodeConfig{
				Type: configs.FilesNodeType,
				Auth: configs.NodeAuthConfig{
					Username: "valid",
					Password: string(hashedPassword),
				},
			},
			AuthManager: mockAuthManager,
		},
	}

	t.Run("Successful Authentication", func(t *testing.T) {
		mockAuthManager.EXPECT().Generate("valid", fileNode.Type).Return("valid-token", nil).Times(1)

		req := &pb.AuthenticateRequest{Username: "valid", Password: "valid"}
		res, err := fileNode.Authenticate(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "valid-token", res.Token)
	})

	t.Run("Invalid Username", func(t *testing.T) {
		req := &pb.AuthenticateRequest{Username: "invalid", Password: "valid"}
		res, err := fileNode.Authenticate(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "invalid username/password")
	})

	t.Run("Invalid Password", func(t *testing.T) {
		req := &pb.AuthenticateRequest{Username: "valid", Password: "invalid"}
		res, err := fileNode.Authenticate(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "invalid username/password")
	})

	t.Run("Generate Token Failure", func(t *testing.T) {
		mockAuthManager.EXPECT().Generate("valid", fileNode.Type).Return("", fmt.Errorf("failed to generate token")).Times(1)

		req := &pb.AuthenticateRequest{Username: "valid", Password: "valid"}
		res, err := fileNode.Authenticate(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, "failed to generate token")
	})
}
