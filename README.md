# Instalação
```sh
go mod tidy
```

# Execução
1. cp .env_sample .env
2. Preencher o GITHUB_TOKEN no .env com seu token

Minerar os dados:
```sh
go run ./cmd/fetch/fetch.go
```

Leitura dos dados:
```sh
go run ./cmd/read/read.go
```
