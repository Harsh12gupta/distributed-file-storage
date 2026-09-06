package coordinator

import "testing"

func TestNewNode(t *testing.T) {
	nodeID := "node-1"
	address := "localhost:9001"

	node := NewNode(nodeID, address, 1000, 1000)

	if node.ID != nodeID {
		t.Errorf("expected node ID %q, got %q", nodeID, node.ID)
	}

	if node.Address != address {
		t.Errorf("expected address %q, got %q", address, node.Address)
	}

	if node.Status != NodeUp {
		t.Errorf("expected status %q, got %q", NodeUp, node.Status)
	}

	if node.TotalSpace != 1000 {
		t.Errorf("expected TotalSpace to be 1000, got %d", node.TotalSpace)
	}

	if node.AvailableSpace != 1000 {
		t.Errorf("expected AvailableSpace to be 1000, got %d", node.AvailableSpace)
	}
}
