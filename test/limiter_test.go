package test

import (
	"testing"
	"time"
	"Desafio-GO/internal/limiter"
)

type mockStore struct {
	count int
	blocked bool
}

func (m *mockStore) Increment(key string, expire time.Duration) (int, error) {
	m.count++
	return m.count, nil
}
func (m *mockStore) Get(key string) (int, error) { return m.count, nil }
func (m *mockStore) SetBlock(key string, d time.Duration) error { m.blocked = true; return nil }
func (m *mockStore) IsBlocked(key string) (bool, time.Duration, error) { return m.blocked, 10 * time.Second, nil }

func TestLimiter(t *testing.T) {
	store := &mockStore{}
	cfg := limiter.LimiterConfig{RateLimitIP: 2, BlockDurationIP: 10 * time.Second, RateLimitToken: map[string]int{"abc": 3}, BlockDurationToken: 10 * time.Second}
	lim := limiter.NewLimiter(store, cfg)

	// IP sem token
	for i := 0; i < 2; i++ {
		allowed, _ := lim.Check("1.1.1.1", "")
		if !allowed { t.Fatal("should allow") }
	}
	allowed, _ := lim.Check("1.1.1.1", "")
	if allowed { t.Fatal("should block after limit") }

	// Token
	store = &mockStore{}
	lim = limiter.NewLimiter(store, cfg)
	for i := 0; i < 3; i++ {
		allowed, _ := lim.Check("2.2.2.2", "abc")
		if !allowed { t.Fatal("should allow token") }
	}
	allowed, _ = lim.Check("2.2.2.2", "abc")
	if allowed { t.Fatal("should block token after limit") }
}
