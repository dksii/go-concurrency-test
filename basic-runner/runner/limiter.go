package runner

import (
	"context"
	"log"
)

// Limiter limits parallel actions.
type Limiter struct {
	quota chan struct{}
}

func NewLimiter(limit uint) *Limiter {
	return &Limiter{
		quota: make(chan struct{}, limit),
	}
}

func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case l.quota <- struct{}{}:
		return nil
	}
}

func (l *Limiter) Release() {
	select {
	case <-l.quota:
	default:
		log.Println("[WARN]: Release() call fallback to default case")
	}
}
