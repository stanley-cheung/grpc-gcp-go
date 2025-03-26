package main

import (
	"fmt"
	"os"

	"cloud.google.com/go/storage/internal/apiv2/storagepb"

	"google.golang.org/grpc"
	grpcgoogle "google.golang.org/grpc/credentials/google"
)

func getGrpcClient() storagepb.StorageClient {
	var grpcOpts []grpc.DialOption
	endpoint = "dns:///storage.googleapis.com:443"
	grpcOpts = []grpc.DialOption{
		grpc.WithCredentialsBundle(
			grpcgoogle.NewComputeEngineCredentials(),
		),
		grpc.WithDisableServiceConfig(),
		grpc.WithDefaultServiceConfig(
			`{"loadBalancingConfig":[{"grpclb":{"childPolicy":[{"pick_first":{}}]}}]}`,
		),
	}
	conn, err := grpc.Dial(endpoint, grpcOpts...)
	if err != nil {
		fmt.Println("Failed to create clientconn: %v", err)
		os.Exit(1)
	}
	return storagepb.NewStorageClient(conn)
}

func main() {
	grpcClient := getGrpcClient()
}
