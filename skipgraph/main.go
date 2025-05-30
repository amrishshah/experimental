package main

import (
	"fmt"
	"math/rand"
)

const MaxLevel = 5 // Maximum levels in the skip graph

// Node represents an element in the Skip Graph
type Node struct {
	key  int
	next []*Node // Pointers to next nodes at different levels
}

// SkipGraph represents the skip graph structure
type SkipGraph struct {
	head *Node
}

// NewNode creates a new node with a given key
func NewNode(key, level int) *Node {
	// Creates a new Node with the given key and level
	// key: The value stored in this node
	// level: The maximum level this node will participate in
	// Returns: A pointer to the new Node
	//
	// The node's next array has size level+1 because:
	// - Levels are 0-indexed (level 0 is the base list)
	// - We need space for levels 0 through 'level'
	// - So we need (level + 1) total spaces
	return &Node{
		key:  key,                    // Store the key value
		next: make([]*Node, level+1), // Create array of next pointers, initially all nil
	}
}

// NewSkipGraph initializes a Skip Graph
func NewSkipGraph() *SkipGraph {
	return &SkipGraph{
		head: NewNode(-1, MaxLevel), // Head node with -1 as a dummy key
	}
}

// RandomLevel generates a random level for node insertion
func RandomLevel() int {
	//rand.Seed(time.Now().UnixNano())
	level := 0
	for rand.Float32() < 0.5 && level < MaxLevel {
		level++
	}
	return level
}

// Insert adds a key to the Skip Graph
func (sg *SkipGraph) Insert(key int) {
	// Create an array to store the update positions at each level
	update := make([]*Node, MaxLevel+1)
	current := sg.head

	// Start from the highest level and find the right place to insert
	// For each level, traverse horizontally until we find a node with a larger key
	// This gives us O(log n) search time on average
	for i := MaxLevel; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
		// Store the node at this level where we need to update pointers
		update[i] = current
	}

	// Generate a random level for the new node (probabilistic balancing)
	// Each level has 50% chance of being included, up to MaxLevel
	level := RandomLevel()
	newNode := NewNode(key, level)

	// Insert the new node at each level up to its randomly chosen level
	// For each level:
	// 1. Set new node's next pointer to the node that update[i] was pointing to
	// 2. Update the update[i] node to point to our new node
	// This maintains the skip graph's linked structure at each level
	for i := 0; i <= level; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}
	fmt.Printf("Inserted: %d (Level: %d)\n", key, level)
}

// Search finds a key in the Skip Graph
func (sg *SkipGraph) Search(key int) bool {
	current := sg.head

	// Start from the highest level and move down
	for i := MaxLevel; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
	}

	// Move to the lowest level and check for key
	current = current.next[0]
	if current != nil && current.key == key {
		fmt.Printf("Found: %d\n", key)
		return true
	}
	fmt.Printf("Not Found: %d\n", key)
	return false
}

// Print displays the skip graph levels
func (sg *SkipGraph) Print() {
	for i := MaxLevel; i >= 0; i-- {
		current := sg.head.next[i]
		fmt.Printf("Level %d: ", i)
		for current != nil {
			fmt.Printf("%d -> ", current.key)
			current = current.next[i]
		}
		fmt.Println("nil")
	}
}

func main() {
	sg := NewSkipGraph()

	// Insert random keys
	keys := []int{10, 20, 30, 5, 15, 25, 35}
	for _, key := range keys {
		sg.Insert(key)
		sg.Print()
	}

	// Search for some keys
	sg.Search(15)
	sg.Search(40)

	// Print the Skip Graph structure
	//sg.Print()
}
