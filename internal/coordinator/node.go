package coordinator

import "time"

type NodeStatus string

const (
	NodeUp   NodeStatus = "up"
	NodeDown NodeStatus = "down"
)

type Node struct {
	ID            string
	Address       string
	Status        NodeStatus
	LastHeartbeat time.Time
}

func NewNode(id string, address string) *Node {
	return &Node{
		ID:            id,
		Address:       address,
		Status:        NodeUp,
		LastHeartbeat: time.Now(),
	}
}
