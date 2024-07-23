package p2p_test

import (
	"bytes"
	"context"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/pb"
	"github.com/tomp332/p2p-agent/src/utils"
	testUtils "github.com/tomp332/p2p-agent/tests/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
	"testing"
)

var (
	client pb.FilesNodeServiceClient
	conn   *grpc.ClientConn
)

func TestMain(m *testing.M) {
	testUtils.SetupConfig()
	utils.SetupLogger()
	utils.Logger.Info().Msg("Starting p2p_test file system tests.")
	// Create temporary storage
	localStorage, cleanupStorage := testUtils.CreateTempStorage(nil)
	defer cleanupStorage()

	// Initialize FilesNode with temporary storage
	filesNode := &p2p.FilesNode{
		Storage: localStorage,
	}

	// Start gRPC server
	var cleanupServer func()
	conn, cleanupServer = testUtils.StartTestGRPCServer(nil, filesNode)
	defer cleanupServer()

	// Create gRPC client
	client = pb.NewFilesNodeServiceClient(conn)

	// Run tests
	code := m.Run()

	// Cleanup
	conn.Close()

	os.Exit(code)
}

func TestUploadFile(t *testing.T) {
	ctx := context.Background()
	stream, err := client.UploadFile(ctx)
	if err != nil {
		t.Fatalf("UploadFile(_) = _, %v", err)
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

	expected := &pb.UploadFileResponse{
		FileId:   reply.GetFileId(),
		FileSize: 18,
		Success:  true,
		Message:  "File uploaded successfully",
	}

	if !proto.Equal(reply, expected) {
		t.Errorf("UploadFile() = %v, want %v", reply, expected)
	}
}

func TestDownloadFile(t *testing.T) {
	ctx := context.Background()

	// First, upload a file to download later
	uploadStream, err := client.UploadFile(ctx)
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
	downloadStream, err := client.DownloadFile(ctx, &pb.DownloadFileRequest{FileId: uploadResponse.GetFileId()})
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
	ctx := context.Background()

	// First, upload a file to delete later
	uploadStream, err := client.UploadFile(ctx)
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
	deleteResponse, err := client.DeleteFile(ctx, &pb.DeleteFileRequest{FileId: uploadResponse.GetFileId()})
	if err != nil {
		t.Fatalf("DeleteFile(_) = _, %v", err)
	}

	expected := &pb.DeleteFileResponse{FileId: uploadResponse.GetFileId()}

	if !proto.Equal(deleteResponse, expected) {
		t.Errorf("DeleteFile() = %v, want %v", deleteResponse, expected)
	}
}
