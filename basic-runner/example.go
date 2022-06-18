package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/kogutich/go-concurrency-test/basic-runner/runner"
)

func main() {
	r := runner.New[int](5)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	for i := 0; i < 10; i++ {
		j := i

		r.Add(func(context.Context) (int, error) {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
			return j, nil
		})
	}

	fmt.Println(r.Run(ctx))
}
