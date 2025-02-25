package main

func temp(ch chan int) {
	ch <- 43
	ch <- 44
	close(ch)
}

func main() {
	ch := make(chan int, 100)

	go temp(ch)

	for val := range ch {
		println(val)
	}
}
