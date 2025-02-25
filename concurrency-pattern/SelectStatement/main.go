package main

import "fmt"

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		ch1 <- "d"
	}()

	go func() {
		ch2 <- "d1"
	}()

	for i := 0; i < 2; i++ { //This line is important for getting data from all channel. Also if not added then
		select {
		case msg1 := <-ch1:
			fmt.Println(msg1)
		case msg2 := <-ch2:
			fmt.Println(msg2)
		}
	}
}
