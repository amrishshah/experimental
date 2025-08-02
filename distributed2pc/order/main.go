package main

import (
	"log"
	"sync"

	order "github.com/amrishshah/distributed2pc/order/svc"
)

func main() {
	foodID := 1
	var wg sync.WaitGroup
	wg.Add(10) // 100 concurrent orders

	for i := 0; i < 10; i++ {
		go func() {
			order, err := order.PlaceOrder(foodID)
			defer wg.Done()
			if err != nil {
				log.Println("order not placed:", err.Error())
			} else {
				log.Println("order placed:", order.ID)
			}
		}()
	}

	wg.Wait()
}
