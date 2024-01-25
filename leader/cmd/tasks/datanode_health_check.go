package tasks

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kysre/TurtleMQ/leader/internal/models"
)

func RunHealthChecks(dataNodeDirectory *models.DataNodeDirectory, healthCheckPeriod int) {
	tickerPeriod := time.Duration(healthCheckPeriod) * time.Second
	ticker := time.NewTicker(tickerPeriod)
	for {
		select {
		case <-ticker.C:
			logrus.Info("Running DataNode health-check")
			checkNodes(dataNodeDirectory)
		}
	}
}

func checkNodes(dataNodeDirectory *models.DataNodeDirectory) {
	dataNodes := dataNodeDirectory.DataNodes
	for i, node := range dataNodes {
		if node.State == models.DataNodeStateUNHEALTHY {
			isHealthy := node.Client.IsHealthy()
			if isHealthy {
				dataNodeDirectory.UpdateDataNodeState(i, models.DataNodeStateAVAILABLE)
				logrus.Info(fmt.Sprintf("DataNode [%d] became healthy!", i))
			}
		}
	}
}