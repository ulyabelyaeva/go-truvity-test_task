package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

const INF = int((^uint(0)) >> 1)
const MAXN = 500    // maximum number of vertices
const MAXM = 300000 // maximum number of edges

var n, m, s, t int              // number of vertices, number of edges, source and drain vertices
var graph = make([][]int, MAXN) // in form of an adjacency list
var distance [MAXN]int          // distances from s to all other vertices, required to build a level graph
var pointer [MAXN]int           // pointer to the first remaining edge in an adjacency list

type edge struct {
	from, to, capacity, flow int
}

var edges = []edge{}

func add_edge(u, v, c int) {
	u -= 1
	v -= 1
	// add the edge into the adjacency list and the list of all edges
	graph[u] = append(graph[u], len(edges))
	edges = append(edges, edge{from: u, to: v, capacity: c, flow: 0})
	graph[v] = append(graph[v], len(edges))
	edges = append(edges, edge{from: v, to: u, capacity: 0, flow: 0})
}

func bfs() bool {
	/*
		Breadth-first search: here we calculate the distances from s to all other vertices
		and build a level graph.
	*/
	queue := make([]int, 0)
	queue = append(queue, s)
	for i := 0; i < n; i++ {
		distance[i] = -1
	}
	distance[s] = 0
	for len(queue) > 0 && distance[t] == -1 {
		u := queue[0]
		queue = queue[1:]
		// iterate over the adjacency list of vertex u
		for i := 0; i < len(graph[u]); i++ {
			edge_id := graph[u][i]
			v := edges[edge_id].to
			// if flow can be increased, update distance to v and add v to the queue
			if distance[v] == -1 && edges[edge_id].flow < edges[edge_id].capacity {
				queue = append(queue, v)
				distance[v] = distance[u] + 1
			}
		}
	}
	return distance[t] != -1 // return true if t can be reached
}

func dfs(u int, flow int) int {
	/*
		Depth-first search: find a blocking flow in the level graph and remove (by increasing pointer)
		all edges that do not lead to t.
	*/
	if flow == 0 || u == t {
		return flow
	}
	// iterate over the adjacency list of vertex u, starting with the first edge that has not yet been removed
	for pointer[u] < len(graph[u]) {
		edge_id := graph[u][pointer[u]]
		v := edges[edge_id].to
		// check if the edge is part of the level graph
		if distance[v] != distance[u]+1 {
			pointer[u] += 1
			continue
		}
		// calculate the flow that can be added
		added_flow := dfs(v, min(flow, edges[edge_id].capacity-edges[edge_id].flow))
		if added_flow != 0 {
			/*
				since we add corresponding edges simultaneously, their indices differ by one,
				so we can use xor to find the corresponding edge. (tip from e-maxx.ru)
			*/
			edges[edge_id].flow += added_flow
			edges[edge_id^1].flow -= added_flow
			return added_flow
		}
		pointer[u] += 1
	}
	return 0
}

func dinic() int {
	/*
		Dinic's algorithm: construct a level graph, find a blocking flow, repeat while
		the blocking flow exists.
	*/
	flow := 0
	// do algorithm's iterations until no blocking flow in the level graph exists anymore
	for true {
		// check if t is reachable from s, calculate distances from s to all other vertices and build a level graph
		if !bfs() {
			break
		}
		for i := 0; i < n; i++ {
			pointer[i] = 0
		}
		// find a blocking flow in the level graph
		added_flow := dfs(s, INF)
		for added_flow != 0 {
			flow += added_flow
			added_flow = dfs(s, INF)
		}
	}
	return flow
}

func input(filename string) {
	/*
		Input graph data from input file.
	*/
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	in := bufio.NewReader(f)
	fmt.Fscan(in, &n, &m, &s, &t)
	s -= 1
	t -= 1
	edges = []edge{}
	graph = make([][]int, n)
	for i := 0; i < m; i++ {
		var u, v, c int
		fmt.Fscan(in, &u, &v, &c)
		add_edge(u, v, c)
	}
}

func main() {
	input("example.txt")
	fmt.Println(dinic())
}
