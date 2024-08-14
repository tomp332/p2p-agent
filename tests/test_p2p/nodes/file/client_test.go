package file

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tomp332/p2p-agent/src/node/p2p/file_node"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/tests/mocks"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestClient_SearchFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
	fileNodeClient := file_node.NewFileNodeClient(mockClient, 5*time.Second)

	mockClient.EXPECT().SearchFile(gomock.Any(), gomock.Any()).Return(&pb.SearchFileResponse{Exists: true}, nil).Times(1)

	exists, err := fileNodeClient.SearchFile(context.Background(), "fileId")
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestClient_DirectDownloadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
	mockStream := mocks.NewMockFilesNodeService_DirectDownloadFileClient(ctrl)
	fileNodeClient := file_node.NewFileNodeClient(mockClient, 5*time.Second)

	mockClient.EXPECT().DirectDownloadFile(gomock.Any(), gomock.Any()).Return(mockStream, nil).Times(1)

	mockStream.EXPECT().Recv().Return(&pb.DirectDownloadFileResponse{Chunk: []byte("chunk1")}, nil).Times(1)
	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	dataChan, errChan := fileNodeClient.DirectDownloadFile(context.Background(), "fileId", "nodeId")

	receivedData := <-dataChan
	assert.Equal(t, []byte("chunk1"), receivedData)
	assert.NoError(t, <-errChan)
}

func TestClient_UploadFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // Clean up the file afterward

	// Write some test data to the temporary file
	testData := []byte("this is a test file")
	_, err = tempFile.Write(testData)
	assert.NoError(t, err)
	tempFile.Close()

	// Mock the FilesNodeServiceClient and UploadFileClient
	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
	mockStream := mocks.NewMockFilesNodeService_UploadFileClient(ctrl)
	fileNodeClient := file_node.NewFileNodeClient(mockClient, 5*time.Second)

	// Expect the UploadFile method to return the mocked stream
	mockClient.EXPECT().UploadFile(gomock.Any()).Return(mockStream, nil).Times(1)

	// Expect the client to send the file data in chunks
	mockStream.EXPECT().Send(gomock.Any()).Return(nil).Times(1)

	// Expect the client to call CloseAndRecv and return a successful response
	mockStream.EXPECT().CloseAndRecv().Return(&pb.UploadFileResponse{
		Success:  true,
		FileId:   "fileId",
		FileSize: float64(len(testData)),
	}, nil).Times(1)

	// Run the UploadFile method with the temporary file
	errChan := fileNodeClient.UploadFile(context.Background(), tempFile.Name())

	// Assert that no errors occurred during upload
	assert.NoError(t, <-errChan)
}
