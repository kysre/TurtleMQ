package core

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kysre/TurtleMQ/leader/pkg/queue"
)

type core struct {
}

func New() queue.QueueServer {
	return &core{}
}

func (c *core) Push(ctx context.Context, request *queue.PushRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (c *core) Pull(ctx context.Context, request *emptypb.Empty) (*queue.PullResponse, error) {
	return &queue.PullResponse{Key: "test", Value: make([]byte, 0)}, nil
}
