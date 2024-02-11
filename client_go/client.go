package client_go

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kysre/TurtleMQ/client_go/queue"
)

const HOST = "64.226.122.208"

type SubscribeFunction func(key string, value []byte)

type QueueClient interface {
	Push(key string, value []byte)
	Pull() (string, []byte)
	Subscribe(function SubscribeFunction)
}

func GetQueueClient() QueueClient {
	leaderAddr := fmt.Sprintf("%s:8000", HOST)
	leaderReplicaAddr := fmt.Sprintf("%s:8001", HOST)

	conn, err := grpc.Dial(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Can't connect to TurtleMQ cluster")
	}
	client := queue.NewQueueClient(conn)
	repConn, err := grpc.Dial(leaderReplicaAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Can't connect to TurtleMQ cluster")
	}
	replicaClient := queue.NewQueueClient(repConn)

	return &queueClient{
		client:        client,
		replicaClient: replicaClient,
	}
}

type queueClient struct {
	client        queue.QueueClient
	replicaClient queue.QueueClient
}

func (c *queueClient) Push(key string, value []byte) {
	ctx := context.Background()
	req := queue.PushRequest{Key: key}
	req.Value = append(req.Value, value...)
	_, err := c.client.Push(ctx, &req)
	if err != nil {
		_, _ = c.replicaClient.Push(ctx, &req)
	}
}

func (c *queueClient) Pull() (string, []byte) {
	ctx := context.Background()
	res, err := c.client.Pull(ctx, &emptypb.Empty{})
	if err != nil {
		res, err = c.replicaClient.Pull(ctx, &emptypb.Empty{})
		if err != nil {
			return "", nil
		}
	}
	defer c.acknowledgePull(ctx, res.GetKey())
	return res.GetKey(), res.GetValue()
}

func (c *queueClient) Subscribe(function SubscribeFunction) {
	go c.runSubscribe(function)
}

func (c *queueClient) acknowledgePull(ctx context.Context, key string) {
	req := queue.AcknowledgePullRequest{Key: key}
	_, err := c.client.AcknowledgePull(ctx, &req)
	if err != nil {
		_, _ = c.replicaClient.AcknowledgePull(ctx, &req)
	}
}

func (c *queueClient) runSubscribe(function SubscribeFunction) {
	tickerPeriod := time.Duration(1) * time.Second
	ticker := time.NewTicker(tickerPeriod)
	for {
		select {
		case <-ticker.C:
			key, value := c.Pull()
			if key == "" {
				continue
			}
			function(key, value)
		}
	}
}
