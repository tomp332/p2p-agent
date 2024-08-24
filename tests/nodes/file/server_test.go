package file

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tomp332/p2p-agent/pkg/nodes"
	"github.com/tomp332/p2p-agent/pkg/nodes/fsNode"
	"github.com/tomp332/p2p-agent/pkg/pb"
	"github.com/tomp332/p2p-agent/pkg/utils/configs"
	"github.com/tomp332/p2p-agent/pkg/utils/types"
	"github.com/tomp332/p2p-agent/tests/mocks"
	"io"
	"sync"
	"testing"
)

func Test_UploadFile_Server(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Create mocks
	mockStorage := mocks.NewMockStorage(ctrl)
	clientStream := mocks.NewMockUploadFileServer(ctrl)

	nodeConfig := &configs.NodeConfig{
		Type: configs.FilesNodeType,
	}
	fileNode := fsNode.NewP2PFilesNode(nodes.NewBaseNode(nodeConfig), mockStorage)

	// Sub-test for successful upload with correct file size
	t.Run("Successful upload with correct file size", func(t *testing.T) {
		mockResult := &pb.UploadFileResponse{}

		clientStream.EXPECT().Context().Return(context.Background()).AnyTimes()
		callRecv1 := clientStream.EXPECT().Recv().
			Return(&pb.UploadFileRequest{FileName: "test.txt", ChunkData: []byte("chunk1")}, nil).
			Times(1)
		callRecv2 := clientStream.EXPECT().Recv().
			Return(&pb.UploadFileRequest{FileName: "test.txt", ChunkData: []byte("chunk2")}, nil).
			Times(1).After(callRecv1)
		clientStream.EXPECT().Recv().Return(nil, io.EOF).After(callRecv2)
		clientStream.EXPECT().SendAndClose(gomock.Any()).DoAndReturn(
			func(result *pb.UploadFileResponse) error {
				mockResult = result
				return nil
			})

		mockStorage.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fileName string, fileDataChan <-chan *types.TransferChunkData, wg *sync.WaitGroup) error {
				for range fileDataChan {
					// Simulate consuming the data
				}
				return nil
			},
		).Times(1)

		err := fileNode.UploadFile(clientStream)
		assert.Nil(t, err)
		assert.NotEmpty(t, mockResult)
		if mockResult.FileSize != 12 {
			t.Fatalf("wrong file size, expected 12, got: %d", mockResult.FileSize)
		}
	})

	// Sub-test for missing file name
	t.Run("Upload with missing file name", func(t *testing.T) {
		clientStream.EXPECT().Context().Return(context.Background()).AnyTimes()
		clientStream.EXPECT().Recv().
			Return(&pb.UploadFileRequest{FileName: "", ChunkData: []byte("chunk1")}, nil).
			Times(1)
		mockStorage.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fileName string, fileDataChan <-chan *types.TransferChunkData, wg *sync.WaitGroup) error {
				for range fileDataChan {
					// Simulate consuming the data
				}
				return nil
			},
		).Times(1)
		err := fileNode.UploadFile(clientStream)
		assert.NotNil(t, err)
		assert.Equal(t, "no file name was specified in the upload request", err.Error())
	})

	// Sub-test for SendAndClose failure
	t.Run("Upload with SendAndClose failure", func(t *testing.T) {
		clientStream.EXPECT().Context().Return(context.Background()).AnyTimes()
		callRecv1 := clientStream.EXPECT().Recv().
			Return(&pb.UploadFileRequest{FileName: "test.txt", ChunkData: []byte("chunk1")}, nil).
			Times(1)
		clientStream.EXPECT().Recv().Return(nil, io.EOF).After(callRecv1)
		clientStream.EXPECT().SendAndClose(gomock.Any()).DoAndReturn(
			func(result *pb.UploadFileResponse) error {
				return errors.New("SendAndClose failed")
			})

		mockStorage.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fileName string, fileDataChan <-chan *types.TransferChunkData, wg *sync.WaitGroup) error {
				for range fileDataChan {
					// Simulate consuming the data
				}
				return nil
			},
		).Times(1)

		err := fileNode.UploadFile(clientStream)
		assert.NotNil(t, err)
		assert.Equal(t, "SendAndClose failed", err.Error())
	})

	// Sub-test for error in Put
	t.Run("Upload with error in Put", func(t *testing.T) {

		clientStream.EXPECT().Context().Return(context.Background()).AnyTimes()
		callRecv1 := clientStream.EXPECT().Recv().
			Return(&pb.UploadFileRequest{FileName: "test.txt", ChunkData: []byte("chunk1")}, nil).
			Times(1)
		clientStream.EXPECT().Recv().Return(nil, io.EOF).After(callRecv1)
		mockStorage.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fileName string, fileDataChan <-chan *types.TransferChunkData, wg *sync.WaitGroup) error {
				for range fileDataChan {
					// Simulate consuming the data
				}
				return errors.New("storage error")
			},
		).Times(1)

		err := fileNode.UploadFile(clientStream)
		assert.NotNil(t, err)
		assert.Equal(t, "storage error", err.Error())
	})
}

