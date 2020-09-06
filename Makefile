.PHONY: generate tidy test e2e_test

generate:
	go generate ./collector

tidy:
	go mod tidy

test:
	go test ./collector

build/collector:
	go build -o build/. cmd/collector.go

integration_test:
	go test --tags=integration ./collector

e2e_test:
	go test ./e2e_test
