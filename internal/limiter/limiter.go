package limiter

import (
	"time"
)

type Limiter struct {
	store LimiterStore
	cfg   LimiterConfig
}

type LimiterConfig struct {
	RateLimitIP         int
	BlockDurationIP     time.Duration
	RateLimitToken      map[string]int
	BlockDurationToken  time.Duration
}

func NewLimiter(store LimiterStore, cfg LimiterConfig) *Limiter {
	return &Limiter{store: store, cfg: cfg}
}

// Check verifica se o IP ou Token excedeu o limite.
func (l *Limiter) Check(ip, token string) (allowed bool, blockTime time.Duration) {
	var key string
	var limit int
	var blockDuration time.Duration

	// Prioridade: Token
	if token != "" {
		if tLimit, ok := l.cfg.RateLimitToken[token]; ok {
			key = "token:" + token
			limit = tLimit
			blockDuration = l.cfg.BlockDurationToken
		} else {
			// Token não configurado, cai para IP
			key = "ip:" + ip
			limit = l.cfg.RateLimitIP
			blockDuration = l.cfg.BlockDurationIP
		}
	} else {
		key = "ip:" + ip
		limit = l.cfg.RateLimitIP
		blockDuration = l.cfg.BlockDurationIP
	}

	// Checar bloqueio
	blocked, ttl, err := l.store.IsBlocked(key)
	if err != nil {
		return false, 0
	}
	if blocked {
		return false, ttl
	}

	// Incrementar contador (janela de 1s)
	count, err := l.store.Increment(key+":count", time.Second)
	if err != nil {
		return false, 0
	}
	if count > limit {
		_ = l.store.SetBlock(key, blockDuration)
		return false, blockDuration
	}
	return true, 0
}
