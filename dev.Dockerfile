FROM --platform=arm64 golang:1.23

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 && \
    go install github.com/golang/mock/mockgen@v1.6.0 && \
    curl -L https://github.com/sqldef/sqldef/releases/download/v0.17.23/sqlite3def_linux_amd64.tar.gz | tar xz -C /usr/local/bin

WORKDIR /app/