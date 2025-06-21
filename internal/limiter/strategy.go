package limiter

import "time"

type LimiterStore interface {
	Increment(key string, expire time.Duration) (int, error)
	Get(key string) (int, error)
	SetBlock(key string, duration time.Duration) error
	IsBlocked(key string) (bool, time.Duration, error)
}
