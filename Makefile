format:
	@go fmt ./...

vet: format
	@go vet ./...

run: format
	@DATABASE_URL=postgres://postgres:postgres@localhost:5432/softwarecraft?sslmode=disable APP_PORT=:3000 go run ./cmd/web

build: format
	@go build -o ./tmp/main ./cmd/web

test:
	@go test -v ./...

migrateup:
	@migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/softwarecraft?sslmode=disable" -verbose up

migratedown:
	@migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/softwarecraft?sslmode=disable" -verbose down

.PHONY: run vet format build test migrateup migratedown

