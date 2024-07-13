package main

import (
	"context"
	"time"

	pb "github.com/tomp332/p2fs/src/protocol"

	"google.golang.org/grpc"
)

func storeChunk(client pb.FileSharingClient, chunk []byte) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.StoreChunk(ctx, &pb.StoreChunkRequest{Chunk: chunk})
	if err != nil {
		return "", err
	}
	return res.GetHash(), nil
}

func retrieveChunk(client pb.FileSharingClient, hash string) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.RetrieveChunk(ctx, &pb.RetrieveChunkRequest{Hash: hash})
	if err != nil {
		return nil, false, err
	}
	return res.GetChunk(), res.GetExists(), nil
}

func connectToNode(address string) (pb.FileSharingClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := pb.NewFileSharingClient(conn)
	return client, nil
}
