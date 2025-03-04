package main

import (
	"fmt"
	"hash"

	"github.com/spaolacci/murmur3"
)

var mhasher hash.Hash32

func init() {
	fmt.Println("hi 1")
	mhasher = murmur3.New32WithSeed(uint32(10))
}

func murmurhash(key string, size int32) int32 {
	mhasher.Write([]byte(key))
	result := mhasher.Sum32() % uint32(size)
	return int32(result)
}

type BloomFilter struct {
	filter []bool
	size   int32
}

func NewBloomFilter(size int32) *BloomFilter {
	return &BloomFilter{
		filter: make([]bool, size),
		size:   size,
	}
}

func (b *BloomFilter) Add(key string) {
	//Hash String to INT32
	idx := murmurhash(key, b.size)
	b.filter[idx] = true
}

func (b *BloomFilter) Exists(key string) bool {
	idx := murmurhash(key, b.size)
	return b.filter[idx]
}

func (b *BloomFilter) print() {
	fmt.Println(b.filter)
}

func main() {
	b := NewBloomFilter(10)
	var keys = []string{"a", "b", "c"}

	for _, key := range keys {
		b.Add(key)
	}

	b.print()
}
