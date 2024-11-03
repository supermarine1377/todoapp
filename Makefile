.PHONY: dev
dev: 
	docker build --pull --rm -f "dev.Dockerfile" -t todoapp:dev "."
	docker run -it -v .:/app -v vscode-server:/root/.vscode-server --rm todoapp:dev

generate:
	go generate ./...

test:
	go test ./...

lint:
	golangci-lint run

migrate-sqlite:
	sqlite3def _data/sqlite.db -f _migration/sqlite/schema.sql 
