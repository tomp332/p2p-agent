package test_p2p

import (
	"bytes"
	"context"
	"github.com/tomp332/p2p-agent/src"
	"github.com/tomp332/p2p-agent/src/node"
	"github.com/tomp332/p2p-agent/src/pb"
	testUtils "github.com/tomp332/p2p-agent/tests/utils"
	"google.golang.org/protobuf/proto"
	"io"
	"testing"
)

func baseFileClient(t *testing.T) pb.FilesNodeServiceClient {
	_, server := initializeFileNode(t)
	conn, err := server.ClientConnection("bufnet")
	if err != nil {
		t.Fatalf("Failed to create client connection: %v", err)
	}
	return pb.NewFilesNodeServiceClient(conn)
}

func initializeFileNode(t *testing.T) (src.P2PNoder, *testUtils.TestGRPCable) {
	grpcTestServer := testUtils.NewTestAgentServer(t)
	n, err := node.InitializeNode(grpcTestServer, testUtils.FileNodeConfig(t))
	if err != nil {
		t.Fatalf("Failed to initialize node: %v", err)
	}
	err = grpcTestServer.Start()
	if err != nil {
		t.Fatalf("Failed to start node: %v", err)
	}
	return n, grpcTestServer
}

func TestUploadFile(t *testing.T) {
	client := baseFileClient(t)
	stream, err := client.UploadFile(context.Background())
	if err != nil {
		t.Fatalf("UploadFile(_) = _, %v", err)
	}
	if stream == nil {
		t.Fatalf("Stream is nil")
	}

	chunks := [][]byte{
		[]byte("chunk1"),
		[]byte("chunk2"),
		[]byte("chunk3"),
	}

	for _, chunk := range chunks {
		if err := stream.Send(&pb.UploadFileRequest{ChunkData: chunk}); err != nil {
			t.Fatalf("Failed to send chunk: %v", err)
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("UploadFile(_) = _, %v", err)
	}
	if reply == nil {
		t.Fatalf("Reply is nil")
	}

	expected := &pb.UploadFileResponse{
		FileId:   reply.GetFileId(),
		FileSize: getExpectedFileSize(chunks),
		Success:  true,
		Message:  "File uploaded successfully",
	}

	if !proto.Equal(reply, expected) {
		t.Errorf("UploadFile() = %v, want %v", reply, expected)
	}
}

func TestDownloadFile(t *testing.T) {
	client := baseFileClient(t)
	uploadStream, err := client.UploadFile(context.Background())
	if err != nil {
		t.Fatalf("UploadFile(_) = _, %v", err)
	}

	chunks := [][]byte{
		[]byte("chunk1"),
		[]byte("chunk2"),
		[]byte("chunk3"),
	}

	for _, chunk := range chunks {
		if err := uploadStream.Send(&pb.UploadFileRequest{ChunkData: chunk}); err != nil {
			t.Fatalf("Failed to send chunk: %v", err)
		}
	}

	uploadResponse, err := uploadStream.CloseAndRecv()
	if err != nil {
		t.Fatalf("UploadFile(_) = _, %v", err)
	}

	// Now, download the file
	downloadStream, err := client.DownloadFile(context.Background(), &pb.DownloadFileRequest{FileId: uploadResponse.GetFileId()})
	if err != nil {
		t.Fatalf("DownloadFile(_) = _, %v", err)
	}

	var receivedChunks [][]byte
	for {
		chunk, err := downloadStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("DownloadFile(_) = _, %v", err)
		}
		receivedChunks = append(receivedChunks, chunk.Chunk)
	}

	if !bytes.Equal(bytes.Join(receivedChunks, nil), bytes.Join(chunks, nil)) {
		t.Errorf("Downloaded chunks do not match uploaded chunks")
	}
}

func TestDeleteFile(t *testing.T) {
	client := baseFileClient(t)
	uploadStream, err := client.UploadFile(context.Background())
	if err != nil {
		t.Fatalf("UploadFile(_) = _, %v", err)
	}

	chunks := [][]byte{
		[]byte("chunk1"),
		[]byte("chunk2"),
		[]byte("chunk3"),
	}

	for _, chunk := range chunks {
		if err := uploadStream.Send(&pb.UploadFileRequest{ChunkData: chunk}); err != nil {
			t.Fatalf("Failed to send chunk: %v", err)
		}
	}

	uploadResponse, err := uploadStream.CloseAndRecv()
	if err != nil {
		t.Fatalf("UploadFile(_) = _, %v", err)
	}

	// Now, delete the file
	deleteResponse, err := client.DeleteFile(context.Background(), &pb.DeleteFileRequest{FileId: uploadResponse.GetFileId()})
	if err != nil {
		t.Fatalf("DeleteFile(_) = _, %v", err)
	}

	expected := &pb.DeleteFileResponse{FileId: uploadResponse.GetFileId()}

	if !proto.Equal(deleteResponse, expected) {
		t.Errorf("DeleteFile() = %v, want %v", deleteResponse, expected)
	}
}

func TestFileNodeBootstrapPeers(t *testing.T) {
	n, _ := initializeFileNode(t)
	err := n.ConnectToBootstrapPeers()
	if err != nil {
		t.Fatalf("Failed to connect to bootstrap peers: %v", err)
	}
}

func getExpectedFileSize(chunks [][]byte) float64 {
	totalSize := 0
	for _, chunk := range chunks {
		totalSize += len(chunk)
	}
	return float64(totalSize)
}
