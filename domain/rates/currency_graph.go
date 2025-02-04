/*
Dijkstra algorithm to determine a path from 1 code to another code through a linked list.
*/

package rates

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/CoreumFoundation/CoreDEX-API/domain/denom"
	tradegrpc "github.com/CoreumFoundation/CoreDEX-API/domain/trade"
	tradeclient "github.com/CoreumFoundation/CoreDEX-API/domain/trade/client"
	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

type node struct {
	name    string
	value   int
	through *node
}

type edge struct {
	node   *node
	weight int
}

type weightedGraph struct {
	nodes []*node
	edges map[string][]*edge
	mutex sync.RWMutex
}

type heap struct {
	elements []*node
}

func key(denom1 *denom.Denom) string {
	return fmt.Sprintf("%s-%s", denom1.Currency, denom1.Issuer)
}

func getCurrencyPath(graph *weightedGraph, denom1 *denom.Denom, target string) []string {
	// Lock the graph: We are working against a pointer and data is returned as rendered nodes for a given source
	// (Alternative would be a deep copy of the graph, or decouple of the nodes and edges from the weightedGraph)
	graph.mutex.Lock()
	// Clone graph in workGraph to prevent mutation of the source graph
	workGraph := cloneGraph(graph)
	graph.mutex.Unlock()
	source := key(denom1)
	var err error
	workGraph, err = dijkstra(workGraph, source)
	if err != nil {
		logger.Errorf("getCurrencyPath: Error getting path: %v", err)
		return nil
	}
	// find all the possible solutions
	possibleSolution := make(map[string][]string)
	for _, node := range workGraph.nodes {
		for n := node; n.through != nil; n = n.through {
			possibleSolution[node.name] = append(possibleSolution[node.name], n.name)
		}
	}
	if possibleSolutionPath, ok := possibleSolution[target]; ok {
		// Check if the last value is the source (The algorithm has a deviation which does not always add the source to the path)
		if possibleSolutionPath[len(possibleSolutionPath)-1] != source {
			// Add the source
			possibleSolutionPath = append(possibleSolutionPath, source)
		}
		// re-order the solution so it is from source to target
		for i, j := 0, len(possibleSolutionPath)-1; i < j; i, j = i+1, j-1 {
			possibleSolutionPath[i], possibleSolutionPath[j] = possibleSolutionPath[j], possibleSolutionPath[i]
		}
		logger.Infof("Solution: %s to %s target is %v", denom1.Currency, target, possibleSolutionPath)
		return possibleSolutionPath
	}
	return nil
}

func dijkstra(graph *weightedGraph, lookupKey string) (*weightedGraph, error) {
	visited := make(map[string]bool)
	heap := &heap{}

	startNode := graph.getNode(lookupKey)
	if startNode == nil {
		return nil, fmt.Errorf("node %s not found in graph", lookupKey)
	}
	startNode.value = 0
	heap.push(startNode)

	for heap.size() > 0 {
		current := heap.pop()
		visited[current.name] = true
		edges := graph.edges[current.name]
		for _, edge := range edges {
			if !visited[edge.node.name] {
				heap.push(edge.node)
				if current.value+edge.weight < edge.node.value {
					edge.node.value = current.value + edge.weight
					edge.node.through = current
				}
			}
		}
	}
	return graph, nil
}

