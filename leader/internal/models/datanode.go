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
	id                int
	address           string
	state             DataNodeState
	remainingMsgCount int
	client            clients.DataNodeClient
}

type DataNodeDirectory struct {
	dataNodes []*DataNode
	mx        *sync.Mutex
}

func NewDataNodeDirectory() *DataNodeDirectory {
	return &DataNodeDirectory{
		dataNodes: make([]*DataNode, 0),
		mx:        &sync.Mutex{},
	}
}

func (d *DataNodeDirectory) AddDataNode(dataNode *DataNode) error {
	d.mx.Lock()
	d.dataNodes = append(d.dataNodes, dataNode)
	d.mx.Unlock()
	return nil
}

func (d *DataNodeDirectory) GetDataNode(index int) (*DataNode, error) {
	d.mx.Lock()
	dataNode := d.dataNodes[index]
	healthy := dataNode.client.IsHealthy()
	if !healthy {
		dataNode.state = DataNodeStateUNHEALTHY
	}
	d.mx.Unlock()

	if dataNode.state == DataNodeStatePENDING {
		return nil, errors.New("PENDING")
	}
	if dataNode.state == DataNodeStateUNHEALTHY {
		return nil, errors.New("PENDING")
	}
	return dataNode, nil
}

func (d *DataNodeDirectory) UpdateDataNodeState(index int, state DataNodeState) {
	d.mx.Lock()
	dataNode := d.dataNodes[index]
	dataNode.state = state
	d.mx.Unlock()
}
