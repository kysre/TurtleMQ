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
	GetRemainingMessagesCount(ctx context.Context) (int32, error)

	Push(ctx context.Context, request *datanode.PushRequest) (*empty.Empty, error)
	Pull(ctx context.Context, request *empty.Empty) (*datanode.PullResponse, error)
	AcknowledgePull(ctx context.Context, request *datanode.AcknowledgePullRequest) (*empty.Empty, error)

	ReadPartition(ctx context.Context, request *datanode.ReadPartitionRequest) (*datanode.ReadPartitionResponse, error)
	WritePartition(ctx context.Context, request *datanode.WritePartitionRequest) error
	PurgeReplicaData(ctx context.Context) error
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

func (d *dataNodeClient) GetRemainingMessagesCount(ctx context.Context) (int32, error) {
	res, err := d.client.GetRemainingMessagesCount(ctx, &emptypb.Empty{})
	if err != nil {
		return 0, err
	}
	return res.GetRemainingMessagesCount(), nil
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

func (d *dataNodeClient) ReadPartition(
	ctx context.Context, request *datanode.ReadPartitionRequest,
) (*datanode.ReadPartitionResponse, error) {
	res, err := d.client.ReadPartition(ctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *dataNodeClient) WritePartition(
	ctx context.Context, request *datanode.WritePartitionRequest,
) error {
	_, err := d.client.WritePartition(ctx, request)
	return err
}

func (d *dataNodeClient) PurgeReplicaData(ctx context.Context) error {
	_, err := d.client.PurgeReplicaData(ctx, &emptypb.Empty{})
	return err
}
