package loadbalancer

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/kysre/TurtleMQ/leader/internal/clients"
	"github.com/kysre/TurtleMQ/leader/internal/models"
	"github.com/kysre/TurtleMQ/leader/pkg/errors"
)

type Balancer interface {
	AddDataNodeToHashCircle(datanode *models.DataNode) error
	GetPushDataNodeAndReplicaClient(
		ctx context.Context, token string) (clients.DataNodeClient, clients.DataNodeClient, error)
	GetPullDataNodeClient(ctx context.Context) (clients.DataNodeClient, error)

	GetPreviousAndAfterDataNodesForSync(addr string) (*models.DataNode, *models.DataNode, error)
	RemoveDataNodeFromHashCircle(addr string) error
}

type balancer struct {
	logger    *logrus.Logger
	directory *models.DataNodeDirectory

	datanodeHashSortedSlice []string
	dataNodeHashMap         map[string]int
}

func NewBalancer(logger *logrus.Logger, directory *models.DataNodeDirectory) Balancer {
	return &balancer{
		logger:    logger,
		directory: directory,

		datanodeHashSortedSlice: make([]string, 0),
		dataNodeHashMap:         make(map[string]int),
	}
}

func (b *balancer) GetPushDataNodeAndReplicaClient(
	ctx context.Context, token string) (clients.DataNodeClient, clients.DataNodeClient, error) {
	hash, err := b.getHash(token)
	if err != nil {
		return nil, nil, err
	}
	// Calc datanode & it's replica indexes
	index, err := b.getLessOrEqualIndexInHashCircle(hash)
	if err != nil {
		return nil, nil, err
	}
	// Get datanode
	dataNodeHash := b.datanodeHashSortedSlice[index]
	dataNode := b.directory.GetDataNode(b.dataNodeHashMap[dataNodeHash])
	if dataNode.State == models.DataNodeStatePENDING {
		return nil, nil, err
	}
	if dataNode.State == models.DataNodeStateUNHEALTHY {
		index = index + 1
		dataNodeHash = b.datanodeHashSortedSlice[index]
		dataNode = b.directory.GetDataNode(b.dataNodeHashMap[dataNodeHash])
	}
	// Get datanode replica
	replicaIndex := b.getDataNodeReplicaIndex(index)
	dataNodeReplicaHash := b.datanodeHashSortedSlice[replicaIndex]
	dataNodeReplica := b.directory.GetDataNode(b.dataNodeHashMap[dataNodeReplicaHash])
	if dataNodeReplica.State == models.DataNodeStatePENDING {
		return nil, nil, err
	}
	if dataNodeReplica.State == models.DataNodeStateUNHEALTHY {
		replicaIndex = replicaIndex + 1
		dataNodeReplicaHash = b.datanodeHashSortedSlice[replicaIndex]
		dataNodeReplica = b.directory.GetDataNode(b.dataNodeHashMap[dataNodeReplicaHash])
	}
	return dataNode.Client, dataNodeReplica.Client, nil
}

// TODO: Add randomness to loadbalancer for pull

func (b *balancer) GetPullDataNodeClient(ctx context.Context) (clients.DataNodeClient, error) {
	maxRemainingMsgHash := ""
	maxRemainingMsgCount := 0
	for i := 0; i < 2*len(b.datanodeHashSortedSlice); i++ {
		dnHash := b.datanodeHashSortedSlice[i%len(b.datanodeHashSortedSlice)]
		dn := b.directory.GetDataNode(b.dataNodeHashMap[dnHash])
		if dn.State != models.DataNodeStateAVAILABLE {
			continue
		}
		if dn.RemainingMsgCount > maxRemainingMsgCount {
			maxRemainingMsgHash = dnHash
			maxRemainingMsgCount = dn.RemainingMsgCount
		}
	}
	if maxRemainingMsgHash == "" {
		return nil, errors.New("No AVAILABLE datanode!")
	}
	dn := b.directory.GetDataNode(b.dataNodeHashMap[maxRemainingMsgHash])
	return dn.Client, nil
}

func (b *balancer) AddDataNodeToHashCircle(datanode *models.DataNode) error {
	datanodeHash, err := b.getHash(datanode.Address)
	if err != nil {
		return err
	}
	insertionIndex, err := b.getLessOrEqualIndexInHashCircle(datanodeHash)
	b.dataNodeHashMap[datanodeHash] = datanode.ID
	b.datanodeHashSortedSlice = b.insertInSlice(b.datanodeHashSortedSlice, insertionIndex, datanodeHash)
	b.logger.Info(fmt.Sprintf("New Hash Circle SortedSlice: %v", b.datanodeHashSortedSlice))
	b.logger.Info(fmt.Sprintf("New Hash Circle HashMap: %v", b.dataNodeHashMap))
	return nil
}

func (b *balancer) getHash(value string) (string, error) {
	hash := md5.New()
	_, err := io.WriteString(hash, value)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (b *balancer) isHashLessOrEqual(hash, otherHash string) (bool, error) {
	decodedHash, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}
	decodedOtherHash, err := hex.DecodeString(otherHash)
	if err != nil {
		return false, err
	}
	return bytes.Compare(decodedHash, decodedOtherHash) <= 0, nil
}

func (b *balancer) insertInSlice(s []string, index int, value string) []string {
	if len(s) == index { // nil or empty slice or after last element
		return append(s, value)
	}
	s = append(s[:index+1], s[index:]...) // index < len(s)
	s[index] = value
	return s
}

func (b *balancer) removeFromSlice(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func (b *balancer) getLessOrEqualIndexInHashCircle(hash string) (int, error) {
	index := 0
	for i, sortedDnHash := range b.datanodeHashSortedSlice {
		isLessOrEqual, err := b.isHashLessOrEqual(hash, sortedDnHash)
		if err != nil {
			return 0, err
		}
		if isLessOrEqual {
			index = i
			break
		}
	}
	return index, nil
}

// Hash Ring implementation (Datanode[i]'s data is replicated in Datanode[i+1])
func (b *balancer) getDataNodeReplicaIndex(i int) int {
	return (i + 1) % len(b.datanodeHashSortedSlice)
}

func (b *balancer) GetPreviousAndAfterDataNodesForSync(addr string) (*models.DataNode, *models.DataNode, error) {
	datanodeHash, err := b.getHash(addr)
	if err != nil {
		return nil, nil, err
	}
	index, err := b.getLessOrEqualIndexInHashCircle(datanodeHash)
	if err != nil {
		return nil, nil, err
	}
	prevIndex := (index - 1 + len(b.datanodeHashSortedSlice)) % len(b.datanodeHashSortedSlice)
	prevDataNodeHash := b.datanodeHashSortedSlice[prevIndex]
	afterIndex := (index + 1) % len(b.datanodeHashSortedSlice)
	afterDataNodeHash := b.datanodeHashSortedSlice[afterIndex]
	prevDataNode := b.directory.GetDataNode(b.dataNodeHashMap[prevDataNodeHash])
	afterDataNode := b.directory.GetDataNode(b.dataNodeHashMap[afterDataNodeHash])
	return prevDataNode, afterDataNode, nil
}

func (b *balancer) RemoveDataNodeFromHashCircle(addr string) error {
	datanodeHash, err := b.getHash(addr)
	if err != nil {
		return err
	}
	index, err := b.getLessOrEqualIndexInHashCircle(datanodeHash)
	if err != nil {
		return err
	}
	b.datanodeHashSortedSlice = b.removeFromSlice(b.datanodeHashSortedSlice, index)
	delete(b.dataNodeHashMap, datanodeHash)
	return nil
}
