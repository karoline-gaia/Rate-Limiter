package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"Desafio-GO/internal/config"
	"Desafio-GO/internal/limiter"
	"Desafio-GO/internal/middleware"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	store, err := limiter.NewRedisStore(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}
	limiterCfg := limiter.LimiterConfig{
		RateLimitIP:        cfg.RateLimitIP,
		BlockDurationIP:    cfg.BlockDurationIP,
		RateLimitToken:     cfg.RateLimitToken,
		BlockDurationToken: cfg.BlockDurationToken,
	}
	lim := limiter.NewLimiter(store, limiterCfg)
	limiterMw := middleware.LimiterMiddleware(lim)

	http.Handle("/", limiterMw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK!"))
	})))

	log.Printf("Servidor iniciado na porta %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
