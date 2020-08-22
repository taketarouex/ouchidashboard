.PHONY: generate tidy test

generate:
	go generate

tidy:
	go mod tidy

test:
	go test

build:
	go build cmd/main.go
