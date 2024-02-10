package models

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kysre/TurtleMQ/leader/internal/clients"
	"github.com/kysre/TurtleMQ/leader/pkg/leader"
)

const MasterAddr = "leader_0"

type LeaderConsensusHandler struct {
	replicaAddr   string
	replicaClient clients.LeaderClient
	isMaster      bool
}

func NewLeaderConsensusHandler(addr string) *LeaderConsensusHandler {
	client, err := clients.NewLeaderClient(addr)
	if err != nil {
		logrus.Error(err)
	}
	return &LeaderConsensusHandler{
		replicaAddr:   addr,
		replicaClient: client,
		isMaster:      addr == MasterAddr,
	}
}

func (l *LeaderConsensusHandler) AmIMaster() bool {
	return l.isMaster
}

func (l *LeaderConsensusHandler) AddDataNodesToReplica(ctx context.Context, dataNodes []*DataNode) {
	tickerPeriod := time.Duration(500) * time.Millisecond
	ticker := time.NewTicker(tickerPeriod)
	for _, node := range dataNodes {
		select {
		case <-ticker.C:
			err := l.replicaClient.AddDataNode(ctx, &leader.AddDataNodeRequest{Address: node.Address})
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
