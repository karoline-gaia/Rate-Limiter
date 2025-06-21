# Rate Limiter Go

## Objetivo
Rate limiter configurável por IP ou Token, usando Redis como backend, pronto para rodar com Docker Compose.

## Como funciona
- Limita requisições por IP ou por token (token tem precedência)
- Configuração via variáveis de ambiente ou arquivo `.env`
- Persistência dos contadores e bloqueios em Redis
- Middleware HTTP pronto para uso
- Fácil trocar Redis por outro backend implementando a interface `LimiterStore`

## Configuração
Veja o arquivo `.env.example`:

```
RATE_LIMIT_IP=10
BLOCK_DURATION_IP=300
RATE_LIMIT_TOKEN_abc123=100
BLOCK_DURATION_TOKEN=120
REDIS_ADDR=redis:6379
PORT=8080
```

## Uso
1. Suba o Redis e o app:
   ```
   docker-compose up --build
   ```
2. Faça requisições para `localhost:8080`.
   - Para testar limitação por token, envie header `API_KEY: <TOKEN>`

## Testes
Execute os testes com:
```
go test ./test/...
```

## Trocar backend
Implemente a interface `LimiterStore` em `internal/limiter/strategy.go`.

## Resposta ao limite excedido
- HTTP 429
- Mensagem: `you have reached the maximum number of requests or actions allowed within a certain time frame`
