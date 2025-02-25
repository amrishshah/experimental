package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

// <-chan int -- receive-only channel (read)
// chan<- int -- send-only channel (write)
func workerTask(id int, tasks <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d processing task %s\n", id, task)
		//time.Sleep(time.Second) // Simulate work
		results <- task // Send result back
	}
}

func main() {
	const numWorkers = 2
	const numTasks = 9
	var wg sync.WaitGroup

	//Buffered channel
	tasks := make(chan string, numTasks)
	results := make(chan string, numTasks)

	//creating go routine.
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // this is need to add for worker.
		go workerTask(i, tasks, results, &wg)
	}

	filePath := "large.csv"
	file, err := os.Open(filePath)

	if err != nil {
		log.Panicln(err)
	}

	//Wait for task to get complete
	go func() {
		defer file.Close()
		defer close(tasks)

		scanner := bufio.NewScanner(file)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer for long lines
		//var lines []string
		for scanner.Scan() {
			//println(lines)
			tasks <- scanner.Text()
			//lines = append(lines, scanner.Text())
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	//Getting results from channel -- This will print data as when available in channel.
	for result := range results {
		fmt.Println("Result:", result)
	}

	//println()
}
