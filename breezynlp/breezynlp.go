package breezynlp

import (
	"fmt"
)

type BreezyNode struct {
	Index    int
	Payload  string
	Children []BreezyNeighborObject
}

func (brNode BreezyNode) AddChild(newChild BreezyNeighborObject) bool {
	// Add child inserts a BreezyNeighberObject in to the Children array
	// If successful returns true else it returns false
	for i := 0; i < len(brNode.Children); i++ {
		if brNode.Children[i].Cost == newChild.Cost && brNode.Children[i].Vertex.Payload == newChild.Vertex.Payload {
			return false
		}
	}
	brNode.Children = append(brNode.Children, newChild)
	return true
}

type BreezyNeighborObject struct {
	Vertex BreezyNode
	Cost   int
}

func (brNeighbor BreezyNeighborObject) String() string {
	return fmt.Sprintf("%v  %v", brNeighbor.Vertex.Payload, brNeighbor.Cost)
}

type BreezyGraph struct {
	BreezyADJList     []BreezyNode
	NumberOfVerticies int
	NumberOfEdges     int
}

func (brGraph BreezyGraph) AddVertex(newVertex BreezyNode) {
	// AddVertex inserts a new BreezyNode in to the BreezyADJList array
	brGraph.BreezyADJList = append(brGraph.BreezyADJList, newVertex)
	brGraph.NumberOfVerticies++
}
