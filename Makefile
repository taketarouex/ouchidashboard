.PHONY: generate tidy test build

generate:
	go generate ./collector

tidy:
	go mod tidy

test:
	go test ./collector

build:
	go build -o build/. cmd/collector.go

integration_test:
	go test --tags=integration ./collector
