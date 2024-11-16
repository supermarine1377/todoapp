[![Go](https://github.com/supermarine1377/todoapp/actions/workflows/go.yml/badge.svg)](https://github.com/supermarine1377/todoapp/actions/workflows/go.yml)

## Introduction

I'm implementing basic TODO app using Go.

## How to run app

```sh
go run app/cmd/main.go
```

## Deelopment setup

Open this repository in VSCode and go to the Command pallete (Cmd+Shift+P) and type “Remote-Containers: Reopen in Container” and select it.

Or you can build Docker image with:
```sh
make dev
```

Then log in the container.

In the container, run this command:

```sh
make migrate-sqlite
```

## Development dependencies

- [golangci-lint](https://github.com/golangci/golangci-lint)
- [staticcheck](https://staticcheck.dev/docs/)
- [gomock](https://github.com/uber-go/mock)
- [asdf](https://asdf-vm.com/)
