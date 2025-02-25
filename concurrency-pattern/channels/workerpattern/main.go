package main

import (
	"fmt"
	"sync"
	"time"
)

// <-chan int -- receive-only channel (read)
// chan<- int -- send-only channel (write)
func workerTask(id int, tasks <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d processing task %d\n", id, task)
		time.Sleep(time.Second) // Simulate work
		results <- task * 2     // Send result back
	}
}

func main() {
	const numWorkers = 3
	const numTasks = 9
	var wg sync.WaitGroup

	//Buffered channel
	tasks := make(chan int, (numTasks))
	results := make(chan int, numTasks)

	//creating go routine.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // this is need to add for worker.
		go workerTask(i, tasks, results, &wg)
	}

	//Sending Task to channel
	for i := 0; i < numTasks; i++ {
		tasks <- i
	}
	close(tasks) // no more task;

	//Wait for task to get complete
	go func() {
		wg.Wait()
		close(results)
	}()

	//Getting results from channel -- This will print data as when available in channel.
	for result := range results {
		fmt.Println("Result:", result)
	}
}
