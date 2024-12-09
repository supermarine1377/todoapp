FROM --platform=arm64 golang:1.23

RUN useradd -ms /bin/bash dev && \
    chown -R dev:dev /go && \
    chown dev:dev /usr/local/bin

USER dev

# Install development dependencies    
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install go.uber.org/mock/mockgen@v0.5.0 && \
    curl -L https://github.com/sqldef/sqldef/releases/download/v0.17.23/sqlite3def_linux_amd64.tar.gz \
        -o sqlite3def_linux_amd64.tar.gz && \
    tar xz -C /usr/local/bin -f sqlite3def_linux_amd64.tar.gz && \
    go install -v github.com/cweill/gotests/gotests@v1.6.0 && \
    go install github.com/swaggo/swag/cmd/swag@latest
  
ENV DATABASE_DSN="/app/_data/sqlite.db"
ENV TZ="Asia/Tokyo"

WORKDIR /app/
