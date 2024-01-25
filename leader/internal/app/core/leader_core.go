package core

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kysre/TurtleMQ/leader/internal/app/loadbalancer"
	"github.com/kysre/TurtleMQ/leader/internal/clients"
	"github.com/kysre/TurtleMQ/leader/internal/models"
	"github.com/kysre/TurtleMQ/leader/pkg/leader"
)

type leaderCore struct {
	logger    *logrus.Logger
	directory *models.DataNodeDirectory

	balancer loadbalancer.Balancer
}

func NewLeaderCore(logger *logrus.Logger, directory *models.DataNodeDirectory) leader.LeaderServer {
	return &leaderCore{
		logger:    logger,
		directory: directory,
	}
}

func (lc *leaderCore) IsHealthy(ctx context.Context, request *empty.Empty) (*empty.Empty, error) {
	return &emptypb.Empty{}, nil // Will not return anything if not healthy
}

func (lc *leaderCore) AddDataNode(ctx context.Context, request *leader.AddDataNodeRequest) (*empty.Empty, error) {
	address := request.GetAddress()
	client, err := clients.NewDataNodeClient(address)
	if err != nil {
		return nil, err
	}
	dataNode := models.DataNode{
		ID:                lc.directory.GetDataNodeCount(),
		Address:           address,
		State:             models.DataNodeStateAVAILABLE,
		RemainingMsgCount: 0,
		Client:            client,
	}
	err = lc.directory.AddDataNode(&dataNode)
	if err != nil {
		return nil, err
	}
	err = lc.balancer.AddDataNodeToHashCircle(&dataNode)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
