build:
	@go build -o bin/json_parser cmd/main.go

run: build
	@./bin/json_parser

test:
	@go test -v ./... -short
