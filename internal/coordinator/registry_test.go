package coordinator

import "testing"

func TestNodeRegistry(t *testing.T) {
	node := NewNode("node-1", "localhost:9001")
	registry := NewNodeRegistry()
	registry.AddNode(node)

	retrievedNode, exists := registry.GetNode("node-1")
	if !exists {
		t.Errorf("expected node to exist in registry")
	}

	if retrievedNode.ID != node.ID {
		t.Errorf("expected node ID %q, got %q", node.ID, retrievedNode.ID)
	}

	if retrievedNode.Address != node.Address {
		t.Errorf("expected address %q, got %q", node.Address, retrievedNode.Address)
	}

	if retrievedNode.Status != node.Status {
		t.Errorf("expected status %q, got %q", node.Status, retrievedNode.Status)
	}

	if retrievedNode.LastHeartbeat.IsZero() {
		t.Errorf("expected LastHeartbeat to be set, but it is zero")
	}

	//Test non-existent node
	_, exists = registry.GetNode("node-2")
	if exists {
		t.Errorf("expected node to not exist in registry")
	}
}

func TestDuplicateNodeRegistration(t *testing.T) {
	registry := NewNodeRegistry()

	node1 := NewNode("node-1", "localhost:9001")
	registry.AddNode(node1)

	node2 := NewNode("node-1", "localhost:9002")
	registry.AddNode(node2)

	storedNode, exists := registry.GetNode("node-1")

	if !exists {
		t.Fatal("expected node to exist")
	}

	if storedNode.Address != "localhost:9002" {
		t.Errorf(
			"expected address localhost:9002, got %s",
			storedNode.Address,
		)
	}
}
