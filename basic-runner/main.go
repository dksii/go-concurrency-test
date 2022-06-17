package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Result[T any] struct {
	Value T
	Error error
}

type Fn[T any] func() (T, error)

type Runner[T any] struct {
	functions []Fn[T]
}

func (r *Runner[T]) Add(f Fn[T]) {
	r.functions = append(r.functions, f)
}

func (r *Runner[T]) Reset() {
	r.functions = nil
}

func (r *Runner[T]) Run() []Result[T] {
	wg := &sync.WaitGroup{}
	m := &sync.Mutex{}
	result := make([]Result[T], 0, len(r.functions))

	wg.Add(len(r.functions))

	for _, f := range r.functions {
		go func(f Fn[T]) {
			defer wg.Done()

			v, err := f()
			m.Lock()
			result = append(result, Result[T]{v, err})
			m.Unlock()
		}(f)
	}

	wg.Wait()
	return result
}

func main() {
	r := &Runner[int]{}

	for i := 0; i < 10; i++ {
		j := i

		r.Add(func() (int, error) {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
			return j, nil
		})
	}

	fmt.Println(r.Run())
}
