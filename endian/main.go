package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	// Example: Read first 2 bytes from a file
	// Define a 16-bit integer with value 0xFEFF
	var test uint16 = 0xFEFF

	fmt.Println(test)

	// Read the first 2 bytes
	bom := make([]byte, 2)
	//binary.LittleEndian.PutUint16(bom, test) // Store using Little-Endian format
	binary.BigEndian.PutUint16(bom, test)

	fmt.Println((bom))

	// Check for BOM
	if bom[0] == 0xFE && bom[1] == 0xFF {
		fmt.Println("File is in Big-Endian UTF-16")
	} else if bom[0] == 0xFF && bom[1] == 0xFE {
		fmt.Println("File is in Little-Endian UTF-16")
	} else {
		fmt.Println("No BOM found or unknown encoding")
	}
}
