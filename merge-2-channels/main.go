package main

import "fmt"

func MergeTwoChannels[T any](ch1, ch2 <-chan T) <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)

		for ch1 != nil || ch2 != nil {
			select {
			case v1, ok := <-ch1:
				if ok {
					output <- v1
				} else {
					ch1 = nil
				}
			case v2, ok := <-ch2:
				if ok {
					output <- v2
				} else {
					ch2 = nil
				}
			}
		}
	}()

	return output
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	sendFn := func(a, b int, ch chan<- int) {
		for i := a; i <= b; i++ {
			ch <- i
		}
		close(ch)
	}

	go sendFn(1, 10, ch1)
	go sendFn(11, 20, ch2)

	for v := range MergeTwoChannels(ch1, ch2) {
		fmt.Println(v)
	}
}
