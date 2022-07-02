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

// splits channel (outputs 1 chan to 2 chans)
func tee[T any](
	done <-chan struct{},
	in <-chan T,
) (_, _ <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)

	go func() {
		defer close(out1)
		defer close(out2)

		for val := range orDone(done, in) {
			var out1, out2 = out1, out2

			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				}
			}
		}
	}()

	return out1, out2
}

func main() {
}
