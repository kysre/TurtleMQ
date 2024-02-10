package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kysre/TurtleMQ/leader/internal/models"
)

func RunRemainingCheck(dataNodeDirectory *models.DataNodeDirectory, remainingCheckPeriod int) {
	tickerPeriod := time.Duration(remainingCheckPeriod) * time.Second
	ticker := time.NewTicker(tickerPeriod)
	for {
		select {
		case <-ticker.C:
			logrus.Info("Running DataNode remaining-check")
			updateNodesRemaining(dataNodeDirectory)
		}
	}
}

func updateNodesRemaining(dataNodeDirectory *models.DataNodeDirectory) {
	dataNodes := dataNodeDirectory.DataNodes
	for _, node := range dataNodes {
		if node.State == models.DataNodeStateAVAILABLE {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			count, err := node.Client.GetRemainingMessagesCount(ctx)
			if err != nil {
				logrus.Error(fmt.Sprintf("Error while remaining-check: %v", err))
			}
			node.RemainingMsgCount = int(count)
			cancel()
		}
	}
}
