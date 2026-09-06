package coordinator

import (
	"sort"
	"time"
)

type NodeRegistry struct {
	nodes map[string]*Node
}

func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		nodes: make(map[string]*Node),
	}
}

func (r *NodeRegistry) AddNode(node *Node) {
	node.LastHeartbeat = time.Now()
	node.Status = NodeUp
	r.nodes[node.ID] = node
}

func (r *NodeRegistry) GetNode(id string) (*Node, bool) {
	node, exists := r.nodes[id]
	if !exists {
		return nil, false
	}
	return node, true
}

func (r *NodeRegistry) UpdateHeartbeat(id string) bool {
	node, exists := r.nodes[id]
	if !exists {
		return false
	}
	node.LastHeartbeat = time.Now()
	node.Status = NodeUp
	return true
}

func (r *NodeRegistry) CheckNodeTimeouts(currentTime time.Time, timeoutDuration time.Duration) {
	for _, node := range r.nodes {
		if currentTime.Sub(node.LastHeartbeat) > timeoutDuration {
			node.Status = NodeDown
		}
	}
}

func (r *NodeRegistry) GetActiveNodes() []*Node {
	activeNodes := make([]*Node, 0)
	for _, node := range r.nodes {
		if node.Status == NodeUp {
			activeNodes = append(activeNodes, node)
		}
	}

	sort.Slice(activeNodes, func(i, j int) bool {
		return activeNodes[i].ID < activeNodes[j].ID
	})
	return activeNodes
}
