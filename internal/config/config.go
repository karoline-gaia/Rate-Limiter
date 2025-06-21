package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	RateLimitIP        int
	BlockDurationIP    time.Duration
	RateLimitToken     map[string]int
	BlockDurationToken time.Duration
	RedisAddr          string
	ServerPort         string
}

func Load() Config {
	ipLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT_IP", "10"))
	blockIP, _ := strconv.Atoi(getEnv("BLOCK_DURATION_IP", "300"))
	blockToken, _ := strconv.Atoi(getEnv("BLOCK_DURATION_TOKEN", "120"))
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	port := getEnv("PORT", "8080")

	tokenLimits := make(map[string]int)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "RATE_LIMIT_TOKEN_") {
			parts := strings.SplitN(env, "=", 2)
			key := strings.TrimPrefix(parts[0], "RATE_LIMIT_TOKEN_")
			val, _ := strconv.Atoi(parts[1])
			tokenLimits[key] = val
		}
	}

	return Config{
		RateLimitIP:        ipLimit,
		BlockDurationIP:    time.Duration(blockIP) * time.Second,
		RateLimitToken:     tokenLimits,
		BlockDurationToken: time.Duration(blockToken) * time.Second,
		RedisAddr:          redisAddr,
		ServerPort:         port,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