// Loads the trade pairs from the trade store and updates the graph
// Reloads every 60 minutes in case a new pool was created
func (f *Fetcher) loadTradePairs() {
	f.graph = newGraph() // Specific initialization is required to prevent an impossible to lock scenario and potential lock crashes
	for {
		// Load all the pairs
		pairs, err := f.cl.tradeStore.GetTradePairs(tradeclient.AuthCtx(context.Background()),
			&tradegrpc.TradePairFilter{Network: f.network})
		if err != nil {
			logger.Errorf("loadTradePairs: Degraded functionality, not all values can be converted to USD: Error getting pairs: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}
		graph := newGraph()
		// Transform the pairs into the nodes:
		currencies := make(map[string]bool, 0)
		for _, pair := range pairs.TradePairs {
			currencies[key(pair.Denom1)] = true
			currencies[key(pair.Denom2)] = true
			logger.Infof("loadTradePairs: Adding trade pair %s-%s %s-%s",
				pair.Denom1.Currency, pair.Denom1.Issuer, pair.Denom2.Currency, pair.Denom2.Issuer)
		}
		nodes := addNodes(graph, currencies)
		// Add the edges between the nodes:
		for _, pair := range pairs.TradePairs {
			graph.addEdge(nodes[key(pair.Denom1)], nodes[key(pair.Denom2)], 0)
			logger.Infof("loadTradePairs: Adding edge %s %s", pair.Denom1.ToString(), pair.Denom2.ToString())
		}
		graph.mutex.Lock()
		f.graph = graph
		graph.mutex.Unlock()
		time.Sleep(60 * time.Minute)
	}
}

func (f *Fetcher) initGraph() {
	go f.loadTradePairs()
}

func newGraph() *weightedGraph {
	return &weightedGraph{
		edges: make(map[string][]*edge),
	}
}

func (g *weightedGraph) getNode(name string) (node *node) {
	for _, n := range g.nodes {
		if n.name == name {
			node = n
		}
	}
	return
}

func (g *weightedGraph) addNode(n *node) {
	g.nodes = append(g.nodes, n)
}

func addNodes(graph *weightedGraph, names map[string]bool) (nodes map[string]*node) {
	nodes = make(map[string]*node)
	for name := range names {
		n := &node{name, math.MaxInt, nil}
		graph.addNode(n)
		nodes[name] = n
	}
	return
}

func (g *weightedGraph) addEdge(n1, n2 *node, weight int) {
	g.edges[n1.name] = append(g.edges[n1.name], &edge{n2, weight})
	g.edges[n2.name] = append(g.edges[n2.name], &edge{n1, weight})
}

func (n *node) String() string {
	return n.name
}

func (e *edge) String() string {
	return e.node.String() + "(" + strconv.Itoa(e.weight) + ")"
}

func (g *weightedGraph) String() (s string) {
	g.mutex.RLock()
	for _, n := range g.nodes {
		s = s + n.String() + " ->"
		for _, c := range g.edges[n.name] {
			s = s + " " + c.node.String() + " (" + strconv.Itoa(c.weight) + ")"
		}
		s = s + "\n"
	}
	g.mutex.RUnlock()
	return
}

func (h *heap) size() int {
	return len(h.elements)
}

// push an element to the heap, re-arrange the heap
func (h *heap) push(element *node) {
	h.elements = append(h.elements, element)
	i := len(h.elements) - 1
	for ; h.elements[i].value < h.elements[h.parent(i)].value; i = h.parent(i) {
		h.swap(i, h.parent(i))
	}
}

// pop the top of the heap, which is the min value
func (h *heap) pop() (i *node) {
	i = h.elements[0]
	h.elements[0] = h.elements[len(h.elements)-1]
	h.elements = h.elements[:len(h.elements)-1]
	h.rearrange(0)
	return
}

// rearrange the heap
func (h *heap) rearrange(i int) {
	smallest := i
	left, right, size := h.leftChild(i), h.rightChild(i), len(h.elements)
	if left < size && h.elements[left].value < h.elements[smallest].value {
		smallest = left
	}
	if right < size && h.elements[right].value < h.elements[smallest].value {
		smallest = right
	}
	if smallest != i {
		h.swap(i, smallest)
		h.rearrange(smallest)
	}
}

func (h *heap) swap(i, j int) {
	h.elements[i], h.elements[j] = h.elements[j], h.elements[i]
}

func (*heap) parent(i int) int {
	return (i - 1) / 2
}

func (*heap) leftChild(i int) int {
	return 2*i + 1
}

func (*heap) rightChild(i int) int {
	return 2*i + 2
}

func (h *heap) String() (str string) {
	return fmt.Sprintf("%q\n", getNames(h.elements))
}

func getNames(nodes []*node) (names []string) {
	for _, node := range nodes {
		names = append(names, node.name)
	}
	return
}

// Clone the graph to prevent pointer data mutation
func cloneGraph(graph *weightedGraph) *weightedGraph {
	newGraph := newGraph()
	for _, nd := range graph.nodes {
		n := &node{name: nd.name, value: nd.value, through: nd.through}
		newGraph.addNode(n)
	}
	for key, edges := range graph.edges {
		for _, edge := range edges {
			newGraph.addEdge(newGraph.getNode(key), newGraph.getNode(edge.node.name), edge.weight)
		}
	}
	return newGraph
}
