package core

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kysre/TurtleMQ/leader/internal/app/loadbalancer"
	"github.com/kysre/TurtleMQ/leader/internal/models"
	"github.com/kysre/TurtleMQ/leader/pkg/datanode"
	"github.com/kysre/TurtleMQ/leader/pkg/queue"
)

type queueCore struct {
	logger    *logrus.Logger
	directory *models.DataNodeDirectory

	balancer loadbalancer.Balancer
}

func NewQueueCore(logger *logrus.Logger, directory *models.DataNodeDirectory, balancer loadbalancer.Balancer) queue.QueueServer {
	return &queueCore{
		logger:    logger,
		directory: directory,
		balancer:  balancer,
	}
}

func (c *queueCore) Push(
	ctx context.Context, request *queue.PushRequest,
) (*emptypb.Empty, error) {
	key := request.GetKey()
	client, err := c.balancer.GetPushDataNodeClient(ctx, key)
	c.logger.Info(fmt.Sprintf("Push key: %s to DataNode %v", key, client))
	if err != nil {
		return nil, err
	}
	messagePb := datanode.QueueMessage{Key: key}
	messagePb.Value = append(messagePb.Value, request.GetValue()...)
	dataNodeReq := datanode.PushRequest{Message: &messagePb}
	return client.Push(ctx, &dataNodeReq)
}

func (c *queueCore) Pull(
	ctx context.Context, request *emptypb.Empty,
) (*queue.PullResponse, error) {
	client, err := c.balancer.GetPullDataNodeClient(ctx)
	if err != nil {
		return nil, err
	}
	dataNodeRes, err := client.Pull(ctx, request)
	if err != nil {
		return nil, err
	}
	message := dataNodeRes.GetMessage()
	response := queue.PullResponse{Key: message.GetKey()}
	response.Value = append(response.Value, message.GetValue()...)
	return &response, nil
}

func (c *queueCore) AcknowledgePull(
	ctx context.Context, request *queue.AcknowledgePullRequest,
) (*emptypb.Empty, error) {
	key := request.GetKey()
	client, err := c.balancer.GetPushDataNodeClient(ctx, key)
	if err != nil {
		return nil, err
	}
	dataNodeReq := datanode.AcknowledgePullRequest{Key: key}
	return client.AcknowledgePull(ctx, &dataNodeReq)
}
