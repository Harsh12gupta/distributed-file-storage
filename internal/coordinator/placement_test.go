package coordinator

import (
	"testing"
)

func TestPlaceChunks(t *testing.T) {
	// Create a NodeRegistry and add some nodes
	registry := NewNodeRegistry()
	node1 := NewNode("node-1", "localhost:9001", 1000, 500)
	node2 := NewNode("node-2", "localhost:9002", 1000, 300)
	node3 := NewNode("node-3", "localhost:9003", 1000, 700)

	registry.AddNode(node1)
	registry.AddNode(node2)
	registry.AddNode(node3)

	chunk1 := &Chunk{Index: 0, Size: 200}
	chunk2 := &Chunk{Index: 1, Size: 200}
	chunk3 := &Chunk{Index: 2, Size: 100}

	chunks := []*Chunk{chunk1, chunk2, chunk3}

	placer := NewPlacer(registry)
	placements, err := placer.PlaceChunks(chunks)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(placements) != 3 {
		t.Fatalf("expected 3 placements, got %d", len(placements))
	}

	if placements[0].ChunkIndex != 0 {
		t.Errorf("expected chunk index 0, got %d", placements[0].ChunkIndex)
	}

	if placements[1].ChunkIndex != 1 {
		t.Errorf("expected chunk index 1, got %d", placements[1].ChunkIndex)
	}

	if placements[2].ChunkIndex != 2 {
		t.Errorf("expected chunk index 2, got %d", placements[2].ChunkIndex)
	}

	if placements[0].NodeID != "node-3" {
		t.Errorf("expected chunk 0 to be placed on node-3, got %s", placements[0].NodeID)
	}

	if placements[1].NodeID != "node-1" {
		t.Errorf("expected chunk 1 to be placed on node-1, got %s", placements[1].NodeID)
	}

	if placements[2].NodeID != "node-3" {
		t.Errorf("expected chunk 2 to be placed on node-3, got %s", placements[2].NodeID)
	}

}

func TestPlaceChunksInsufficientSpace(t *testing.T) {
	registry := NewNodeRegistry()
	node1 := NewNode("node-1", "localhost:9001", 1000, 100)
	node2 := NewNode("node-2", "localhost:9002", 1000, 50)

	registry.AddNode(node1)
	registry.AddNode(node2)

	chunk1 := &Chunk{Index: 0, Size: 100}
	chunk2 := &Chunk{Index: 1, Size: 300}

	chunks := []*Chunk{chunk1, chunk2}

	placer := NewPlacer(registry)
	_, err := placer.PlaceChunks(chunks)

	if err == nil {
		t.Fatalf("expected error due to insufficient space, got nil")
	}

	if err != ErrInsufficientSpace {
		t.Fatalf("expected ErrInsufficientSpace, got %v", err)
	}
}

func TestPlaceChunksIgnoreDownNodes(t *testing.T) {
	registry := NewNodeRegistry()
	node1 := NewNode("node-1", "localhost:9001", 1000, 500)
	node2 := NewNode("node-2", "localhost:9002", 1000, 300)
	node3 := NewNode("node-3", "localhost:9003", 1000, 700)

	registry.AddNode(node1)
	registry.AddNode(node2)
	registry.AddNode(node3)

	// Mark node2 as down
	node2.Status = NodeDown

	chunk1 := &Chunk{Index: 0, Size: 500}
	chunk2 := &Chunk{Index: 1, Size: 500}
	chunk3 := &Chunk{Index: 2, Size: 100}

	chunks := []*Chunk{chunk1, chunk2, chunk3}

	placer := NewPlacer(registry)
	placements, err := placer.PlaceChunks(chunks)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(placements) != 3 {
		t.Fatalf("expected 3 placements, got %d", len(placements))
	}

	if placements[0].ChunkIndex != 0 {
		t.Errorf("expected chunk index 0, got %d", placements[0].ChunkIndex)
	}

	if placements[1].ChunkIndex != 1 {
		t.Errorf("expected chunk index 1, got %d", placements[1].ChunkIndex)
	}

	if placements[2].ChunkIndex != 2 {
		t.Errorf("expected chunk index 2, got %d", placements[2].ChunkIndex)
	}

	if placements[0].NodeID != "node-3" {
		t.Errorf("expected chunk 0 to be placed on node-3, got %s", placements[0].NodeID)
	}

	if placements[1].NodeID != "node-1" {
		t.Errorf("expected chunk 1 to be placed on node-1, got %s", placements[1].NodeID)
	}

	if placements[2].NodeID != "node-3" {
		t.Errorf("expected chunk 2 to be placed on node-3, got %s", placements[2].NodeID)
	}
}

func TestPlaceChunksDoesNotModifyNodeCapacity(t *testing.T) {
	registry := NewNodeRegistry()

	node1 := NewNode("node-1", "localhost:9001", 1000, 500)
	node2 := NewNode("node-2", "localhost:9002", 1000, 300)

	registry.AddNode(node1)
	registry.AddNode(node2)

	chunks := []*Chunk{
		{Index: 0, Size: 200},
		{Index: 1, Size: 100},
	}

	placer := NewPlacer(registry)

	_, err := placer.PlaceChunks(chunks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if node1.AvailableSpace != 500 {
		t.Errorf("expected node1 available space to remain 500, got %d", node1.AvailableSpace)
	}

	if node2.AvailableSpace != 300 {
		t.Errorf("expected node2 available space to remain 300, got %d", node2.AvailableSpace)
	}
}
