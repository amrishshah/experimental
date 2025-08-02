package main

import "fmt"

func main() {
	card := "4539148803436467"
	double := false
	sum := 0
	for i := len(card) - 1; i >= 0; i-- {
		num := (int(card[i]) - '0')
		if num >= 0 {
			if double {
				num = 2 * num
				if num > 9 {
					//fmt.Println(num)
					//num = 1 + (num % 10)
					num = num - 9
				}
			}
			//fmt.Println(num)
			sum += num
			double = !double
		}

	}
	fmt.Println(sum % 10)
}
