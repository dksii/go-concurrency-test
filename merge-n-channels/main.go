package main

import (
	"fmt"
	"sync"
)

func MergeNChannels[T any](channels ...<-chan T) <-chan T {
	output := make(chan T)
	wg := &sync.WaitGroup{}

	wg.Add(len(channels))

	for _, ch := range channels {
		go func(ch <-chan T) {
			defer wg.Done()

			for v := range ch {
				output <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

func main() {
	ch1, ch2, ch3 := make(chan int), make(chan int), make(chan int)

	sendFn := func(a, b int, ch chan<- int) {
		for i := a; i <= b; i++ {
			ch <- i
		}
		close(ch)
	}

	go sendFn(1, 5, ch1)
	go sendFn(6, 10, ch2)
	go sendFn(11, 15, ch3)

	for v := range MergeNChannels(ch1, ch2, ch3) {
		fmt.Println(v)
	}
}
