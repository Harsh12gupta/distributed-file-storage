package coordinator

import "time"

type NodeRegistry struct {
	nodes map[string]Node
}

func NewNodeRegistry() *NodeRegistry {
	return &NodeRegistry{
		nodes: make(map[string]Node),
	}
}

func (r *NodeRegistry) AddNode(node *Node) {
	node.LastHeartbeat = time.Now()
	node.Status = NodeUp
	r.nodes[node.ID] = *node
}

func (r *NodeRegistry) GetNode(id string) (*Node, bool) {
	node, exists := r.nodes[id]
	if !exists {
		return nil, false
	}
	return &node, true
}
