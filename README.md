# Rate Limiter Go

## Objetivo
Rate limiter configurável por IP ou Token, usando Redis como backend, pronto para rodar com Docker Compose.

## Pré-requisitos

- Docker e Docker Compose instalados corretamente.
- Go (opcional, apenas para rodar/testar fora do Docker).

### Instalação do Docker (Linux)

```sh
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io
sudo apt-get install docker-compose-plugin
```

Se houver conflitos de pacotes (como `docker-buildx` ou `docker-compose-plugin`), remova os pacotes conflitantes antes de instalar.

## Como funciona
- Limita requisições por IP ou por token (token tem precedência)
- Configuração via variáveis de ambiente ou arquivo `.env`
- Persistência dos contadores e bloqueios em Redis
- Middleware HTTP pronto para uso
- Fácil trocar Redis por outro backend implementando a interface `LimiterStore`

## Configuração

Copie `.env.example` para `.env` e ajuste conforme necessário:

```
RATE_LIMIT_IP=10                    # Limite de requisições por IP por segundo
BLOCK_DURATION_IP=300               # Tempo de bloqueio para IP em segundos
RATE_LIMIT_TOKEN_abc123=100         # Limite para o token abc123
BLOCK_DURATION_TOKEN=120            # Tempo de bloqueio para tokens em segundos
REDIS_ADDR=redis:6379               # Endereço do Redis (não altere se for usar o docker-compose)
PORT=8080                           # Porta do servidor
```

## Execução

1. **Suba o Redis e o app:**
   ```sh
   docker compose up --build
   ```
   Se ocorrer erro de permissão, tente:
   ```sh
   sudo docker compose up --build
   ```

2. O servidor estará disponível em [http://localhost:8080](http://localhost:8080).

## Testando o Rate Limiter

- **Por IP:** até 10 requisições por segundo por IP.
- **Por Token:** envie o header `API_KEY: abc123` para até 100 requisições por segundo.

### Exemplo com curl

```sh
curl -H "API_KEY: abc123" http://localhost:8080/
```

Se exceder o limite, a resposta será:
```
HTTP/1.1 429 Too Many Requests
Retry-After: 300
you have reached the maximum number of requests or actions allowed within a certain time frame
```

## Testes
Execute os testes com:
```
go test ./test/...
```

## Trocar backend
Implemente a interface `LimiterStore` em `internal/limiter/strategy.go`.

## Resolução de Problemas

- **Erro de conflito de pacotes Docker:**  
  Remova pacotes conflitantes antes de instalar novos:
  ```sh
  sudo apt-get remove docker-buildx-plugin docker-compose-plugin
  sudo apt-get install docker-ce docker-ce-cli docker-compose-plugin
  ```

- **Permissão negada ao rodar Docker:**  
  Adicione seu usuário ao grupo docker:
  ```sh
  sudo usermod -aG docker $USER
  ```

- **Dúvidas ou problemas:**  
  Consulte a [documentação oficial do Docker](https://docs.docker.com/).

## Resposta ao limite excedido
- HTTP 429
- Mensagem: `you have reached the maximum number of requests or actions allowed within a certain time frame`
