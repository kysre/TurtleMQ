package models

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kysre/TurtleMQ/leader/internal/clients"
	"github.com/kysre/TurtleMQ/leader/pkg/leader"
)

const MasterHost = "leader_0"

type LeaderConsensusHandler struct {
	replicaAddr   string
	replicaClient clients.LeaderClient
	isMaster      bool
}

func NewLeaderConsensusHandler(replicaHost string) *LeaderConsensusHandler {
	client, err := clients.NewLeaderClient(fmt.Sprintf("%s:8888", replicaHost))
	if err != nil {
		logrus.Error(err)
	}
	return &LeaderConsensusHandler{
		replicaAddr:   replicaHost,
		replicaClient: client,
		isMaster:      replicaHost != MasterHost,
	}
}

func (l *LeaderConsensusHandler) AmIMaster() bool {
	return l.isMaster
}

func (l *LeaderConsensusHandler) IsReplicaAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return l.replicaClient.IsHealthy(ctx)
}

func (l *LeaderConsensusHandler) AddDataNodesToReplica(ctx context.Context, dataNodes []*DataNode) {
	tickerPeriod := time.Duration(200) * time.Millisecond
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
