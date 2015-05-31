package breezynlp

import (
	"fmt"
)

type BreezyNode struct {
	Index    int
	Payload  string
	Children []BreezyNeighborObject
}

func (brNode BreezyNode) String() string {
	return fmt.Sprintf("%v %v %v\n", brNode.Index, brNode.Payload, brNode.Children)
}

func (brNode *BreezyNode) AddChild(newChild BreezyNeighborObject) bool {
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

func (brGraph BreezyGraph) String() string {
	return fmt.Sprintf("%v %v \n%v", brGraph.NumberOfVerticies, brGraph.NumberOfEdges, brGraph.BreezyADJList)
}

func (brGraph *BreezyGraph) AddVertex(newVertex BreezyNode) {
	// AddVertex inserts a new BreezyNode in to the BreezyADJList array
	//fmt.Println(newVertex)
	brGraph.BreezyADJList = append(brGraph.BreezyADJList, newVertex)
	brGraph.NumberOfVerticies++
}

func (brGraph *BreezyGraph) AddEdge(betweenVertex BreezyNode, andNeighbor BreezyNeighborObject) {
	//AddEdge inserts a link between two nodes that is not directed.
	isInGraph, neighborInGraph := false, false
	//Add link for initial direction
	for i := 0; i < len(brGraph.BreezyADJList); i++ {
		if brGraph.BreezyADJList[i].Index == betweenVertex.Index && brGraph.BreezyADJList[i].Payload == betweenVertex.Payload {
			isInGraph = true
			brGraph.BreezyADJList[i].AddChild(andNeighbor)

		}
		if brGraph.BreezyADJList[i].Index == andNeighbor.Vertex.Index && brGraph.BreezyADJList[i].Payload == andNeighbor.Vertex.Payload {
			neighborInGraph = true
			brGraph.BreezyADJList[i].AddChild(BreezyNeighborObject{betweenVertex, andNeighbor.Cost})
		}

	}
	if !isInGraph {
		brGraph.AddVertex(betweenVertex)
		brGraph.BreezyADJList[len(brGraph.BreezyADJList)].AddChild(andNeighbor)
	}
	if !neighborInGraph {
		brGraph.AddVertex(andNeighbor.Vertex)
		brGraph.BreezyADJList[len(brGraph.BreezyADJList)].AddChild(BreezyNeighborObject{betweenVertex, andNeighbor.Cost})
	}
	brGraph.NumberOfEdges++
}
