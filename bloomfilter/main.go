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
	filter []uint8
	size   int32
}

func NewBloomFilter(size int32) *BloomFilter {
	return &BloomFilter{
		filter: make([]uint8, size),
		size:   size,
	}
}

func (b *BloomFilter) Add(key string) {
	//Hash String to INT32
	idx := murmurhash(key, b.size)
	fmt.Println(idx)
	aidx := idx / 8 // array index

	bidx := idx % 8

	b.filter[aidx] = b.filter[aidx] | (1 << bidx)
}

func (b *BloomFilter) Exists(key string) bool {
	idx := murmurhash(key, b.size)
	aidx := idx / 8 // array index
	bidx := idx % 8
	return b.filter[aidx]&(1<<bidx) > 0
}

func (b *BloomFilter) print() {
	fmt.Println(b.filter)
}

func main() {
	b := NewBloomFilter(10)
	var keys = []string{"aqwertyuiuytrewq", "b", "c"}

	for _, key := range keys {
		b.Add(key)
	}

	b.print()
}
