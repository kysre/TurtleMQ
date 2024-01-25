package loadbalancer

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"

	"github.com/sirupsen/logrus"

	"github.com/kysre/TurtleMQ/leader/internal/clients"
	"github.com/kysre/TurtleMQ/leader/internal/models"
)

type Balancer interface {
	AddDataNodeToHashCircle(datanode *models.DataNode) error
	GetPushDataNodeClient(ctx context.Context, token string) (*clients.DataNodeClient, error)
	GetPullDataNodeClient(ctx context.Context) (*clients.DataNodeClient, error)
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

func (b *balancer) GetPushDataNodeClient(ctx context.Context, token string) (*clients.DataNodeClient, error) {
	hash, err := b.getHash(token)
	if err != nil {
		return nil, err
	}
	index, err := b.getLessOrEqualIndexInHashCircle(hash)
	if err != nil {
		return nil, err
	}
	dnHash := b.datanodeHashSortedSlice[index]
	dn, err := b.directory.GetDataNode(ctx, b.dataNodeHashMap[dnHash])
	if err != nil {
		return nil, err
	}
	return &dn.Client, nil
}

func (b *balancer) GetPullDataNodeClient(ctx context.Context) (*clients.DataNodeClient, error) {
	maxRemainingMsgHash := ""
	maxRemainingMsgCount := 0
	for i := 0; i < 2*len(b.datanodeHashSortedSlice); i++ {
		dnHash := b.datanodeHashSortedSlice[i%len(b.datanodeHashSortedSlice)]
		dn, _ := b.directory.GetDataNode(ctx, b.dataNodeHashMap[dnHash])
		if dn == nil {
			continue
		}
		if dn.RemainingMsgCount > maxRemainingMsgCount {
			maxRemainingMsgHash = dnHash
			maxRemainingMsgCount = dn.RemainingMsgCount
		}
	}
	dn, err := b.directory.GetDataNode(ctx, b.dataNodeHashMap[maxRemainingMsgHash])
	if err != nil {
		return nil, err
	}
	return &dn.Client, nil
}

func (b *balancer) AddDataNodeToHashCircle(datanode *models.DataNode) error {
	datanodeHash, err := b.getHash(datanode.Address)
	if err != nil {
		return err
	}
	insertionIndex, err := b.getLessOrEqualIndexInHashCircle(datanodeHash)
	b.dataNodeHashMap[datanodeHash] = datanode.ID
	b.datanodeHashSortedSlice = b.insertInSlice(b.datanodeHashSortedSlice, insertionIndex, datanodeHash)
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
