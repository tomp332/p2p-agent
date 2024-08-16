package file

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tomp332/p2p-agent/pkg/nodes/file_node"
	"github.com/tomp332/p2p-agent/pkg/pb"
	"github.com/tomp332/p2p-agent/tests/mocks"
	"io"
	"os"
	"testing"
	"time"
)

func TestFileNodeClient_SearchFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
	fileNodeClient := file_node.NewFileNodeClient(mockClient, 5*time.Second)

	// Mocking the gRPC call to SearchFile
	mockClient.EXPECT().SearchFile(gomock.Any(), &pb.SearchFileRequest{FileId: "fileId"}).Return(&pb.SearchFileResponse{Exists: true}, nil).Times(1)

	exists, err := fileNodeClient.SearchFile(context.Background(), "fileId")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestFileNodeClient_DirectDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
	mockStream := mocks.NewMockFilesNodeService_DirectDownloadFileClient(ctrl)
	fileNodeClient := file_node.NewFileNodeClient(mockClient, 5*time.Second)

	mockClient.EXPECT().DirectDownloadFile(gomock.Any(), &pb.DirectDownloadFileRequest{
		NodeId: "test-nodes-id",
		FileId: "fileId",
	}).Return(mockStream, nil).Times(1)

	mockStream.EXPECT().Recv().Return(&pb.DirectDownloadFileResponse{Chunk: []byte("chunk1")}, nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	dataChan, errChan := fileNodeClient.DirectDownloadFile(context.Background(), "fileId", "test-nodes-id")

	receivedData := <-dataChan
	assert.Equal(t, []byte("chunk1"), receivedData)
	assert.NoError(t, <-errChan)
}

func TestFileNodeClient_DownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
	mockStream := mocks.NewMockFilesNodeService_DownloadFileClient(ctrl)
	fileNodeClient := file_node.NewFileNodeClient(mockClient, 5*time.Second)

	mockClient.EXPECT().DownloadFile(gomock.Any(), &pb.DownloadFileRequest{FileId: "fileId"}).Return(mockStream, nil).Times(1)

	mockStream.EXPECT().Recv().Return(&pb.DownloadFileResponse{Chunk: []byte("chunk1")}, nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	dataChan, errChan := fileNodeClient.DownloadFile(context.Background(), "fileId")

	receivedData := <-dataChan
	assert.Equal(t, []byte("chunk1"), receivedData)
	assert.NoError(t, <-errChan)
}

func TestFileNodeClient_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up the file afterward

	// Write some test data to the temporary file
	testData := []byte("this is a test file")
	_, err = tempFile.Write(testData)
	assert.NoError(t, err)
	tempFile.Close()

	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
	mockStream := mocks.NewMockFilesNodeService_UploadFileClient(ctrl)
	fileNodeClient := file_node.NewFileNodeClient(mockClient, 5*time.Second)

	mockClient.EXPECT().UploadFile(gomock.Any()).Return(mockStream, nil).Times(1)

	mockStream.EXPECT().Send(gomock.Any()).Return(nil).Times(1)

	mockStream.EXPECT().CloseAndRecv().Return(&pb.UploadFileResponse{
		Success:  true,
		FileId:   "fileId",
		FileSize: float64(len(testData)),
	}, nil).Times(1)

	errChan := fileNodeClient.UploadFile(context.Background(), tempFile.Name())

	assert.NoError(t, <-errChan)
}
