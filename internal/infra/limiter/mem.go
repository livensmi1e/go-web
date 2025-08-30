package limiter

import (
	"context"
	"go-web/internal/core/ports"
	"sync"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter *rate.Limiter
}

type memLimiter struct {
	visitors map[string]*visitor
	r        rate.Limit
	b        int
	mu       sync.Mutex
}

func NewMemLimiter(r rate.Limit, b int) ports.RateLimiter {
	l := &memLimiter{
		visitors: make(map[string]*visitor),
		r:        r,
		b:        b,
	}
	return l
}

func (l *memLimiter) getVisitor(key string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	v, exists := l.visitors[key]
	if !exists {
		limiter := rate.NewLimiter(l.r, l.b)
		l.visitors[key] = &visitor{limiter: limiter}
		return limiter
	}

	return v.limiter
}

func (l *memLimiter) Allow(ctx context.Context, key string) (bool, error) {
	limiter := l.getVisitor(key)
	return limiter.Allow(), nil
}
