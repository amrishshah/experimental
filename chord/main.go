package main

import (
	"fmt"
	"sort"
)

const m = 4             // bits
const ringSize = 1 << m // 16

type ChordRing struct {
	IDs     []int         // sorted list of node IDs
	Fingers map[int][]int // finger table: nodeID -> list of fingers
}

// Simplified version: global view
func (cr *ChordRing) findSuccessorSimple(id int) int {
	for _, nodeID := range cr.IDs {
		if nodeID >= id {
			return nodeID
		}
	}
	return cr.IDs[0] // wrap around
}

// Real Chord version: finger-based routing
func (cr *ChordRing) findSuccessorRecursive(startNode, id int) int {
	succ := cr.findSuccessorSimple((startNode + 1) % ringSize)
	if inInterval(id, startNode, succ, true) {
		return succ
	} else {
		closest := cr.closestPrecedingNode(startNode, id)
		if closest == startNode {
			return succ // fallback if no closer node
		}
		return cr.findSuccessorRecursive(closest, id)
	}
}

// Finger-based closest preceding node
func (cr *ChordRing) closestPrecedingNode(n, id int) int {
	for i := m - 1; i >= 0; i-- {
		finger := cr.Fingers[n][i]
		if inInterval(finger, n, id, false) {
			return finger
		}
	}
	return n
}

// Interval check: (a, b] or (a, b)
func inInterval(x, a, b int, inclusive bool) bool {
	if a < b {
		if inclusive {
			return x > a && x <= b
		} else {
			return x > a && x < b
		}
	} else {
		if inclusive {
			return x > a || x <= b
		} else {
			return x > a || x < b
		}
	}
}

func main() {
	cr := ChordRing{
		IDs:     []int{0, 3, 6, 9, 14},
		Fingers: make(map[int][]int),
	}

	// Sort IDs to simulate a proper ring
	sort.Ints(cr.IDs)

	// Manually create a simple finger table for demo (real one uses 2^i jumps)
	for _, id := range cr.IDs {
		var fingers []int
		for i := 0; i < m; i++ {
			fingerID := (id + (1 << i)) % ringSize
			fingers = append(fingers, cr.findSuccessorSimple(fingerID))
		}
		cr.Fingers[id] = fingers
	}

	fmt.Println("Node finger tables:")
	for id, fingers := range cr.Fingers {
		fmt.Printf("Node %d: %v\n", id, fingers)
	}

	// Try resolving key 11 from node 3
	key := 11
	fmt.Printf("\n[Simple] Successor of %d: %d\n", key, cr.findSuccessorSimple(key))
	fmt.Printf("[Chord ] Successor of %d (from Node 3): %d\n", key, cr.findSuccessorRecursive(3, key))
}
