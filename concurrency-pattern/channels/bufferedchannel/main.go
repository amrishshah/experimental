package main

import "fmt"

func main() {
	// Create a buffered channel
	ch := make(chan string, 3)

	// Send data to the channel
	ch <- "Go"
	ch <- "Concurrency"
	ch <- "Channel"

	close(ch) // Close the channel

	// Loop to receive data
	for val := range ch {
		fmt.Println(val)
	}
}
