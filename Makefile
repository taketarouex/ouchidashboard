.PHONY: generate tidy test build

generate:
	go generate

tidy:
	go mod tidy

test:
	go test

build:
	go build cmd/main.go
