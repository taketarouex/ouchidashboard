.PHONY: generate tidy test build

generate:
	go generate ./collector

tidy:
	go mod tidy

test:
	go test ./collector

build:
	go build cmd/main.go
