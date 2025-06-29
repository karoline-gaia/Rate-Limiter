package middleware

import (
	"net/http"
	"strings"
	"net"
	"fmt"
	"time"
)

type Limiter interface {
	Check(ip, token string) (bool, time.Duration)
}

func LimiterMiddleware(lim Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extrai IP real do header X-Forwarded-For, se existir
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			} else {
				// Pode conter múltiplos IPs, pega o primeiro
				ip = strings.Split(ip, ",")[0]
			}
			token := r.Header.Get("API_KEY")
			allowed, blockTime := lim.Check(ip, token)
			if !allowed {
				if blockTime > 0 {
					w.Header().Set("Retry-After", fmt.Sprintf("%d", int(blockTime.Seconds())))
				}
				w.WriteHeader(429)
				w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
