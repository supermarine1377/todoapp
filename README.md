# TODO App

[![Go](https://github.com/supermarine1377/todoapp/actions/workflows/go.yml/badge.svg)](https://github.com/supermarine1377/todoapp/actions/workflows/go.yml)

[![Go Report Card](https://goreportcard.com/badge/github.com/supermarine1377/todoapp)](https://goreportcard.com/report/github.com/supermarine1377/todoapp)

## Introduction

This repository contains a basic TODO application implemented in Go. The application provides a foundation for learning Go development and cloud-native practices.

## How to Run the App

To run the application, execute the following command:

```sh
go run app/cmd/main.go
```

## Development Setup

### Using VSCode with Dev Containers

1. Open this repository in VSCode.
2. Open the Command Palette (Cmd+Shift+P) and search for "Remote-Containers: Reopen in Container." Select it to start working in the development container.

### Using Docker

1. Build the Docker image with:
   ```sh
   make dev
   ```
2. Log into the container.
3. Inside the container, run the following command to initialize the SQLite database:
   ```sh
   make migrate-sqlite
   ```

## Key Features

### API Endpoints

This application provides several API endpoints. The API specification is available [here](https://supermarine1377.github.io/todoapp/).

### Continuous Integration

The application is designed for easy testing and maintenance, featuring comprehensive unit and end-to-end test coverage. All tests are executed automatically via GitHub Actions.

### Cloud-Native Design

The application is being developed with cloud-native principles in mind to ensure smooth deployment in cloud environments. It includes:

- **Graceful Shutdown**: Handles termination signals to ensure proper cleanup of resources.
- **Configuration Separation**: Allows easy management of environment-specific configurations.

## Development Dependencies

To contribute or develop the application, ensure the following dependencies are installed:

- [golangci-lint](https://github.com/golangci/golangci-lint): For linting Go code.
- [staticcheck](https://staticcheck.dev/docs/): For static code analysis.
- [gomock](https://github.com/uber-go/mock): For generating and using mock objects in tests.
- [asdf](https://asdf-vm.com/): For managing tool versions.

---

For more information, please refer to the documentation and API specification linked above.
