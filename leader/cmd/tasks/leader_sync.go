package tasks

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kysre/TurtleMQ/leader/internal/models"
)

type LeaderSyncer struct {
	logger    *logrus.Logger
	handler   *models.LeaderConsensusHandler
	directory *models.DataNodeDirectory
	period    int
}

func NewLeaderSyncer(
	logger *logrus.Logger,
	handler *models.LeaderConsensusHandler,
	directory *models.DataNodeDirectory,
	period int,
) *LeaderSyncer {
	return &LeaderSyncer{
		logger:    logger,
		handler:   handler,
		directory: directory,
		period:    period,
	}
}

func (ls *LeaderSyncer) RunLeaderSync() {
	ls.logger.Debug("Running leader sync")
	tickerPeriod := time.Duration(ls.period) * time.Second
	ticker := time.NewTicker(tickerPeriod)
	for {
		select {
		case <-ticker.C:
			if ls.handler.IsReplicaAvailable() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				ls.handler.AddDataNodesToReplica(ctx, ls.directory.DataNodes)
				cancel()
			}
		}
	}
}
