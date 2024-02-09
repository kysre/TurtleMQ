package models

import (
	"sync"

	"github.com/kysre/TurtleMQ/leader/internal/clients"
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
	directory := DataNodeDirectory{
		DataNodes: make([]*DataNode, 0),
		MX:        &sync.Mutex{},
	}
	return &directory
}

func (d *DataNodeDirectory) AddDataNode(dataNode *DataNode) error {
	d.MX.Lock()
	d.DataNodes = append(d.DataNodes, dataNode)
	d.MX.Unlock()
	return nil
}

func (d *DataNodeDirectory) GetDataNode(index int) *DataNode {
	dataNode := d.DataNodes[index]
	return dataNode
}

func (d *DataNodeDirectory) UpdateDataNodeState(index int, state DataNodeState) {
	d.MX.Lock()
	dataNode := d.DataNodes[index]
	dataNode.State = state
	d.MX.Unlock()
}

func (d *DataNodeDirectory) GetDataNodeCount() int {
	d.MX.Lock()
	defer d.MX.Unlock()
	return len(d.DataNodes)
}
