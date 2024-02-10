package clients

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kysre/TurtleMQ/leader/pkg/leader"
)

type LeaderClient interface {
	IsHealthy(ctx context.Context) bool
	AddDataNode(ctx context.Context, request *leader.AddDataNodeRequest) error
}

type leaderClient struct {
	client leader.LeaderClient
}

func NewLeaderClient(addr string) (LeaderClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := leader.NewLeaderClient(conn)
	return &leaderClient{
		client: client,
	}, nil
}

func (c *leaderClient) IsHealthy(ctx context.Context) bool {
	_, err := c.client.IsHealthy(ctx, &emptypb.Empty{})
	return err == nil
}

func (c *leaderClient) AddDataNode(ctx context.Context, request *leader.AddDataNodeRequest) error {
	_, err := c.client.AddDataNode(ctx, request)
	return err
}