func Test_DeleteFile_Server(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Create mocks
	mockStorage := mocks.NewMockStorage(ctrl)
	nodeConfig := &configs.NodeConfig{
		Type: configs.FilesNodeType,
	}
	fileNode := fsNode.NewP2PFilesNode(nodes.NewBaseNode(nodeConfig), mockStorage)

	t.Run("delete file successfully", func(t *testing.T) {
		mockStorage.EXPECT().Delete(gomock.Any(), "testfile").Return(nil)
		response, err := fileNode.DeleteFile(context.Background(), &pb.DeleteFileRequest{FileName: "testfile"})
		assert.NoError(t, err)
		assert.Equal(t, "testfile", response.FileName)
	})
	t.Run("error removing local file", func(t *testing.T) {
		mockStorage.EXPECT().Delete(gomock.Any(), "testfile").Return(errors.New("storage error"))
		_, err := fileNode.DeleteFile(context.Background(), &pb.DeleteFileRequest{FileName: "testfile"})
		assert.NotNil(t, err)
		assert.Equal(t, "storage error", err.Error())
	})
}

func Test_DownloadFile_Server(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStorage := mocks.NewMockStorage(ctrl)
	nodeConfig := &configs.NodeConfig{
		Type: configs.FilesNodeType,
	}
	fileNode := fsNode.NewP2PFilesNode(nodes.NewBaseNode(nodeConfig), mockStorage)

	t.Run("download file successfully", func(t *testing.T) {
		configs.MainConfig.ID = "1"
		serverStream := mocks.NewMockDownloadFileServer(ctrl)
		serverStream.EXPECT().Context().Return(context.Background()).AnyTimes()
		mockStorage.EXPECT().Exists("testfile").Return(true).AnyTimes()
		serverStream.EXPECT().Send(&pb.DownloadFileResponse{
			FileName:  "testfile",
			Chunk:     []byte("chunk1"),
			ChunkSize: 6,
		}).Return(nil)

		mockStorage.
			EXPECT().
			Get(context.Background(), "testfile").
			DoAndReturn(func(ctx context.Context, fileName string) (<-chan *types.TransferChunkData, error) {
				dataChan := make(chan *types.TransferChunkData, 1)
				defer close(dataChan)
				dataChan <- &types.TransferChunkData{
					ChunkSize: 6,
					ChunkData: []byte("chunk1"),
				}
				return dataChan, nil
			})
		err := fileNode.DownloadFile(&pb.DownloadFileRequest{
			FileName: "testfile",
		}, serverStream)
		assert.NoError(t, err)
	})

	t.Run("file does not exist", func(t *testing.T) {
		serverStream := mocks.NewMockDownloadFileServer(ctrl)
		serverStream.EXPECT().Context().Return(context.Background()).AnyTimes()
		mockStorage.EXPECT().Exists("nonexistentfile").Return(false).AnyTimes()

		err := fileNode.DownloadFile(&pb.DownloadFileRequest{
			FileName: "nonexistentfile",
		}, serverStream)
		assert.Error(t, err)
		assert.Equal(t, "file was not found locally, and there are no connected peers", err.Error())
	})

	t.Run("error retrieving file chunks", func(t *testing.T) {
		serverStream := mocks.NewMockDownloadFileServer(ctrl)
		serverStream.EXPECT().Context().Return(context.Background()).AnyTimes()
		mockStorage.EXPECT().Exists("testfile").Return(true).AnyTimes()
		mockStorage.EXPECT().Get(context.Background(), "testfile").Return(nil, errors.New("error getting file chunks"))

		err := fileNode.DownloadFile(&pb.DownloadFileRequest{
			FileName: "testfile",
		}, serverStream)
		assert.Error(t, err)
		assert.Equal(t, "error getting file chunks", err.Error())
	})

	t.Run("error sending file chunks", func(t *testing.T) {
		serverStream := mocks.NewMockDownloadFileServer(ctrl)
		serverStream.EXPECT().Context().Return(context.Background()).AnyTimes()
		mockStorage.EXPECT().Exists("testfile").Return(true).AnyTimes()
		serverStream.EXPECT().Send(&pb.DownloadFileResponse{
			FileName:  "testfile",
			Chunk:     []byte("chunk1"),
			ChunkSize: 6,
		}).Return(errors.New("error sending chunk"))

		mockStorage.
			EXPECT().
			Get(context.Background(), "testfile").
			DoAndReturn(func(ctx context.Context, fileName string) (<-chan *types.TransferChunkData, error) {
				dataChan := make(chan *types.TransferChunkData, 1)
				defer close(dataChan)
				dataChan <- &types.TransferChunkData{
					ChunkSize: 6,
					ChunkData: []byte("chunk1"),
				}
				return dataChan, nil
			})
		err := fileNode.DownloadFile(&pb.DownloadFileRequest{
			FileName: "testfile",
		}, serverStream)
		assert.Error(t, err)
		assert.Equal(t, "error sending chunk", err.Error())
	})

	t.Run("context canceled during download", func(t *testing.T) {
		// Create a cancellable context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Mock the server stream
		serverStream := mocks.NewMockDownloadFileServer(ctrl)
		serverStream.EXPECT().Context().Return(ctx).AnyTimes()

		// Mock the storage to simulate file existence
		mockStorage.EXPECT().Exists("testfile").Return(true).AnyTimes()

		// Mock the storage's Get method to simulate data retrieval and cancellation
		mockStorage.
			EXPECT().
			Get(ctx, "testfile").
			DoAndReturn(func(ctx context.Context, fileName string) (<-chan *types.TransferChunkData, error) {
				dataChan := make(chan *types.TransferChunkData, 1)
				go func() {
					// Simulate sending a chunk after some delay
					dataChan <- &types.TransferChunkData{
						ChunkData: []byte("data"),
						ChunkSize: 4,
					}

					// Simulate context cancellation after a short delay
					cancel()

					// Close the data channel
					close(dataChan)
				}()
				return dataChan, nil
			})

		// Call DownloadFile and expect it to handle the context cancellation
		err := fileNode.DownloadFile(&pb.DownloadFileRequest{
			FileName: "testfile",
		}, serverStream)

		// Verify that an error occurred due to context cancellation
		assert.Error(t, err)
		assert.Equal(t, context.Canceled.Error(), err.Error())
	})

	t.Run("multiple chunks to send", func(t *testing.T) {
		serverStream := mocks.NewMockDownloadFileServer(ctrl)
		serverStream.EXPECT().Context().Return(context.Background()).AnyTimes()
		mockStorage.EXPECT().Exists("testfile").Return(true).AnyTimes()
		serverStream.EXPECT().Send(&pb.DownloadFileResponse{
			FileName:  "testfile",
			Chunk:     []byte("chunk1"),
			ChunkSize: 6,
		}).Return(nil)
		serverStream.EXPECT().Send(&pb.DownloadFileResponse{
			FileName:  "testfile",
			Chunk:     []byte("chunk2"),
			ChunkSize: 6,
		}).Return(nil)

		mockStorage.
			EXPECT().
			Get(context.Background(), "testfile").
			DoAndReturn(func(ctx context.Context, fileName string) (<-chan *types.TransferChunkData, error) {
				dataChan := make(chan *types.TransferChunkData, 2)
				defer close(dataChan)
				dataChan <- &types.TransferChunkData{
					ChunkSize: 6,
					ChunkData: []byte("chunk1"),
				}
				dataChan <- &types.TransferChunkData{
					ChunkSize: 6,
					ChunkData: []byte("chunk2"),
				}
				return dataChan, nil
			})
		err := fileNode.DownloadFile(&pb.DownloadFileRequest{
			FileName: "testfile",
		}, serverStream)
		assert.NoError(t, err)
	})
}

