package main

import "fmt"

func orDone[T any](done <-chan struct{}, input <-chan T) <-chan T {
	output := make(chan T)

	go func() {
		defer close(output)

		for {
			select {
			case <-done:
				return
			case v, ok := <-input:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case output <- v:
				}
			}
		}
	}()

	return output
}

// reads values from channel of channels
func bridge[T any](
	done <-chan struct{},
	chanStream <-chan <-chan T,
) <-chan T {
	valStream := make(chan T)
	go func() {
		defer close(valStream)
		for {
			var stream <-chan T
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func main() {
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}
	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
}
