package core

import (
	"context"
	"fmt"
	"time"

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
	// Collect metrics
	startTime := time.Now()
	c.collectRpsMetrics("Push")
	defer c.collectLatencyMetrics("Push", startTime)

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
	// Collect metrics
	startTime := time.Now()
	c.collectRpsMetrics("Pull")
	defer c.collectLatencyMetrics("Pull", startTime)

	c.logger.Info("Received Pull request")
	client, err := c.balancer.GetPullDataNodeClient(ctx)
	if err != nil {
		return nil, err
	}
	c.logger.Info(fmt.Sprintf("Pull req datanode client: %v", client))
	dataNodeRes, err := client.Pull(ctx, request)
	if err != nil {
		return nil, err
	}
	message := dataNodeRes.GetMessage()
	c.logger.Info(fmt.Sprintf("Pull response with key: %s", message.GetKey()))
	response := queue.PullResponse{Key: message.GetKey()}
	response.Value = append(response.Value, message.GetValue()...)
	return &response, nil
}

func (c *queueCore) AcknowledgePull(
	ctx context.Context, request *queue.AcknowledgePullRequest,
) (*emptypb.Empty, error) {
	// Collect metrics
	startTime := time.Now()
	c.collectRpsMetrics("Ack")
	defer c.collectLatencyMetrics("Ack", startTime)

	key := request.GetKey()
	c.logger.Info(fmt.Sprintf("Received Ack Pull key=%s", key))
	client, err := c.balancer.GetPushDataNodeClient(ctx, key)
	if err != nil {
		return nil, err
	}
	dataNodeReq := datanode.AcknowledgePullRequest{Key: key}
	return client.AcknowledgePull(ctx, &dataNodeReq)
}

func (c *queueCore) collectLatencyMetrics(methodName string, startTime time.Time) {
	Latency.With(map[string]string{
		"provider": "leader",
		"method":   methodName,
		"status":   "Done",
	}).Observe(time.Since(startTime).Seconds())
}

func (c *queueCore) collectRpsMetrics(methodName string) {
	Total.With(map[string]string{
		"provider": "leader",
		"method":   methodName,
		"status":   "Done",
	}).Observe(1)
}
