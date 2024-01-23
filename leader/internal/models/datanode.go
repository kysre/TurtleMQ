package models

import (
	"sync"

	"github.com/kysre/TurtleMQ/leader/clients"
	"github.com/kysre/TurtleMQ/leader/pkg/errors"
)

type DataNodeState string

const (
	DataNodeStateUNHEALTHY DataNodeState = "UNHEALTHY"
	DataNodeStateAVAILABLE DataNodeState = "AVAILABLE"
	DataNodeStatePENDING   DataNodeState = "PENDING"
)

type DataNode struct {
	ID                int
	Address           string
	State             DataNodeState
	RemainingMsgCount int
	Client            clients.DataNodeClient
}

type DataNodeDirectory struct {
	DataNodes []*DataNode
	MX        *sync.Mutex
}

func NewDataNodeDirectory() *DataNodeDirectory {
	return &DataNodeDirectory{
		DataNodes: make([]*DataNode, 0),
		MX:        &sync.Mutex{},
	}
}

func (d *DataNodeDirectory) AddDataNode(dataNode *DataNode) error {
	d.MX.Lock()
	d.DataNodes = append(d.DataNodes, dataNode)
	d.MX.Unlock()
	return nil
}

func (d *DataNodeDirectory) GetDataNode(index int) (*DataNode, error) {
	d.MX.Lock()
	dataNode := d.DataNodes[index]
	healthy := dataNode.Client.IsHealthy()
	if !healthy {
		dataNode.State = DataNodeStateUNHEALTHY
	}
	d.MX.Unlock()

	if dataNode.State == DataNodeStatePENDING {
		return nil, errors.New("PENDING")
	}
	if dataNode.State == DataNodeStateUNHEALTHY {
		return nil, errors.New("PENDING")
	}
	return dataNode, nil
}

func (d *DataNodeDirectory) UpdateDataNodeState(index int, state DataNodeState) {
	d.MX.Lock()
	dataNode := d.DataNodes[index]
	dataNode.State = state
	d.MX.Unlock()
}
