package runner

import (
	"context"
	"sync"
)

type Result[T any] struct {
	Value T
	Error error
}

type Fn[T any] func(context.Context) (T, error)

// Runner runs functions in parallel.
// Supports limiting the number of running goroutines.
type Runner[T any] struct {
	functions []Fn[T]
	limiter   *Limiter
}

func New[T any](limit uint) *Runner[T] {
	return &Runner[T]{limiter: NewLimiter(limit)}
}

func (r *Runner[T]) Add(f Fn[T]) {
	r.functions = append(r.functions, f)
}

func (r *Runner[T]) Reset() {
	r.functions = nil
}

// Run runs all added functions in parallel and blocks until all executed (or context.Context cancelled).
func (r *Runner[T]) Run(ctx context.Context) []Result[T] {
	wg := &sync.WaitGroup{}
	m := &sync.Mutex{}
	result := make([]Result[T], 0, len(r.functions))

	wg.Add(len(r.functions))

	for _, f := range r.functions {
		if err := r.limiter.Acquire(ctx); err != nil {
			var zero T
			m.Lock()
			result = append(result, Result[T]{zero, err})
			m.Unlock()
			wg.Done()
			continue
		}

		go func(f Fn[T]) {
			defer wg.Done()
			defer r.limiter.Release()

			v, err := f(ctx)
			m.Lock()
			result = append(result, Result[T]{v, err})
			m.Unlock()
		}(f)
	}

	wg.Wait()
	return result
}
