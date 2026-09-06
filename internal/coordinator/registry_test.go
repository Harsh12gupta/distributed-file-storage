package coordinator

import (
	"testing"
	"time"
)

func TestNodeRegistry(t *testing.T) {
	node := NewNode("node-1", "localhost:9001", 1000, 1000)
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

	if retrievedNode.TotalSpace != node.TotalSpace {
		t.Errorf("expected TotalSpace to be %d, got %d", node.TotalSpace, retrievedNode.TotalSpace)
	}

	if retrievedNode.AvailableSpace != node.AvailableSpace {
		t.Errorf("expected AvailableSpace to be %d, got %d", node.AvailableSpace, retrievedNode.AvailableSpace)
	}

	//Test non-existent node
	_, exists = registry.GetNode("node-2")
	if exists {
		t.Errorf("expected node to not exist in registry")
	}
}

func TestDuplicateNodeRegistration(t *testing.T) {
	registry := NewNodeRegistry()

	node1 := NewNode("node-1", "localhost:9001", 1000, 1000)
	registry.AddNode(node1)

	node2 := NewNode("node-1", "localhost:9002", 1000, 1000)
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

func TestUpdateHeartbeat(t *testing.T) {
	registry := NewNodeRegistry()

	node1 := NewNode("node-1", "localhost:9001", 1000, 1000)
	registry.AddNode(node1)

	// Get the heartbeat before updating it
	storedNode, exists := registry.GetNode("node-1")

	if !exists {
		t.Fatal("expected node to exist")
	}

	previousHeartbeat := storedNode.LastHeartbeat

	// Update through the registry
	registry.UpdateHeartbeat("node-1")

	// Get the node again
	updatedNode, exists := registry.GetNode("node-1")

	if !exists {
		t.Fatal("expected node to exist")
	}

	if !previousHeartbeat.Before(updatedNode.LastHeartbeat) {
		t.Error("expected LastHeartbeat to be updated")
	}
}

func TestUpdateHeartbeatUnknownNode(t *testing.T) {
	registry := NewNodeRegistry()

	success := registry.UpdateHeartbeat("node-999")

	if success {
		t.Error("expected heartbeat update to fail for unknown node")
	}
}

func TestNodeStatusAfterHeartbeatUpdate(t *testing.T) {
	node := NewNode("node-1", "localhost:9001", 1000, 1000)
	registry := NewNodeRegistry()
	registry.AddNode(node)
	node.Status = NodeDown // Simulate node being down

	registry.UpdateHeartbeat("node-1")

	updatedNode, _ := registry.GetNode("node-1")

	if updatedNode.Status != NodeUp {
		t.Errorf(
			"expected node status to be %q after heartbeat update, got %q",
			NodeUp,
			updatedNode.Status,
		)
	}
}

func TestNodeTimeOut(t *testing.T) {
	node1 := NewNode("node-1", "localhost:9001", 1000, 1000)
	registry := NewNodeRegistry()
	registry.AddNode(node1)

	node2 := NewNode("node-2", "localhost:9002", 1000, 1000)
	registry.AddNode(node2)

	node3 := NewNode("node-3", "localhost:9003", 1000, 1000)
	registry.AddNode(node3)

	currentTime := time.Now()

	// Simulate node1 being down for more than 15 seconds
	node1.LastHeartbeat = currentTime.Add(-16 * time.Second)

	// Simulate node2 being down for less than 15 seconds
	node2.LastHeartbeat = currentTime.Add(-10 * time.Second)

	node3.LastHeartbeat = currentTime.Add(-15 * time.Second)

	registry.CheckNodeTimeouts(currentTime, 15*time.Second)

	if node1.Status != NodeDown {
		t.Errorf("expected node1 status to be %q, got %q", NodeDown, node1.Status)
	}

	if node2.Status != NodeUp {
		t.Errorf("expected node2 status to be %q, got %q", NodeUp, node2.Status)
	}

	if node3.Status != NodeUp {
		t.Errorf("expected node3 status to be %q, got %q", NodeUp, node3.Status)
	}

}

func TestGetActiveNodes(t *testing.T) {
	node1 := NewNode("node-1", "localhost:9001", 1000, 1000)
	registry := NewNodeRegistry()
	registry.AddNode(node1)

	node2 := NewNode("node-2", "localhost:9002", 1000, 1000)
	registry.AddNode(node2)

	node3 := NewNode("node-3", "localhost:9003", 1000, 1000)
	registry.AddNode(node3)

	node1.Status = NodeUp
	node2.Status = NodeDown
	node3.Status = NodeUp

	activeNodes := registry.GetActiveNodes()
	if len(activeNodes) != 2 {
		t.Errorf("expected 2 active nodes, got %d", len(activeNodes))
	}

	if activeNodes[0].ID != "node-1" {
		t.Errorf("expected node-1, got %s", activeNodes[0].ID)
	}

	if activeNodes[1].ID != "node-3" {
		t.Errorf("expected node-3, got %s", activeNodes[1].ID)
	}
}

func TestGetActiveNodesNoActiveNodes(t *testing.T) {
	registry := NewNodeRegistry()

	node1 := NewNode("node-1", "localhost:9001", 1000, 1000)
	node2 := NewNode("node-2", "localhost:9002", 1000, 1000)

	registry.AddNode(node1)
	registry.AddNode(node2)

	node1.Status = NodeDown
	node2.Status = NodeDown

	activeNodes := registry.GetActiveNodes()

	if len(activeNodes) != 0 {
		t.Errorf("expected 0 active nodes, got %d", len(activeNodes))
	}
}
