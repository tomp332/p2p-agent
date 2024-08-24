package file

//func TestFileNodeClient_DownloadFile(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
//	mockStream := mocks.NewMockFilesNodeService_DownloadFileClient(ctrl)
//	fileNodeClient := fsNode.NewFileNodeClient(mockClient, nil, 5*time.Second, nil)
//
//	mockClient.EXPECT().DownloadFile(gomock.Any(), &pb.DownloadFileRequest{FileId: "fileId"}).Return(mockStream, nil).Times(1)
//
//	mockStream.EXPECT().Recv().Return(&pb.DownloadFileResponse{Chunk: []byte("chunk1")}, nil).Times(1)
//	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
//
//	dataChan, errChan := fileNodeClient.DownloadFile(context.Background(), "fileId")
//
//	receivedData := <-dataChan
//	assert.Equal(t, []byte("chunk1"), receivedData.ChunkData)
//	assert.NoError(t, <-errChan)
//}
//
//func TestFileNodeClient_Authenticate(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
//	fileNodeClient := fsNode.NewFileNodeClient(mockClient, nil, 5*time.Second, nil)
//	mockClient.EXPECT().Authenticate(gomock.Any(), &pb.AuthenticateRequest{Username: "test", Password: "test"}).Return(&pb.AuthenticateResponse{Token: "123"}, nil).Times(1)
//	authResponse, err := fileNodeClient.Authenticate(context.Background(), "test", "test")
//	assert.NoError(t, err)
//	assert.NotEmpty(t, authResponse.Token)
//	assert.Equal(t, "123", authResponse.Token)
//}
//
//func TestFileNodeClient_BadAuthentication(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	// Set up the mock client and the file node client
//	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
//	fileNodeClient := fsNode.NewFileNodeClient(mockClient, nil, 5*time.Second, nil)
//
//	// Set up the expectation for the Authenticate call to return an error
//	mockClient.EXPECT().
//		Authenticate(gomock.Any(), &pb.AuthenticateRequest{Username: "wrong", Password: "wrong"}).
//		Return(nil, errors.New("authentication failed")).
//		Times(1)
//
//	// Call the Authenticate method with incorrect credentials
//	authResponse, err := fileNodeClient.Authenticate(context.Background(), "wrong", "wrong")
//
//	// Assert that an error was returned
//	assert.Error(t, err)
//	assert.EqualError(t, err, "authentication failed")
//
//	// Assert that the response is nil (i.e., no token was returned)
//	assert.Nil(t, authResponse)
//}
//
//func TestFileNodeClient_UploadFile(t *testing.T) {
//	ctrl := gomock.NewController(t)
//
//	// Create a temporary file for testing
//	tempFile, err := os.CreateTemp("", "testfile")
//	assert.NoError(t, err)
//	defer os.Remove(tempFile.Name()) // Clean up the file afterward
//
//	// Write some test data to the temporary file
//	testData := []byte("this is a test file")
//	_, err = tempFile.Write(testData)
//	assert.NoError(t, err)
//	tempFile.Close()
//
//	mockClient := mocks.NewMockFilesNodeServiceClient(ctrl)
//	mockStream := mocks.NewMockFilesNodeService_UploadFileClient(ctrl)
//	fileNodeClient := fsNode.NewFileNodeClient(mockClient, nil, 5*time.Second, nil)
//
//	mockClient.EXPECT().UploadFile(gomock.Any()).Return(mockStream, nil).Times(1)
//
//	mockStream.EXPECT().Send(gomock.Any()).Return(nil).Times(1)
//
//	mockStream.EXPECT().CloseAndRecv().Return(&pb.UploadFileResponse{
//		Success:  true,
//		FileId:   "fileId",
//		FileSize: int64(len(testData)),
//	}, nil).Times(1)
//
//	errChan := fileNodeClient.UploadFile(context.Background(), tempFile.Name())
//
//	assert.NoError(t, <-errChan)
//}
