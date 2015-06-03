package breezynlp

import (
	"fmt"
)

type BreezyNode struct {
	Index    int                    //Index for node
	Payload  string                 // The string to be stored in node
	Children []BreezyNeighborObject // Array of nodes connected to current node with costs to node
}

func (brNode BreezyNode) String() string {
	return fmt.Sprintf("%v %v %v\n", brNode.Index, brNode.Payload, brNode.Children)
}

func (brNode *BreezyNode) AddChild(newChild BreezyNeighborObject) bool {
	// Add child inserts a BreezyNeighberObject in to the Children array
	// If successful returns true else it returns false
	for i := 0; i < len(brNode.Children); i++ { //loop through the children of the node
		if brNode.Children[i].Cost == newChild.Cost && brNode.Children[i].Vertex.Payload == newChild.Vertex.Payload {
			//if child is already present return false
			return false
		} else if brNode.Children[i].Cost > newChild.Cost && brNode.Children[i].Vertex.Payload == newChild.Vertex.Payload {
			// if child is found but new child has lower cost replace with new cost
			brNode.Children[i].Cost = newChild.Cost
			return true
		}
	}
	// if child is not found after loop through child add new neighbor
	brNode.Children = append(brNode.Children, newChild)
	return true
}

func (brNode *BreezyNode) removeChild(childToRemove BreezyNode) {
	// removeChild function removes the neighbor connection for a node
	foundIndex := -1
	var tempArr []BreezyNeighborObject
	for i := 0; i < len(brNode.Children); i++ { // loop through children to find the correct child
		if brNode.Children[i].Vertex.Index == childToRemove.Index && brNode.Children[i].Vertex.Payload == childToRemove.Payload {
			// If child is found store the index of that correct child
			foundIndex = i
		}
	}
	// make temp array of all items after found index and shorten orginal array to all items before index
	tempArr = brNode.Children[foundIndex+1 : len(brNode.Children)]
	brNode.Children = brNode.Children[0 : foundIndex-1]
	for j := 0; j < len(tempArr); j++ { // loop through tempArr to add back elements after the index
		brNode.Children = append(brNode.Children, tempArr[j])
	}
}

type BreezyNeighborObject struct {
	// BreezyNeighborObject stores the neccessary data a connection between two nodes
	Vertex BreezyNode // Connecting neighbor node
	Cost   int        // Cost to go to new vertex from original
}

func (brNeighbor BreezyNeighborObject) String() string {
	return fmt.Sprintf("%v  %v", brNeighbor.Vertex.Payload, brNeighbor.Cost)
}

type BreezyGraph struct {
	// BreezyGraph stores the neccessary data to make a graph
	BreezyADJList     []BreezyNode // Array of vertices within the graph
	NumberOfVerticies int          // number of vertices in the graph
	NumberOfEdges     int          // number of edges in the graph
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

func (brGraph *BreezyGraph) RemoveVertex(vertexToRemove BreezyNode) bool {
	// Make a queue to place other vertices connected to this vertex

	// Need to fix
	foundIndex := -1
	neighborQueue := BreezyQueue{nil, nil, 0}
	for i := 0; i < len(brGraph.BreezyADJList); i++ {
		if brGraph.BreezyADJList[i].Index == vertexToRemove.Index && brGraph.BreezyADJList[i].Payload == vertexToRemove.Payload {
			foundIndex = i
			for j := 0; j < len(brGraph.BreezyADJList[i].Children); j++ {
				neighborQueue.enqueue(BreezyQueueNode{brGraph.BreezyADJList[i].Children[j].Vertex.Index, brGraph.BreezyADJList[i].Children[j].Vertex.Payload, nil})
			}
			// Remove children in queue
			for neighborQueue.Length > 0 {
				tempNode := neighborQueue.dequeue()
				for k := 0; k < len(brGraph.BreezyADJList); k++ {
					if tempNode.Index == brGraph.BreezyADJList[k].Index && tempNode.Payload == brGraph.BreezyADJList[k].Payload {
						brGraph.RemoveEdge(brGraph.BreezyADJList[i], brGraph.BreezyADJList[k])
					}
				}
			}
		}
	}
	if foundIndex > -1 {
		tempArr := brGraph.BreezyADJList[i+1 : len(brGraph.BreezyADJList)]
		brGraph.BreezyADJList = brGraph.BreezyADJList[0 : i-1]
		for i := 0; i < len(tempArr); i++ {
			brGraph.BreezyADJList = append(brGraph.BreezyADJList, tempArr[i])
		}
		brGraph.NumberOfVerticies--
		return true
	} else {
		return false
	}
}

func (brGraph *BreezyGraph) RemoveEdge(fromVertex BreezyNode, andVertex BreezyNode) {
	firstHalf, secondHalf := false, false
	for i := 0; i < len(brGraph.BreezyADJList); i++ {
		if brGraph.BreezyADJList[i].Index == fromVertex.Index && brGraph.BreezyADJList[i].Payload == fromVertex.Payload {
			// remove child from fromVertex
			fromVertex.removeChild(andVertex)
			firstHalf = true
		}
		if brGraph.BreezyADJList[i].Index == andVertex.Index && brGraph.BreezyADJList[i].Payload == andVertex.Payload {
			// remove child from andvVrtex
			andVertex.removeChild(fromVertex)
			secondHalf = true
		}
	}
	if firstHalf == true && secondHalf == true {
		brGraph.NumberOfEdges--
	}
}

// Queue

type BreezyQueueNode struct {
	// BreezyQueueNode stores the node data for a queue
	Index   int              // Index from BreezyNode stored in queue
	Payload string           // Payload from BreezyNode stored in a queue
	Next    *BreezyQueueNode // Pointer to the queue node connected to this node
}
type BreezyQueue struct {
	// BreezyQueue stores the data for a queue data structure
	First  *BreezyQueueNode // Pointer to the first element in the queue
	Last   *BreezyQueueNode // Pointer to the last element in the queue
	Length int              // Value for the length of the queue
}

func (brQueue *BreezyQueue) enqueue(newNode BreezyQueueNode) {
	// Enqueue inserts a new queue node at he end of the queue
	if brQueue.First == nil {
		brQueue.First = &newNode
		brQueue.Last = brQueue.First
	} else {
		brQueue.Last.Next = &newNode
		brQueue.Last = &newNode
	}
	brQueue.Length++
}

func (brQueue *BreezyQueue) dequeue() BreezyQueueNode {
	// Dequeue pops off the first element in the queue
	if brQueue.First != nil {
		returnNode := brQueue.First
		brQueue.First = brQueue.First.Next
		return *returnNode
	}
	return BreezyQueueNode{-1, "", nil}
}
