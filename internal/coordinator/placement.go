package coordinator

type Chunk struct {
	Index int
	Size  uint64
}

type Placement struct {
	ChunkIndex int
	NodeID     string
}

type Placer struct {
	registry *NodeRegistry
}

func NewPlacer(registry *NodeRegistry) *Placer {
	return &Placer{
		registry: registry,
	}
}

func (p *Placer) PlaceChunks(chunks []*Chunk) ([]*Placement, error) {
	placements := make([]*Placement, 0, len(chunks))
	activeNodes := p.registry.GetActiveNodes()

	if len(activeNodes) == 0 {
		return nil, ErrNoActiveNodes
	}

	remainingSpace := make(map[string]uint64)

	for _, node := range activeNodes {
		remainingSpace[node.ID] = node.AvailableSpace
	}

	for _, chunk := range chunks {
		var selectedNode *Node
		for _, node := range activeNodes {
			if remainingSpace[node.ID] >= chunk.Size {
				if selectedNode == nil || remainingSpace[node.ID] > remainingSpace[selectedNode.ID] {
					selectedNode = node
				}
			}
		}

		if selectedNode == nil {
			return nil, ErrInsufficientSpace
		} else {
			remainingSpace[selectedNode.ID] -= chunk.Size
			placements = append(placements, &Placement{
				ChunkIndex: chunk.Index,
				NodeID:     selectedNode.ID,
			})
		}
	}
	return placements, nil
}
