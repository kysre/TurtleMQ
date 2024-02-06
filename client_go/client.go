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

type SubscribeFunction func(key string, value []byte)

type QueueClient interface {
	Push(key string, value []byte)
	Pull() (string, []byte)
	Subscribe(function SubscribeFunction)
}

func GetQueueClient() QueueClient {
	leaderAddr := "localhost:8000"
	conn, err := grpc.Dial(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Can't connect to TurtleMQ cluster")
	}
	client := queue.NewQueueClient(conn)
	return &queueClient{
		client: client,
	}
}

type queueClient struct {
	client queue.QueueClient
}

func (c *queueClient) Push(key string, value []byte) {
	ctx := context.Background()
	req := queue.PushRequest{Key: key}
	req.Value = append(req.Value, value...)
	_, err := c.client.Push(ctx, &req)
	if err != nil {
		fmt.Print(err)
	}
}

func (c *queueClient) Pull() (string, []byte) {
	ctx := context.Background()
	res, err := c.client.Pull(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Print(err)
		return "", nil
	}
	return res.GetKey(), res.GetValue()
}

func (c *queueClient) Subscribe(function SubscribeFunction) {
	go c.runSubscribe(function)
}

func (c *queueClient) runSubscribe(function SubscribeFunction) {
	tickerPeriod := time.Duration(1) * time.Second
	ticker := time.NewTicker(tickerPeriod)
	for {
		select {
		case <-ticker.C:
			key, value := c.Pull()
			if value == nil {
				continue
			}
			function(key, value)
		}
	}
}
