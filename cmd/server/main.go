package main

import (
	"log"
	"net/http"
	"os"

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
	lim := limiter.NewLimiter(store, cfg)
	limiterMw := middleware.LimiterMiddleware(lim)

	http.Handle("/", limiterMw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK!"))
	})))

	log.Printf("Servidor iniciado na porta %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, nil))
}
