package coordinator

import "time"

type NodeStatus string

const (
	NodeUp   NodeStatus = "up"
	NodeDown NodeStatus = "down"
)

type Node struct {
	ID             string
	Address        string
	Status         NodeStatus
	LastHeartbeat  time.Time
	TotalSpace     uint64
	AvailableSpace uint64
}

func NewNode(id string, address string, totalSpace uint64, availableSpace uint64) *Node {
	return &Node{
		ID:             id,
		Address:        address,
		Status:         NodeUp,
		TotalSpace:     totalSpace,
		AvailableSpace: availableSpace,
	}
}
