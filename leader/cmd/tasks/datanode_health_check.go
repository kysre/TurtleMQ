package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kysre/TurtleMQ/leader/internal/models"
)

type DataNodeHealthChecker struct {
	dataNodeDirectory  *models.DataNodeDirectory
	healthCheckPeriod  int
	dataNodeDataSyncer *DataNodeDataSyncer
}

func NewDataNodeHealthChecker(
	dataNodeDirectory *models.DataNodeDirectory,
	healthCheckPeriod int,
	dataNodeDataSyncer *DataNodeDataSyncer,
) *DataNodeHealthChecker {
	return &DataNodeHealthChecker{
		dataNodeDirectory:  dataNodeDirectory,
		healthCheckPeriod:  healthCheckPeriod,
		dataNodeDataSyncer: dataNodeDataSyncer,
	}
}

func (hc *DataNodeHealthChecker) RunHealthChecks() {
	tickerPeriod := time.Duration(hc.healthCheckPeriod) * time.Second
	ticker := time.NewTicker(tickerPeriod)
	for {
		select {
		case <-ticker.C:
			logrus.Debug("Running DataNode health-check")
			hc.checkNodes()
		}
	}
}

func (hc *DataNodeHealthChecker) checkNodes() {
	dataNodes := hc.dataNodeDirectory.DataNodes
	for i, node := range dataNodes {
		if node.State != models.DataNodeStatePENDING {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			isHealthy := node.Client.IsHealthy(ctx)
			if node.State == models.DataNodeStateAVAILABLE && !isHealthy {
				hc.dataNodeDirectory.UpdateDataNodeState(i, models.DataNodeStateUNHEALTHY)
				logrus.Info(fmt.Sprintf("DataNode [%d] is not healthy!", i))
				go hc.dataNodeDataSyncer.SyncData(node)
			} else if node.State == models.DataNodeStateUNHEALTHY && isHealthy {
				hc.dataNodeDirectory.UpdateDataNodeState(i, models.DataNodeStateAVAILABLE)
				logrus.Info(fmt.Sprintf("DataNode [%d] became healthy!", i))
			}
			cancel()
		}
	}
}
