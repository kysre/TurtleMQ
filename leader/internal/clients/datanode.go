package clients

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kysre/TurtleMQ/leader/pkg/datanode"
)

type DataNodeClient interface {
	IsHealthy(ctx context.Context) bool
	Push(ctx context.Context, request *datanode.PushRequest) (*empty.Empty, error)
	Pull(ctx context.Context, request *empty.Empty) (*datanode.PullResponse, error)
	AcknowledgePull(ctx context.Context, request *datanode.AcknowledgePullRequest) (*empty.Empty, error)
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

func (d *dataNodeClient) Push(ctx context.Context, request *datanode.PushRequest) (*empty.Empty, error) {
	res, err := d.client.Push(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *dataNodeClient) Pull(ctx context.Context, request *empty.Empty) (*datanode.PullResponse, error) {
	res, err := d.client.Pull(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *dataNodeClient) AcknowledgePull(
	ctx context.Context, request *datanode.AcknowledgePullRequest) (*empty.Empty, error) {
	res, err := d.client.AcknowledgePull(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
