[![Go](https://github.com/supermarine1377/todoapp/actions/workflows/go.yml/badge.svg)](https://github.com/supermarine1377/todoapp/actions/workflows/go.yml)

# Introduction

I'm implementing basic TODO app using Go.

# How to run app

```sh
go run app/cmd/main.go
```

# How to run server with gracefully shutdown feature

```sh
go run app/cmd/main.go --gracefully-shutdown 
```

Send requrest like this:
```sh
curl -i localhost:8080/healthz
```

ターミナルで^+cしても、サーバーはすぐに終了せず、リクエストに対する処理を終了させてからサーバーが終了します