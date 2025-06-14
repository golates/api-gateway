# Migration
MIGRATION_NAME ?= $(shell read -p "Migration name: " migration_name; echo $$migration_name)
DB_CONNECTION ?= postgres://postgres:root@localhost:5432/go-api?sslmode=disable

build:
	go build -o bin/api cmd/main.go
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/api cmd/api/main.go
test:
	go test ./...
test-cover:
	go test ./... -coverprofile=coverage.out
test-cover-open:
	go tool cover -html=coverage.out
run:
	go run cmd/main.go
create-migration:
	@migrate create -ext sql -dir migrations -seq $(MIGRATION_NAME)
up-migration:
	@migrate -database $(DB_CONNECTION) -path migrations up
down-migration:
	@migrate -database $(DB_CONNECTION) -path migrations down
generate-grpc:
	@protoc --go_out=./  --go_opt=MarshalOptions.emitUnpopulated=true --go-grpc_out=./ proto/*.proto