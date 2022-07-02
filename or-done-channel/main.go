package main

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

func main() {
	// usage:
	// for val := range orDone(done, myChan) {
	// // Do something with val
	// }
}
