/* # The Graph-Based Approaches involve converting time series data into graph structures, where nodes represent specific
data points or states, and edges represent transitions between these states. This method can be powerful for
identifying complex patterns, especially when the relationships between data points are non-linear or involve
intricate dependencies.
- Graph Representation: Convert the time series into a graph where nodes represent data points or specific states
	(e.g., certain price levels in candlestick data), and edges represent transitions between these nodes
	(e.g., price changes).
- Pattern Matching: The target pattern is also represented as a graph. The goal is to find a subgraph within
	the larger time series graph that matches the pattern graph.
- Graph Search Algorithms: Techniques such as Depth-First Search (DFS), Breadth-First Search (BFS), or more advanced
	methods like Subgraph Isomorphism are used to find matches between the pattern and the time series graph.*/

/* # Enhancing the Graph-Based Approach
1. Multi-Edge and Multi-Node Patterns
- Multi-Edge: If your patterns involve more complex relationships(e.g. ,multi possible transitions from a single node),
	you might represent these relationships with multiple edges.
- Multi-Node: For patterns involving more than simple linear transitions, consider using a more complex graph structure
	where nodes represent states or ranges of values rather than single data points.
2.Advanced Graph Matching Algorithms
- Subgraph Isomorphism: For more complex patterns, use advanced subgraph isomorphism algorithms, which allow for more
	flexible pattern matching, including patterns with varying lengths and structures.
- Graph Embeddings: Represent nodes and edges as vectors in a high-dimensional space, allowing for more sophisticated
	pattern matching using machine learning techniques.
3. Graph-Based Anomaly Detection
- Anomaly Detection: Extend the graph-based approach to detect anomalies by identifying unexpected subgraphs
	or transitions that do not match any known patterns.*/

package gba

import "fmt"

// Example candlestick data (price changes)
var data = []float64{100.0, 101.5, 102.0, 101.8, 102.3, 103.5, 102.8, 104.0, 103.5, 103.0, 102.0}

// Pattern to match (a specific price transition pattern)
var pattern = []float64{101.8, 102.3, 103.5}

// Graph Node structure
type Node struct {
	Value float64
	Edges []*Node
}

// Function to add an edge between nodes
func (n *Node) AddEdge(neighbor *Node) {
	n.Edges = append(n.Edges, neighbor)
}

// BuildGraph creates a graph from the given data. Construct a graph from the candlestick data. Each node will represent
// a price point, and edges will represent price changes between consecutive points.
func BuildGraph(data []float64) []*Node {
	nodes := make([]*Node, len(data))
	for i, value := range data {
		nodes[i] = &Node{Value: value}
		if i > 0 {
			nodes[i-1].AddEdge(nodes[i])
		}
	}
	return nodes
}

// DFS for Subgraph Isomorphism
func DFS(node *Node, pattern []float64, index int) bool {
	if index == len(pattern) {
		return true
	}
	for _, neighbor := range node.Edges {
		if neighbor.Value == pattern[index] && DFS(neighbor, pattern, index+1) {
			return true
		}
	}
	return false
}

// FindPatternInGraph searches for the pattern in the graph
func FindPatternInGraph(graph []*Node, pattern []float64) *Node {
	for _, node := range graph {
		if node.Value == pattern[0] && DFS(node, pattern, 1) {
			return node
		}
	}
	return nil
}

func Analysis() {
	// Build the graph from the data
	graph := BuildGraph(data)

	// Search for the pattern in the graph
	matchNode := FindPatternInGraph(graph, pattern)

	if matchNode != nil {
		fmt.Printf("Pattern matched starting at value: %f\n", matchNode.Value)
	} else {
		fmt.Println("Pattern not found.")
	}
}

//func ConvertForAnal()  {}
