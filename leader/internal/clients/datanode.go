package clients

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kysre/TurtleMQ/leader/pkg/datanode"
)

type DataNodeClient interface {
	IsHealthy(ctx context.Context) bool
}

type dataNodeClient struct {
	client datanode.DataNodeClient
}

func NewDataNodeClient(addr string) (DataNodeClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := datanode.NewDataNodeClient(conn)
	return &dataNodeClient{
		client: client,
	}, nil
}

func (d *dataNodeClient) IsHealthy(ctx context.Context) bool {
	_, err := d.client.IsHealthy(ctx, &emptypb.Empty{})
	return err == nil
}