//func TestFileNode_Authenticate(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	mockAuthManager := mocks.NewMockAuthenticationManager(ctrl)
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("valid"), bcrypt.DefaultCost)
//	fileNode := &fsNode.FileNode{
//		BaseNode: &nodes.BaseNode{
//			NodeConfig: configs.NodeConfig{
//				Type: configs.FilesNodeType,
//				Auth: configs.NodeAuthConfig{
//					Username: "valid",
//					Password: string(hashedPassword),
//				},
//			},
//			AuthManager: mockAuthManager,
//		},
//	}
//
//	t.Run("Successful Authentication", func(t *testing.T) {
//		mockAuthManager.EXPECT().Generate("valid", fileNode.Type).Return("valid-token", nil).Times(1)
//
//		req := &pb.AuthenticateRequest{Username: "valid", Password: "valid"}
//		res, err := fileNode.Authenticate(context.Background(), req)
//
//		assert.NoError(t, err)
//		assert.NotNil(t, res)
//		assert.Equal(t, "valid-token", res.Token)
//	})
//
//	t.Run("Invalid Username", func(t *testing.T) {
//		req := &pb.AuthenticateRequest{Username: "invalid", Password: "valid"}
//		res, err := fileNode.Authenticate(context.Background(), req)
//
//		assert.Error(t, err)
//		assert.Nil(t, res)
//		assert.EqualError(t, err, "invalid username/password")
//	})
//
//	t.Run("Invalid Password", func(t *testing.T) {
//		req := &pb.AuthenticateRequest{Username: "valid", Password: "invalid"}
//		res, err := fileNode.Authenticate(context.Background(), req)
//
//		assert.Error(t, err)
//		assert.Nil(t, res)
//		assert.EqualError(t, err, "invalid username/password")
//	})
//
//	t.Run("Generate Token Failure", func(t *testing.T) {
//		mockAuthManager.EXPECT().Generate("valid", fileNode.Type).Return("", fmt.Errorf("failed to generate token")).Times(1)
//
//		req := &pb.AuthenticateRequest{Username: "valid", Password: "valid"}
//		res, err := fileNode.Authenticate(context.Background(), req)
//
//		assert.Error(t, err)
//		assert.Nil(t, res)
//		assert.EqualError(t, err, "failed to generate token")
//	})
//}

