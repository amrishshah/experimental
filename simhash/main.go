package main

import (
	"fmt"
	"hash/fnv"
	"strings"
)

// Simhash (return 64 bit number)
// --> String
// --> break word
// --> create hash (FNV-1a 64-bit hash object)
// --> for each token
// --> check if bit is set or not
// --> if set increment vector for that position
//

// It return finger print of 64 bit unsign
func Simhash(tokens []string) uint64 {
	//
	var v [64]int

	for _, token := range tokens {
		hash := hashToken(token)
		for i := 0; i < 64; i++ {
			if hash&(1<<i) != 0 {
				v[i]++
			} else {
				v[i]--
			}
		}
	}

	//Now for each words

	var fingerprint uint64
	for i := 0; i < 64; i++ {
		if v[i] > 0 {
			fingerprint |= (1 << i)
		}
	}

	return fingerprint
}

func hashToken(token string) uint64 {
	h := fnv.New64a()      // Create a new FNV-1a 64-bit hash object
	h.Write([]byte(token)) // Convert the token (string) to bytes and hash it
	return h.Sum64()       // Return the final 64-bit hash value
}

// HammingDistance computes the number of differing bits between two Simhash fingerprints.
func HammingDistance(hash1, hash2 uint64) int {
	xor := hash1 ^ hash2
	distance := 0
	for xor != 0 {
		distance++
		xor = xor & (xor - 1)
		fmt.Println(xor)
	}
	return distance
}
func main() {

	tokens1 := strings.Fields("This is a sample document for Simhash computation")
	tokens2 := strings.Fields("This is another example document for Simhash computation")

	hash1 := Simhash(tokens1)
	hash2 := Simhash(tokens2)

	fmt.Println(hash1)
	fmt.Println(hash2)

	distance := HammingDistance(hash1, hash2)

	fmt.Println(distance)

}
