package main

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"sort"
)

type StorageNode struct {
	Name string
	Host string
}

func hashFn(content []byte, totalSlots int) (int, error) {
	if totalSlots <= 0 {
		return 0, fmt.Errorf("totalSlots must be > 0")
	}

	sum := sha256.Sum256(content)        // [32]byte
	num := new(big.Int).SetBytes(sum[:]) // convert to big integer
	mod := new(big.Int).Mod(num, big.NewInt(int64(totalSlots)))

	return int(mod.Int64()), nil
}

type Ring struct {
	TotalSlots int
	Slots      []int          // sorted slot positions
	Lookup     map[int]string // slot -> nodeID/host
}

func (rg *Ring) AddNode(host string) error {
	key, err := hashFn([]byte(host), rg.TotalSlots)
	if err != nil {
		return err
	}
	rg.Slots = append(rg.Slots, key)
	if rg.Lookup == nil {
		rg.Lookup = make(map[int]string)
	}
	rg.Lookup[key] = host
	return nil
}

func (rg *Ring) DeleteNode(host string) (bool, error) {
	key, err := hashFn([]byte(host), rg.TotalSlots)
	if err != nil {
		return false, err
	}
	// find slot position in sorted slice
	i := sort.SearchInts(rg.Slots, key)
	if i >= len(rg.Slots) || rg.Slots[i] != key {
		return false, nil // not found
	}
	// delete from map and slice
	delete(rg.Lookup, key)
	rg.Slots = append(rg.Slots[:i], rg.Slots[i+1:]...)
	return true, nil
}

func (rg *Ring) Assign(content []byte) (string, error) {
	key, err := hashFn(content, rg.TotalSlots)
	if err != nil {
		return "", err
	}

	index := 0

	sort.Ints(rg.Slots)
	for _, slot_index := range rg.Slots {
		if key < slot_index {
			index = slot_index
			break
		}
	}

	if index == 0 {
		index = rg.Slots[0]
	}

	return rg.Lookup[index], nil

}

var total_slots int = 50

func main() {
	storageNodes := []StorageNode{
		{Name: "A", Host: "239.67.52.72"},
		{Name: "B", Host: "137.70.131.229"},
		{Name: "C", Host: "98.5.87.182"},
		{Name: "D", Host: "11.225.158.95"},
		{Name: "E", Host: "203.187.116.210"},
		{Name: "F", Host: "107.117.238.203"},
		{Name: "G", Host: "27.161.219.131"},
	}

	rg := &Ring{TotalSlots: 50}

	// Example: print nodes
	for _, node := range storageNodes {
		fmt.Printf("Node %s -> %s\n", node.Name, node.Host)
		rg.AddNode(node.Host)
	}

	sort.Ints(rg.Slots)
	for key, index := range rg.Slots {
		fmt.Printf("SLot Node Data %d -> %d\n", key, index)
	}

	// for key, index := range rg.Lookup {
	// 	fmt.Printf("Node Data %d -> %s\n", key, index)
	// }

	//Fixed Hashing, when node need to be add then data need to be migrated as per new node hash.
	filepath.WalkDir("dir", func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // skip directories
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// var sum int = 0

		// for _, b := range data {
		// 	sum += int(b)
		// }

		// hash_key := sum % len(storageNodes)

		//hash_key_c, err := hashFn(data, total_slots)

		//Add data in the ring

		host, err := rg.Assign(data)
		fmt.Println(host)

		return nil
	})

	//Consistent hash: we take very large has 2
}
