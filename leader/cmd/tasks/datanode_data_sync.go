package tasks

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kysre/TurtleMQ/leader/internal/app/loadbalancer"
	"github.com/kysre/TurtleMQ/leader/internal/clients"
	"github.com/kysre/TurtleMQ/leader/internal/models"
	"github.com/kysre/TurtleMQ/leader/pkg/datanode"
)

type DataNodeDataSyncer struct {
	logger            *logrus.Logger
	dataNodeDirectory *models.DataNodeDirectory
	balancer          loadbalancer.Balancer
	partitionCount    int
	dataSyncTimeout   int
	shouldSync        bool
}

func NewDataNodeSyncer(
	logger *logrus.Logger,
	dataNodeDirectory *models.DataNodeDirectory,
	balancer loadbalancer.Balancer,
	partitionCount int,
	dataSyncTimeout int,
	shouldSync bool,
) *DataNodeDataSyncer {
	return &DataNodeDataSyncer{
		logger:            logger,
		dataNodeDirectory: dataNodeDirectory,
		balancer:          balancer,
		partitionCount:    partitionCount,
		dataSyncTimeout:   dataSyncTimeout,
		shouldSync:        shouldSync,
	}
}

// Data of datanode[i] is replicated in datanode[i+1]
// After datanode[i] fails:
//  - Remove datanode[i] from load balancer
//  - Replica data of datanode[i+1] should be read -> push to datanode[i+1] data
//  - Replica data of datanode[i+1] should be purged
//  - Data of datanode[i-1] should be read -> write to datanode[i+1] replica data

func (s *DataNodeDataSyncer) SyncData(failedDataNode *models.DataNode) {
	// Not sync data in leader's replica
	if !s.shouldSync {
		return
	}
	s.logger.Debug(fmt.Sprintf("Start datasync for datanode[%d]", failedDataNode.ID))
	// Create context for requests
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.dataSyncTimeout)*time.Second)
	defer cancel()
	// Get dataNodes that are present in data sync
	prevDataNode, afterDataNode, err := s.balancer.GetPreviousAndAfterDataNodesForSync(failedDataNode.Address)
	if err != nil {
		s.logger.Error(err)
	}
	// Halt push & pull to dataNodes that are present in data sync
	s.dataNodeDirectory.UpdateDataNodeState(failedDataNode.ID, models.DataNodeStatePENDING)
	s.dataNodeDirectory.UpdateDataNodeState(prevDataNode.ID, models.DataNodeStatePENDING)
	s.dataNodeDirectory.UpdateDataNodeState(afterDataNode.ID, models.DataNodeStatePENDING)
	// Step 1
	err = s.balancer.RemoveDataNodeFromHashCircle(failedDataNode.Address)
	if err != nil {
		logrus.Error(err)
		os.Exit(-1)
	}
	// Step 2
	for i := 0; i < s.partitionCount; i++ {
		res, err := afterDataNode.Client.ReadPartition(
			ctx, &datanode.ReadPartitionRequest{PartitionIndex: int32(i), IsReplica: true})
		if err != nil {
			s.logger.Error(err)
			continue
		}
		s.pushMessagesToDataNode(ctx, afterDataNode.Client, res.PartitionMessages)
	}
	// Step 3
	err = afterDataNode.Client.PurgeReplicaData(ctx)
	if err != nil {
		s.logger.Error(err)
	}
	// Step 4
	for i := 0; i < s.partitionCount; i++ {
		res, err := prevDataNode.Client.ReadPartition(
			ctx, &datanode.ReadPartitionRequest{PartitionIndex: int32(i), IsReplica: false})
		if err != nil {
			s.logger.Error(err)
			continue
		}
		req := datanode.WritePartitionRequest{PartitionIndex: int32(i), IsReplica: true}
		req.PartitionMessages = append(req.PartitionMessages, res.PartitionMessages...)
		err = afterDataNode.Client.WritePartition(ctx, &req)
		if err != nil {
			s.logger.Error(err)
		}
	}
	// Resume push & pull to dataNodes that are present in data sync
	s.dataNodeDirectory.UpdateDataNodeState(failedDataNode.ID, models.DataNodeStateUNHEALTHY)
	s.dataNodeDirectory.UpdateDataNodeState(prevDataNode.ID, models.DataNodeStateAVAILABLE)
	s.dataNodeDirectory.UpdateDataNodeState(afterDataNode.ID, models.DataNodeStateAVAILABLE)
	s.logger.Debug(fmt.Sprintf("Done datasync for datanode[%d]", failedDataNode.ID))
}

func (s *DataNodeDataSyncer) pushMessagesToDataNode(
	ctx context.Context, client clients.DataNodeClient, messages []*datanode.QueueMessage) {
	for _, message := range messages {
		req := datanode.PushRequest{IsReplica: false, Message: message}
		reqCtx, cancel := context.WithTimeout(ctx, time.Duration(10)*time.Second)
		_, err := client.Push(reqCtx, &req)
		if err != nil {
			s.logger.Error(err)
		}
		cancel()
	}
}
