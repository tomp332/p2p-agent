package main

import (
	pb "github.com/tomp332/p2fs/src/protocol"
	"log"
	"sync"
)

var bootstrapAddresses = []string{
	"localhost:5001",
	"localhost:5002",
}

type Node struct {
	storage map[string][]byte
	peers   map[string]pb.FileSharingClient
	mu      sync.Mutex
}

func NewNode(port string, bootstrap bool) *Node {
	node := &Node{
		storage: make(map[string][]byte),
		peers:   make(map[string]pb.FileSharingClient),
	}

	// Start gRPC server
	go startGRPCServer(port, node.storage)

	// Connect to bootstrap nodes
	if bootstrap {
		for _, address := range bootstrapAddresses {
			client, err := connectToNode(address)
			if err != nil {
				log.Printf("failed to connect to bootstrap node %s: %v", address, err)
				continue
			}
			node.mu.Lock()
			node.peers[address] = client
			node.mu.Unlock()
		}
	}

	return node
}

func (node *Node) StoreChunk(chunk []byte) (string, error) {
	hash := hashChunk(chunk)
	node.mu.Lock()
	node.storage[hash] = chunk
	node.mu.Unlock()

	// Propagate to peers
	for address, client := range node.peers {
		go func(address string, client pb.FileSharingClient) {
			if _, err := storeChunk(client, chunk); err != nil {
				log.Printf("failed to store chunk on peer %s: %v", address, err)
			}
		}(address, client)
	}

	return hash, nil
}

func (node *Node) RetrieveChunk(hash string) ([]byte, bool) {
	node.mu.Lock()
	chunk, exists := node.storage[hash]
	node.mu.Unlock()

	if exists {
		return chunk, true
	}

	// Try to retrieve from peers
	for address, client := range node.peers {
		chunk, exists, err := retrieveChunk(client, hash)
		if err != nil {
			log.Printf("failed to retrieve chunk from peer %s: %v", address, err)
			continue
		}
		if exists {
			return chunk, true
		}
	}

	return nil, false
}

func main() {
	// Create node
	node := NewNode("5003", true)

	// Example: Store and retrieve a file chunk
	chunk := []byte("example data")
	hash, err := node.StoreChunk(chunk)
	if err != nil {
		log.Fatalf("failed to store chunk: %v", err)
	}
	log.Printf("stored chunk with hash: %s", hash)

	retrievedChunk, exists := node.RetrieveChunk(hash)
	if exists {
		log.Printf("retrieved chunk: %s", string(retrievedChunk))
	} else {
		log.Printf("chunk not found")
	}
}