//func TestFileNode_DownloadFile_RemotePeer(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	// Create mocks
//	remoteMockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
//	mockStorage := mocks.NewMockStorage(ctrl)
//	mockStream := mocks.NewMockFilesNodeService_DownloadFileClient(ctrl) // Mock stream
//	// Create FileNode instance with mocks
//	fileNode := &fsNode.FileNode{
//		BaseNode: &nodes.BaseNode{
//			ConnectedPeers: map[string]nodes.PeerConnection{
//				"1": {
//					ID: "1",
//					ConnectionInfo: &configs.BootStrapNodeConnection{
//						Host: "test",
//						Port: 1111,
//					},
//					NodeClient: fsNode.NewFileNodeClient(remoteMockClient, nil, 5*time.Minute, nil),
//				},
//			},
//		},
//		Storage: mockStorage,
//	}
//	t.Run("TestRemoteFile", func(t *testing.T) {
//		// Set up mock response for Storage.Get method
//		mockStorage.EXPECT().
//			Get(gomock.Any(), "testfile").
//			Return(nil, fmt.Errorf("no file found")) // Simulate file not found locally
//		mockStorage.EXPECT().
//			Exists("testfile").
//			Return(false)
//
//		// Mock the client to return the mock stream
//		remoteMockClient.EXPECT().
//			DownloadFile(gomock.Any(), gomock.Any()).
//			Return(mockStream, nil)
//		remoteMockClient.EXPECT().SearchFile(gomock.Any(), &pb.SearchFileRequest{
//			FileId: "testfile",
//		}).Return(&pb.SearchFileResponse{
//			FileId: "testfile",
//			NodeId: "1",
//		}, nil)
//		// Setup expectations for the mock stream
//		mockStream.EXPECT().Recv().Return(&pb.DownloadFileResponse{Chunk: []byte("chunk1")}, nil).Times(1)
//		mockStream.EXPECT().Recv().Return(&pb.DownloadFileResponse{Chunk: []byte("chunk2")}, nil).Times(1)
//		mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
//
//		// Create a mock DownloadFileServer to act as the remote peer
//		mockDownloadFileServer := mocks.NewMockFilesNodeService_DownloadFileServer(ctrl)
//
//		// Set up the mock DownloadFileServer to expect Context and Send calls
//		mockDownloadFileServer.EXPECT().
//			Context().
//			Return(context.Background()).
//			AnyTimes()
//
//		// Setup mock to expect Send calls for each chunk
//		mockDownloadFileServer.EXPECT().
//			Send(&pb.DownloadFileResponse{
//				FileId: "testfile",
//				Exists: true,
//				Chunk:  []byte("chunk1"),
//			}).Return(nil).Times(1)
//
//		mockDownloadFileServer.EXPECT().
//			Send(&pb.DownloadFileResponse{
//				FileId: "testfile",
//				Exists: true,
//				Chunk:  []byte("chunk2"),
//			}).Return(nil).Times(1)
//
//		// Perform the test
//		err := fileNode.DownloadFile(&pb.DownloadFileRequest{FileId: "testfile"}, mockDownloadFileServer)
//		assert.NoError(t, err)
//	})
//	t.Run("DownloadFileLocal", func(t *testing.T) {
//		// Set the file as existing
//		mockStorage.EXPECT().
//			Exists("testfile").
//			Return(true)
//		mocks.NewMockFilesNodeService_Sear
//		mockStorage.EXPECT().
//			Get(gomock.Any(), "testfile").
//			Return(nil, fmt.Errorf("no file found")) // Simulate file not found locally
//	})
//}
