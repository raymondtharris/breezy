package breezynlp

import (
	"fmt"
)

type BreezyNode struct {
	Index    int
	Payload  string
	Children []BreezyNode
}

type BreezyADJObject struct {
	Vertex BreezyNode
}

type BreezyGraph struct {
}
