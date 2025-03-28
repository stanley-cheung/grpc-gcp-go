package main

import (
	"context"
	"fmt"
	"hash/crc32"
	"os"

	gcspb "github.com/stanley-cheung/grpc-gcp-go/e2e-examples/gcs/cloud.google.com/go/storage/genproto/apiv2/storagepb"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/balancer/rls"
	"google.golang.org/grpc/credentials/google"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	_ "google.golang.org/grpc/xds/googledirectpath"
)

func getGrpcClient() gcspb.StorageClient {
	resolver.SetDefaultScheme("dns")
	endpoint := "google-c2p:///storage.googleapis.com"

	var grpcOpts []grpc.DialOption
	grpcOpts = []grpc.DialOption{
		grpc.WithCredentialsBundle(
			google.NewComputeEngineCredentials(),
		),
	}
	conn, err := grpc.Dial(endpoint, grpcOpts...)
	if err != nil {
		fmt.Println("Failed to create clientconn: %v", err)
		os.Exit(1)
	}
	return gcspb.NewStorageClient(conn)
}

func readRequest(client gcspb.StorageClient) {
	ctx := context.Background()
	req := gcspb.ReadObjectRequest{
		Bucket: "projects/_/buckets/stanleycheung-bucket",
		Object: "test01.txt",
	}
	ctx = metadata.AppendToOutgoingContext(ctx,
		"x-goog-request-params", "bucket=projects/_/buckets/stanleycheung-bucket")
	stream, err := client.ReadObject(ctx, &req)
	if err != nil {
		fmt.Println("ReadObject got error: ", err)
		os.Exit(1)
	}
	resp, err := stream.Recv()
	if err != nil {
		fmt.Println("ReadObject Recv error: ", err)
		os.Exit(1)
	}
	fmt.Println("ReadObject result: ", resp.ChecksummedData.String())
}

func writeRequest(client gcspb.StorageClient) {
	data := []byte("test test test\ntest test test \n")
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := crc32.Checksum(data, crc32c)
	fmt.Println("CRC32C checksum: ", checksum)
}

func main() {
	grpcClient := getGrpcClient()
	readRequest(grpcClient)
	writeRequest(grpcClient)
}
