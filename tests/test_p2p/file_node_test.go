package test_p2p

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tomp332/p2p-agent/src/node/p2p"
	"github.com/tomp332/p2p-agent/src/node/p2p/p2p_file_node"
	"github.com/tomp332/p2p-agent/src/pb"
	testUtils "github.com/tomp332/p2p-agent/tests/utils"
	"google.golang.org/protobuf/proto"
	"io"
	"testing"
)

func initializeFileNode(t *testing.T) (*p2p_file_node.FilesNode, *testUtils.TestGRPCable) {
	grpcTestServer := testUtils.NewTestAgentServer(t)
	config := testUtils.FileNodeConfig(t)
	baseNode := p2p.NewBaseNode(config)
	fileNode := p2p_file_node.NewP2PFilesNode(baseNode, &config.P2PFilesNodeConfig)
	fileNode.Register(grpcTestServer.BaseServer)
	err := grpcTestServer.Start()
	if err != nil {
		t.Fatalf("Failed to start node: %v", err)
	}
	return fileNode, grpcTestServer
}

func baseFileClient(t *testing.T) pb.FilesNodeServiceClient {
	_, server := initializeFileNode(t)
	conn, err := server.ClientConnection(testUtils.TestServerHostname)
	if err != nil {
		t.Fatalf("Failed to create client connection: %v", err)
	}
	return pb.NewFilesNodeServiceClient(conn)
}

func uploadFile(t *testing.T, client pb.FilesNodeServiceClient) (*pb.UploadFileResponse, [][]byte) {
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
	response, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("UploadFile(_) = _, %v", err)
	}
	if response == nil {
		t.Fatalf("Reply is nil")
	}
	expected := &pb.UploadFileResponse{
		FileId:   response.GetFileId(),
		FileSize: getExpectedFileSize(chunks),
		Success:  true,
		Message:  "File uploaded successfully",
	}

	if !proto.Equal(response, expected) {
		t.Errorf("UploadFile() = %v, want %v", response, expected)
	}
	t.Log("Uploaded test file to FilesNode")
	return response, chunks
}

func TestGroupNode(t *testing.T) {
	client := baseFileClient(t)
	uploadResponse, chunks := uploadFile(t, client)
	t.Run("DownloadFile", func(t *testing.T) {
		downloadStream, err := client.DownloadFile(context.Background(), &pb.DownloadFileRequest{FileId: uploadResponse.GetFileId()})
		if err != nil {
			t.Fatalf("DownloadFile(_) = _, %v", err)
		}

		var receivedChunks [][]byte
		for {
			chunk, downloadErr := downloadStream.Recv()
			if downloadErr == io.EOF {
				break
			}
			if downloadErr != nil {
				t.Fatalf("DownloadFile(_) = _, %v", downloadErr)
			}
			receivedChunks = append(receivedChunks, chunk.Chunk)
		}

		if !bytes.Equal(bytes.Join(receivedChunks, nil), bytes.Join(chunks, nil)) {
			t.Errorf("Downloaded chunks do not match uploaded chunks")
		}
	})
	t.Run("DownloadFileFromNetwork", func(t *testing.T) {
		n1, server := initializeFileNode(t)
		n1.BootstrapPeerAddrs = []string{testUtils.TestServerHostname}
		err := n1.ConnectToBootstrapPeers(server)
		if err != nil {
			t.Fatalf("Failed to connect to bootstrap peers: %v", err)
		}
		// Now, download the file
		downloadStream, err := client.DownloadFile(context.Background(), &pb.DownloadFileRequest{FileId: uploadResponse.GetFileId()})
		if err != nil {
			t.Fatalf("DownloadFile(_) = _, %v", err)
		}
		var receivedChunks [][]byte
		for {
			chunk, downloadErr := downloadStream.Recv()
			if downloadErr == io.EOF {
				break
			}
			if downloadErr != nil {
				t.Fatalf("DownloadFile(_) = _, %v", downloadErr)
			}
			receivedChunks = append(receivedChunks, chunk.Chunk)
		}
		if !bytes.Equal(bytes.Join(receivedChunks, nil), bytes.Join(chunks, nil)) {
			t.Errorf("Downloaded chunks do not match uploaded chunks")
		}
	})
	t.Run("SearchFile", func(t *testing.T) {
		searchResponse, err := client.SearchFile(context.Background(), &pb.SearchFileRequest{
			FileId: uploadResponse.GetFileId(),
		})
		if err != nil {
			t.Fatalf("SearchFile(_) = _, %v", err)
		}
		if searchResponse == nil {
			t.Fatalf("Response is nil")
		}
		if uploadResponse.FileId != searchResponse.GetFileId() || !searchResponse.Exists {
			t.Errorf("SearchFile(_) = %v, want %v", uploadResponse.FileId, uploadResponse.GetFileId())
		}
	})
	t.Run("DeleteFile", func(t *testing.T) {
		deleteResponse, err := client.DeleteFile(context.Background(), &pb.DeleteFileRequest{FileId: uploadResponse.GetFileId()})
		if err != nil {
			t.Fatalf("DeleteFile(_) = _, %v", err)
		}

		expected := &pb.DeleteFileResponse{FileId: uploadResponse.GetFileId()}

		if !proto.Equal(deleteResponse, expected) {
			t.Errorf("DeleteFile() = %v, want %v", deleteResponse, expected)
		}
		t.Log("Successfully deleted test file.")
	})
}

func TestValidBootstrapPeer(t *testing.T) {
	n, server := initializeFileNode(t)
	n.BootstrapPeerAddrs = []string{testUtils.TestServerHostname}
	err := n.ConnectToBootstrapPeers(server)
	if err != nil {
		t.Fatalf("Failed to connect to bootstrap peers: %v", err)
	}
	assert.Equal(t, 1, len(n.ConnectedPeers), "There should be 1 connected peer.")
}

func TestNonValidBootstrapPeer(t *testing.T) {
	n, server := initializeFileNode(t)
	n.BootstrapPeerAddrs = []string{"non-valid"}
	n.BootstrapNodeTimeout = 3
	err := n.ConnectToBootstrapPeers(server)
	assert.NotNil(t, err, "Invalid bootstrap peer")
}
