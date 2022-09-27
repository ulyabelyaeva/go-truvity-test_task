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

var n, m, s, t int               // number of vertices, number of edges, source and drain vertices
var graph = make([][]int, 0)     // in form of an adjacency list
var graph_inv = make([][]int, 0) // inverted edges
var distance [MAXN]int           // distances from s to all other vertices, required to build a level graph
var potential [MAXN]int          // potential of a vertex
var banned [MAXN]bool            // true if vertex is temporarily deleted\
var pointer [MAXN]int            // pointer to the first not deleted edge
var pointer_inv [MAXN]int        // same for the graph with inverted edges

type edge struct {
	from, to, capacity, flow int
}

var edges = []edge{}

func add_edge(u, v, c int) {
	u -= 1
	v -= 1
	// add the edge into the adjacency list and the list of all edges
	graph[u] = append(graph[u], len(edges))
	graph_inv[v] = append(graph_inv[v], len(edges))
	edges = append(edges, edge{from: u, to: v, capacity: c, flow: 0})
	graph[v] = append(graph[v], len(edges))
	graph_inv[u] = append(graph_inv[u], len(edges))
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
		banned[i] = false
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

func calc_potential(u int) int {
	/*
		Calculate potential of vertex u
	*/
	sum_out := 0
	for i := 0; i < len(graph[u]); i++ {
		edge_id := graph[u][i]
		v := edges[edge_id].to
		if distance[v] != distance[u]+1 || banned[v] {
			continue
		}
		sum_out += edges[edge_id].capacity - edges[edge_id].flow
	}
	sum_in := 0
	for i := 0; i < len(graph_inv[u]); i++ {
		edge_id := graph_inv[u][i]
		v := edges[edge_id].from
		if distance[v] != distance[u]-1 || banned[v] {
			continue
		}
		sum_in += edges[edge_id].capacity - edges[edge_id].flow
	}
	if u == t {
		return sum_in
	}
	if u == s {
		return sum_out
	}
	return min(sum_in, sum_out)
}

func dfs(u int, flow int, inv bool) {
	/*
		Depth-first search: find a blocking flow in the level graph and remove (by increasing pointer)
		all edges that do not lead to t.
	*/
	if flow == 0 || (!inv && u == t) || (inv && u == s) {
		return
	}
	// iterate over the adjacency list of vertex u, starting with the first edge that has not yet been removed
	if !inv {
		for pointer[u] < len(graph[u]) {
			edge_id := 0
			v := 0
			edge_id = graph[u][pointer[u]]
			v = edges[edge_id].to
			if distance[v] != distance[u]+1 || banned[v] {
				pointer[u] += 1
				continue
			}
			// calculate the flow that can be added
			added_flow := min(flow, edges[edge_id].capacity-edges[edge_id].flow)
			flow -= added_flow
			dfs(v, added_flow, inv)
			if added_flow != 0 {
				edges[edge_id].flow += added_flow
				edges[edge_id^1].flow -= added_flow
			} else if flow != 0 {
				pointer[u] += 1
			}
			if flow == 0 {
				break
			}
		}
	} else {
		for pointer_inv[u] < len(graph_inv[u]) {
			edge_id := 0
			v := 0
			edge_id = graph_inv[u][pointer_inv[u]]
			v = edges[edge_id].from
			if distance[v] != distance[u]-1 || banned[v] {
				pointer_inv[u] += 1
				continue
			}
			// calculate the flow that can be added
			added_flow := min(flow, edges[edge_id].capacity-edges[edge_id].flow)
			flow -= added_flow
			dfs(v, added_flow, inv)
			if added_flow != 0 {
				edges[edge_id].flow += added_flow
				edges[edge_id^1].flow -= added_flow
			} else if flow != 0 {
				pointer_inv[u] += 1
			}
			if flow == 0 {
				break
			}
		}
	}
}

func recalc_potential(u int, inv bool) {
	/*
		If we've found a vertex with zero potential, we should delete it and
		recalculate potentials of adjacent vertices. This function does that.
	*/
	queue := make([]int, 0)
	queue = append(queue, u)
	banned[u] = true
	// some kind of bfs again
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		len_graph_u := len(graph[u])
		if inv {
			len_graph_u = len(graph_inv[u])
		}
		// iterate over the adjacency list of vertex u
		for i := 0; i < len_graph_u; i++ {
			edge_id := 0
			v := 0
			if !inv {
				edge_id = graph[u][i]
				v = edges[edge_id].to
				if distance[v] != distance[u]+1 {
					continue
				}
			} else {
				edge_id = graph_inv[u][i]
				v = edges[edge_id].from
				if distance[v] != distance[u]-1 {
					continue
				}
			}
			potential[v] = calc_potential(v)
			// check if v has to be deleted
			if s != v && t != v && potential[v] == 0 && !banned[v] {
				banned[v] = true
				queue = append(queue, v)
			}
		}
	}
}

func find_blocking_flow() int {
	/*
		Calculate potentials, find a ('best') vertex with the smallest non-zero potential,
		update flow through this vertex.
	*/
	// calculate potentials
	for u := 0; u < n; u++ {
		potential[u] = calc_potential(u)
	}
	// remove all vertices with zero potential
	for u := 0; u < n; u++ {
		if potential[u] == 0 && !banned[u] {
			recalc_potential(u, false)
			recalc_potential(u, true)
		}
	}
	// find a vertex with smallest non-zero potential
	u_best := s
	for u := 0; u < n; u++ {
		if potential[u] != 0 && potential[u] < potential[u_best] {
			u_best = u
		}
	}
	// update flow through the 'best' vertex
	// from u_best to f:
	//update_flow(u_best, false)
	dfs(u_best, potential[u_best], false)
	// from u_best to s via inverted edges
	//update_flow(u_best, true)
	dfs(u_best, potential[u_best], true)
	return potential[u_best]
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
			pointer_inv[i] = 0
		}
		// find a blocking flow in the level graph
		added_flow := find_blocking_flow()
		for added_flow != 0 {
			flow += added_flow
			added_flow = find_blocking_flow()
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
	graph_inv = make([][]int, n)
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
