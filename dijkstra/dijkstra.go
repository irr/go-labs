package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Graph struct {
	VertexArray []*Vertex
}

type Vertex struct {
	Id      string
	Visited bool
	AdjEdge []*Edge
}

type Edge struct {
	Source      *Vertex
	Destination *Vertex
	Weight      int
}

func NewGraph() *Graph {
	return &Graph{
		make([]*Vertex, 0),
	}
}

func NewVertex(input_id string) *Vertex {
	return &Vertex{
		Id:      input_id,
		Visited: false,
		AdjEdge: make([]*Edge, 0),
	}
}

func NewEdge(source, destination *Vertex, weight int) *Edge {
	return &Edge{
		source,
		destination,
		weight,
	}
}

func StrToInt(input_str string) int {
	result, err := strconv.Atoi(input_str)
	if err != nil {
		panic("failed to convert string")
	}
	return result
}

func (G *Graph) AddVertexs(more ...*Vertex) {
	for _, vertex := range more {
		G.VertexArray = append(G.VertexArray, vertex)
	}
}

func (G *Graph) GetVertexByID(id string) *Vertex {
	for _, vertex := range G.VertexArray {
		if vertex.Id == id {
			return vertex
		}
	}
	return nil
}

//Find the node with the id, or create it.
func (G *Graph) GetOrConst(id string) *Vertex {
	vertex := G.GetVertexByID(id)
	if vertex == nil {
		vertex = NewVertex(id)
		G.AddVertexs(vertex)
	}
	return vertex
}

func (A *Vertex) AddEdges(more ...*Edge) {
	for _, edge := range more {
		A.AdjEdge = append(A.AdjEdge, edge)
	}
}

func ImportData(input_str string) *Graph {
	input_str = strings.TrimSpace(input_str)
	lines := strings.Split(input_str, "\n")

	new_graph := NewGraph()

	for _, line := range lines {
		fields := strings.Split(line, "|")

		SourceID := fields[0]
		edgepairs := fields[1:]

		new_graph.GetOrConst(SourceID)

		for _, pair := range edgepairs {
			if len(strings.Split(pair, ",")) == 1 {
				//to skip
				continue
			}
			DestinationID := strings.Split(pair, ",")[0]
			weight := StrToInt(strings.Split(pair, ",")[1])

			src_vertex := new_graph.GetOrConst(SourceID)
			des_vertex := new_graph.GetOrConst(DestinationID)

			//Connect bi-direction
			edge1 := NewEdge(src_vertex, des_vertex, weight)
			src_vertex.AddEdges(edge1)

			edge2 := NewEdge(des_vertex, src_vertex, weight)
			des_vertex.AddEdges(edge2)
		}
	}
	return new_graph
}

func (A *Vertex) GetAdEdg() chan *Edge {
	edgechan := make(chan *Edge)

	go func() {
		defer close(edgechan)
		for _, edge := range A.AdjEdge {
			edgechan <- edge
		}
	}()

	return edgechan
}

func DFS(StartSource *Vertex) {
	if StartSource.Visited {
		return
	}

	StartSource.Visited = true
	fmt.Printf("%v ", StartSource.Id)

	for edge := range StartSource.GetAdEdg() {
		DFS(edge.Destination)
	}
}

const MAXWEIGHT = 1000000

type MinDistanceFromSource map[*Vertex]int

func (G *Graph) Dijks(StartSource, TargetSource *Vertex) MinDistanceFromSource {
	D := make(MinDistanceFromSource)
	for _, vertex := range G.VertexArray {
		D[vertex] = MAXWEIGHT
	}
	D[StartSource] = 0

	for edge := range StartSource.GetAdEdg() {
		D[edge.Destination] = edge.Weight
	}
	CalculateDistance(StartSource, TargetSource, D)
	return D
}

func CalculateDistance(StartSource, TargetSource *Vertex, D MinDistanceFromSource) {
	for edge := range StartSource.GetAdEdg() {
		if D[edge.Destination] > D[edge.Source]+edge.Weight {
			D[edge.Destination] = D[edge.Source] + edge.Weight
		} else if D[edge.Destination] < D[edge.Source]+edge.Weight {
			continue
		}
		CalculateDistance(edge.Destination, TargetSource, D)
	}
}

func main() {

	//A ~ F : IDs of vertices
	//A|B,7|C,10|D,15
	//weight(or distance) is 7 from A to B
	//, 10 from A to C
	//, 15 from A to D
	str1 := `
A|B,7|C,9|F,20
B|A,7|C,10|D,15
C|A,9|B,10|D,11|E,30|F,2
D|B,15|C,11|E,2
E|C,30|D,2|F,9
F|A,20|C,2|E,9
`

	G1 := ImportData(str1)
	DFS(G1.GetVertexByID("A"))

	distmap1 := G1.Dijks(G1.GetVertexByID("A"), G1.GetVertexByID("E"))

	fmt.Println()
	for vertex1, distance1 := range distmap1 {
		fmt.Println(vertex1.Id, "=", distance1)
	}
	/*
	   A B C D E F
	   A = 0
	   B = 7
	   C = 9
	   F = 11
	   D = 20
	   E = 20
	*/
	fmt.Println()

	str2 := `
S|A,15|B,14|C,9
A|S,15|B,5|D,20|T,44
B|S,14|A,5|D,30|E,18
C|S,9|E,24
D|A,20|B,30|E,2|F,11|T,16
E|B,18|C,24|D,2|F,6|T,19
F|D,11|E,6|T,6
T|A,44|D,16|F,6|E,19
`

	G2 := ImportData(str2)
	DFS(G2.GetVertexByID("S"))

	distmap2 := G2.Dijks(G2.GetVertexByID("S"), G2.GetVertexByID("T"))

	fmt.Println()
	for vertex2, distance2 := range distmap2 {
		fmt.Println(vertex2.Id, "=", distance2)
	}
	/*
	   A S B D E C F T
	   S = 0
	   A = 15
	   B = 14
	   C = 9
	   D = 34
	   T = 44
	   E = 32
	   F = 38
	*/
}
