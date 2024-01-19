package core

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kysre/TurtleMQ/leader/internal/models"
	"github.com/kysre/TurtleMQ/leader/pkg/queue"
)

type queueCore struct {
	logger    *logrus.Logger
	directory *models.DataNodeDirectory
}

func NewQueueCore(logger *logrus.Logger, directory *models.DataNodeDirectory) queue.QueueServer {
	return &queueCore{
		logger:    logger,
		directory: directory,
	}
}

func (c *queueCore) Push(
	ctx context.Context, request *queue.PushRequest,
) (*emptypb.Empty, error) {
	// TODO: Get Related DN from LB and Push to it
	return &emptypb.Empty{}, nil
}

func (c *queueCore) Pull(
	ctx context.Context, request *emptypb.Empty,
) (*queue.PullResponse, error) {
	// TODO: Get DataNode from LB and call Pull
	return &queue.PullResponse{Key: "test", Value: make([][]byte, 0)}, nil
}

func (c *queueCore) AcknowledgePull(
	ctx context.Context, request *queue.AcknowledgePullRequest,
) (*emptypb.Empty, error) {
	// TODO: Get related DN from LB and call it's AcknowledgePull
	return &emptypb.Empty{}, nil
}
